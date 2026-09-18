package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"neomatica/neosync-tcp/infra/messaging/rabbitmq"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/pkg/protocol/common"
	"neomatica/neosync-tcp/types"
	"neomatica/neosync-tcp/util"
	"net"
	"net/http"
	"sync"
	"time"
)

type SessionMemory struct {
	store      store.Storage
	rmq        *rabbitmq.RabbitMQ
	LimiterDb  chan struct{}
	sessions   sync.Map
	requestMap sync.Map
}

func NewSessionMemory(store store.Storage, rmq *rabbitmq.RabbitMQ, limiterDb chan struct{}) *SessionMemory {
	return &SessionMemory{
		store:      store,
		rmq:        rmq,
		LimiterDb:  limiterDb,
		sessions:   sync.Map{},
		requestMap: sync.Map{},
	}
}

/* NewDeviceSession добавить/обновить сессию устройства */
func (sm *SessionMemory) NewDeviceSession(imei string, device *adm.ADMDevice) {
	imeiPrepare := util.PrepareImei(imei)

	if oldValue, ok := sm.sessions.LoadAndDelete(imeiPrepare); ok {
		sm.shutdownSession(oldValue.(*types.TCPSession), imeiPrepare, false)
	}

	session := &types.TCPSession{
		DeviceID:  device.DeviceId,
		RequestId: nil,
		Device:    device,
		Busy:      false,
		InFlight:  nil,
	}

	sm.sessions.Store(imeiPrepare, session)

	if err := sm.store.Devices.Update_DeviceStatusByImei(context.Background(), imeiPrepare, true); err != nil {
		logger.Error("NewDeviceSession ip={%s}: Failed to update device status true/false: %s", device.Conn.RemoteAddr().(*net.TCPAddr).IP, err.Error())
	}
}

/* ConnectedImeis IMEI с активной TCP-сессией */
func (sm *SessionMemory) ConnectedImeis() []string {
	imeis := make([]string, 0)

	sm.sessions.Range(func(key, _ any) bool {
		if imei, ok := key.(string); ok && imei != "" {
			imeis = append(imeis, imei)
		}

		return true
	})

	return imeis
}

/* GetSession получить сессию по IMEI */
func (sm *SessionMemory) GetDeviceSession(imei string) (*types.TCPSession, bool) {
	imeiPrepare := util.PrepareImei(imei)

	v, ok := sm.sessions.Load(imeiPrepare)
	if !ok {
		return nil, false
	}

	session := v.(*types.TCPSession)

	session.Mu.Lock()
	defer session.Mu.Unlock()

	return session, true
}

/* RemoveDeviceSession удалить сессию устройства */
func (sm *SessionMemory) RemoveDeviceSession(imei string) {
	imeiPrepare := util.PrepareImei(imei)

	v, ok := sm.sessions.LoadAndDelete(imeiPrepare)
	if !ok {
		return
	}

	sm.shutdownSession(v.(*types.TCPSession), imeiPrepare, true)
}

func (sm *SessionMemory) shutdownSession(session *types.TCPSession, imeiPrepare string, updateDeviceStatus bool) {
	session.Mu.Lock()

	if session.RequestId != nil {
		sm.requestMap.Delete(*session.RequestId)
		session.RequestId = nil
	}

	if session.InFlight != nil {
		close(session.InFlight.AnswerCh)
		session.InFlight = nil
	}

	if session.Device != nil {
		session.Device.OnDisconnect()
	}

	session.QueueCommand = nil
	session.QueueTelemetry = nil
	session.Busy = false
	session.Mu.Unlock()

	if err := sm.store.DeviceCommands.Update_DeviceCommandsResetOnConnectOnDisconnect(context.Background(), imeiPrepare); err != nil {
		logger.Error("shutdownSession imei={%s}: Failed to reset on_connect commands: %s", imeiPrepare, err.Error())
	}

	if err := sm.store.DeviceCommands.Update_DeviceCommandsMarkExecutionBatch(context.Background(), imeiPrepare, constants.CMD_STATUS_EXECUTIONERROR, "message.device-not-connected"); err != nil {
		logger.Error("shutdownSession imei={%s}: Failed batch mark execution error for pending/inprogress commands: %s", imeiPrepare, err.Error())
	}

	if updateDeviceStatus && session.Device != nil {
		sm.MarkCompanyTaskDisconnected(session.DeviceID)
	}

	if updateDeviceStatus {
		if err := sm.store.Devices.Update_DeviceStatusByImei(context.Background(), imeiPrepare, false); err != nil {
			logger.Error("shutdownSession imei={%s}: Failed to update device status true/false: %s", imeiPrepare, err.Error())
		}
	}
}

func (sm *SessionMemory) NewBundleSession(rabbitmqTransit types.RabbitMQ_TransitBinary) error {
	imeiPrepare := util.PrepareImei(rabbitmqTransit.Imei)

	v, ok := sm.sessions.Load(imeiPrepare)
	if !ok {
		return fmt.Errorf("the tracker is not connected to the server")
	}

	session := v.(*types.TCPSession)

	/* Check packet type */
	switch rabbitmqTransit.Type {
	/* Command -> command queue */
	case constants.ADM_RC_TYPE_STRING:
		sm.EnqueueCommand(rabbitmqTransit)
		return nil

	/* Set configuration */
	case constants.ADM_RC_TYPE_SET_CFG:
		session.ConfigurationMu.Lock()
		defer session.ConfigurationMu.Unlock()

		task, err := sm.store.Companies.Get_ActiveCompanyTaskByDeviceId(context.Background(), session.Device.DeviceId)
		if err != nil {
			return err
		}

		if task != nil {
			return errors.New("bulk configuration task is active")
		}

		if !sm.allowSetCfg(session, imeiPrepare) {
			return nil
		}

		if err := sm.TransitData(rabbitmqTransit.Imei, rabbitmqTransit.Type, rabbitmqTransit.Data); err != nil {
			return err
		}

		sm.MarkConfigurationPushed(session.Device.DeviceId, rabbitmqTransit.Data, 0, imeiPrepare)
		return nil

	/* Get configuration with request id */
	case constants.ADM_RC_TYPE_GET_CFG:
		session.Mu.Lock()
		if session.RequestId != nil {
			session.Mu.Unlock()
			return fmt.Errorf("device %s is busy", rabbitmqTransit.Imei)
		}

		/* Mark that this IMEI now has an outstanding request */
		session.RequestId = &rabbitmqTransit.RequestID
		session.Mu.Unlock()

		sm.requestMap.Store(rabbitmqTransit.RequestID, imeiPrepare)
		/* Transit data to terminal */
		sm.TransitData(rabbitmqTransit.Imei, rabbitmqTransit.Type, rabbitmqTransit.Data)

		/* Time-out request id */
		time.AfterFunc(constants.RabbitMQ_TimeoutRead, func() {
			sm.RemoveRequestIdSession(rabbitmqTransit.RequestID)
			logger.Warning("Request timeout req={%s} imei={%s}", rabbitmqTransit.RequestID, rabbitmqTransit.Imei)
		})

		return nil
	}

	return nil
}

/* RemoveRequestIdSession удалить requestId у сессии */
func (sm *SessionMemory) RemoveRequestIdSession(requestId string) {
	if imei, ok := sm.requestMap.Load(requestId); ok {
		if session, ok := sm.sessions.Load(imei.(string)); ok {
			sess := session.(*types.TCPSession)

			sess.Mu.Lock()
			if sess.RequestId != nil && *sess.RequestId == requestId {
				sess.RequestId = nil
			}
			sess.Mu.Unlock()
		}

		sm.requestMap.Delete(requestId)
	}
}

/* EnqueueCommand: Add command to queue */
func (sm *SessionMemory) EnqueueCommand(cmd types.RabbitMQ_TransitBinary) {
	imeiPrepare := util.PrepareImei(cmd.Imei)

	v, ok := sm.sessions.Load(imeiPrepare)
	if !ok {
		return
	}

	session := v.(*types.TCPSession)

	session.Mu.Lock()
	session.QueueCommand = append(session.QueueCommand, cmd)

	if session.Busy {
		session.Mu.Unlock()
		return
	}

	session.Busy = true
	session.Mu.Unlock()

	go sm.runQueue(session)
}

func (sm *SessionMemory) ScheduleTelemetryCommands(imei string, commands []string) {
	if len(commands) == 0 {
		return
	}

	imeiPrepare := util.PrepareImei(imei)

	v, ok := sm.sessions.Load(imeiPrepare)
	if !ok {
		return
	}

	session := v.(*types.TCPSession)

	session.Mu.Lock()
	if len(session.QueueTelemetry) > 0 {
		session.Mu.Unlock()
		return
	}

	if session.InFlight != nil && session.InFlight.Cmd.Telemetry {
		session.Mu.Unlock()
		return
	}

	for _, command := range commands {
		session.QueueTelemetry = append(session.QueueTelemetry, types.RabbitMQ_TransitBinary{
			Type:      constants.ADM_RC_TYPE_STRING,
			Imei:      imei,
			NResp:     true,
			Telemetry: true,
			Data:      []byte(command),
		})
	}

	if session.Busy {
		session.Mu.Unlock()
		return
	}

	session.Busy = true
	session.Mu.Unlock()

	go sm.runQueue(session)
}

func (sm *SessionMemory) DeliverCommandAnswer(imei, answer string) {
	session, ok := sm.GetDeviceSession(imei)
	if !ok || session == nil {
		return
	}

	session.Mu.Lock()
	inFlight := session.InFlight
	session.Mu.Unlock()

	if inFlight == nil {
		logger.Warning("DeliverCommandAnswer imei={%s}: dropped unsolicited answer", util.PrepareImei(imei))
		return
	}

	select {
	case inFlight.AnswerCh <- answer:
	default:
		logger.Warning("DeliverCommandAnswer imei={%s}: dropped answer for execution_id={%d}", util.PrepareImei(imei), inFlight.Cmd.ExecutionID)
	}
}

func (sm *SessionMemory) HandleCommandAnswer(answer string, cmd types.RabbitMQ_TransitBinary) {
	if cmd.Telemetry {
		return
	}

	ctx := context.Background()

	if cmd.NResp {
		sm.rmq.SendToRabbitAsync(types.RabbitMQ_TransitBinary{
			Type:      0x00,
			RequestID: cmd.RequestID,
			Imei:      cmd.Imei,
			StatCode:  http.StatusOK,
			Data:      []byte(answer),
		})
		return
	}

	if cmd.ExecutionID == 0 {
		logger.Warning("HandleCommandAnswer imei={%s}: missing execution_id for command={%s}", util.PrepareImei(cmd.Imei), string(cmd.Data))
		return
	}

	sm.LimiterDb <- struct{}{}
	if err := sm.store.DeviceCommands.Update_DeviceCommandExecutionByID(ctx, cmd.ExecutionID, constants.CMD_STATUS_COMPLETED, answer); err != nil {
		logger.Error("HandleCommandAnswer imei={%s} execution_id={%d}: Failed to update response: %s", util.PrepareImei(cmd.Imei), cmd.ExecutionID, err.Error())
	}
	<-sm.LimiterDb
}

func (sm *SessionMemory) prepareCommandExecution(ctx context.Context, cmd *types.RabbitMQ_TransitBinary) error {
	imei := util.PrepareImei(cmd.Imei)

	if cmd.NResp {
		return nil
	}

	if cmd.ExecutionID > 0 {
		return sm.store.DeviceCommands.Update_DeviceCommandExecutionInProgressByID(ctx, cmd.ExecutionID)
	}

	executionID, _, err := sm.store.DeviceCommands.Claim_DeviceCommandExecution(ctx, imei, string(cmd.Data))
	if err != nil {
		return err
	}

	cmd.ExecutionID = executionID
	return nil
}

func (sm *SessionMemory) runQueue(session *types.TCPSession) {
	timeout := constants.RabbitMQ_TimeoutRead

	for {
		if session.Device.Imei != nil {
			if _, ok := sm.sessions.Load(util.PrepareImei(*session.Device.Imei)); !ok {
				return
			}
		}

		session.Mu.Lock()
		var cmd types.RabbitMQ_TransitBinary
		if len(session.QueueCommand) > 0 {
			cmd = session.QueueCommand[0]
			session.QueueCommand = session.QueueCommand[1:]
		} else if len(session.QueueTelemetry) > 0 {
			cmd = session.QueueTelemetry[0]
			session.QueueTelemetry = session.QueueTelemetry[1:]
		} else {
			session.Busy = false
			session.Mu.Unlock()
			return
		}

		inFlight := &types.InFlightCommand{
			Cmd:      cmd,
			AnswerCh: make(chan string, 1),
		}
		session.InFlight = inFlight
		session.Mu.Unlock()

		ctx := context.Background()
		sm.LimiterDb <- struct{}{}
		prepareErr := sm.prepareCommandExecution(ctx, &inFlight.Cmd)
		<-sm.LimiterDb

		if prepareErr != nil {
			if !errors.Is(prepareErr, sql.ErrNoRows) {
				logger.Error("runQueue imei={%s}: Failed to prepare command execution: %s", cmd.Imei, prepareErr.Error())
			}

			session.Mu.Lock()
			session.InFlight = nil
			close(inFlight.AnswerCh)
			session.Mu.Unlock()
			continue
		}

		sm.TransitData(cmd.Imei, cmd.Type, cmd.Data)

		timer := time.NewTimer(timeout)
		select {
		case answer, ok := <-inFlight.AnswerCh:
			timer.Stop()
			if !ok {
				return
			}

			sm.HandleCommandAnswer(answer, inFlight.Cmd)

		case <-timer.C:
			if session.Device != nil {
				session.Device.ReleaseTerminal()
			}

			if inFlight.Cmd.ExecutionID > 0 && !inFlight.Cmd.NResp {
				sm.LimiterDb <- struct{}{}
				if err := sm.store.DeviceCommands.Update_DeviceCommandExecutionTimeoutByID(ctx, inFlight.Cmd.ExecutionID); err != nil {
					logger.Error("runQueue imei={%s} execution_id={%d}: Failed to handle command timeout: %s", cmd.Imei, inFlight.Cmd.ExecutionID, err.Error())
				}
				<-sm.LimiterDb
			}

			select {
			case <-inFlight.AnswerCh:
				logger.Warning("runQueue imei={%s} execution_id={%d}: discarded late answer after timeout", cmd.Imei, inFlight.Cmd.ExecutionID)
			default:
			}
		}

		session.Mu.Lock()
		session.InFlight = nil
		close(inFlight.AnswerCh)
		session.Mu.Unlock()
	}
}

func (sm *SessionMemory) allowSetCfg(session *types.TCPSession, imeiPrepare string) bool {
	if session == nil || session.Device == nil {
		return false
	}

	if util.HasDeviceOwner(session.Device.UserUuid) {
		return true
	}

	device, err := sm.store.Devices.Get_DeviceByImei(context.Background(), imeiPrepare)
	if err != nil {
		logger.Error("allowSetCfg imei={%s}: Failed to get device for SET_CFG owner check: %s", imeiPrepare, err.Error())
		return false
	}

	if device == nil || !util.HasDeviceOwner(device.UserUUID) {
		logger.Info("allowSetCfg imei={%s}: skip SET_CFG, device has no owner", imeiPrepare)
		return false
	}

	session.Device.UserUuid = device.UserUUID
	return true
}

/* MarkConfigurationPushed фиксирует, что этот хеш уже ушёл на трекер */
func (sm *SessionMemory) MarkConfigurationPushed(deviceID uint64, cfgData []byte, fallbackHash uint32, imei string) {
	if deviceID == 0 {
		return
	}

	cfgHash := common.GetCfgHash(cfgData)
	if cfgHash == 0 {
		cfgHash = fallbackHash
	}

	if cfgHash == 0 {
		return
	}

	go func() {
		sm.LimiterDb <- struct{}{}
		defer func() { <-sm.LimiterDb }()

		if err := sm.store.Configurations.Update_ConfigurationPushedHashByDeviceId(context.Background(), deviceID, cfgHash); err != nil {
			logger.Error("MarkConfigurationPushed imei={%s}: Failed to update cfg_pushed_hash: %s", imei, err.Error())
		}
	}()
}

/* TransitData только отправляет пакет на трекер */
func (sm *SessionMemory) TransitData(imei string, type_p uint8, buffer []byte) error {
	imeiPrepare := util.PrepareImei(imei)

	if session, ok := sm.GetDeviceSession(imeiPrepare); ok {
		return session.Device.TransitData(type_p, buffer)
	}

	logger.Error("TransitData imei={%s}: Session not found", imei)
	return errors.New("device session not found")
}

/* MarkCompanyTaskDisconnected сохраняет срок ожидания ранее отправленной конфигурации. */
func (sm *SessionMemory) MarkCompanyTaskDisconnected(deviceID uint64) {
	ctx := context.Background()

	task, err := sm.store.Companies.Get_ActiveCompanyTaskByDeviceId(ctx, deviceID)
	if err != nil {
		logger.Error("MarkCompanyTaskDisconnected device_id={%d}: %s", deviceID, err.Error())
		return
	}

	if task == nil || task.Status == constants.COMPANY_TASK_STATUS_PENDING {
		return
	}

	_, err = sm.store.Companies.Update_CompanyTaskDelivery(ctx, task, &models.CompanyTaskDelivery{
		Status:               constants.COMPANY_TASK_STATUS_PENDING,
		ErrorMessage:         task.ErrorMessage,
		AttemptToken:         task.AttemptToken,
		ConfirmationDeadline: task.ConfirmationDeadline,
		NextAttemptAt:        task.NextAttemptAt,
		LeaseUntil:           task.LeaseUntil,
	})
	if err != nil {
		logger.Error("MarkCompanyTaskDisconnected device_id={%d}: %s", deviceID, err.Error())
	}
}
