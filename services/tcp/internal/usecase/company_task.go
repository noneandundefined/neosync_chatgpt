package usecase

import (
	"context"
	"database/sql"

	"neomatica/neosync-tcp/config"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/memory"
	"neomatica/neosync-tcp/infra/store/postgres/models"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/pkg/protocol/adm"
	"neomatica/neosync-tcp/pkg/protocol/common"
	"neomatica/neosync-tcp/util"

	"github.com/google/uuid"
)

type CompanyTaskUseCase struct {
	Device *adm.ADMDevice
	Store  store.Storage
	Db     *sql.DB
}

func (u *CompanyTaskUseCase) ObserveDeviceConfiguration(ctx context.Context, deviceID uint64, cfgHash uint32) (bool, error) {
	task, err := u.Store.Companies.Get_ActiveCompanyTaskByDeviceId(ctx, deviceID)
	if err != nil || task == nil {
		return err != nil, err
	}

	if u.isTaskExpired(task) {
		return true, u.finishTask(ctx, task, constants.COMPANY_TASK_STATUS_EXPIRED, constants.COMPANY_TASK_ERROR_EXPIRED)
	}

	if err := u.Store.Companies.Update_CompanyTaskObservation(ctx, task.ID, cfgHash); err != nil {
		return true, err
	}

	current, err := u.Store.Companies.Get_ActiveCompanyTaskByDeviceId(ctx, deviceID)
	if err != nil || current == nil || current.ID != task.ID {
		return true, err
	}

	if cfgHash != 0 &&
		len(current.CfgData) > 0 &&
		cfgHash == common.GetCfgHash(current.CfgData) {
		return true, u.finishTask(ctx, current, constants.COMPANY_TASK_STATUS_COMPLETED, "")
	}

	return true, nil
}

/* TryApplyPendingTask сохраняет приоритет массовой задачи до её завершения. */
func (u *CompanyTaskUseCase) TryApplyPendingTask(ctx context.Context, session *memory.SessionMemory, d *models.Device_DeviceConf_Sync) (bool, error) {
	if d == nil || session == nil || u.Device == nil {
		return false, nil
	}

	currentSession, exists := session.GetDeviceSession(d.IMEI)
	if !exists || currentSession.Device != u.Device {
		return true, nil
	}

	task, err := u.Store.Companies.Get_ActiveCompanyTaskByDeviceId(ctx, d.ID)
	if err != nil || task == nil {
		return err != nil, err
	}

	if u.isTaskExpired(task) {
		return true, u.finishTask(
			ctx,
			task,
			constants.COMPANY_TASK_STATUS_EXPIRED,
			constants.COMPANY_TASK_ERROR_EXPIRED,
		)
	}

	if task.LeaseUntil != nil || task.ConfirmationDeadline != nil {
		if task.Status == constants.COMPANY_TASK_STATUS_PENDING {
			status := constants.COMPANY_TASK_STATUS_AWAITING_CONFIRMATION

			if task.LeaseUntil != nil {
				status = constants.COMPANY_TASK_STATUS_SENDING
			}

			_, err := u.Store.Companies.Update_CompanyTaskDelivery(
				ctx,
				task,
				&models.CompanyTaskDelivery{
					Status:               status,
					ErrorMessage:         task.ErrorMessage,
					AttemptToken:         task.AttemptToken,
					ConfirmationDeadline: task.ConfirmationDeadline,
					NextAttemptAt:        task.NextAttemptAt,
					LeaseUntil:           task.LeaseUntil,
				},
			)

			return true, err
		}

		return true, u.ProcessTaskDeadline(ctx, task)
	}

	if task.NextAttemptAt != nil && task.DatabaseNow.Before(*task.NextAttemptAt) {
		return true, nil
	}

	if !util.HasDeviceOwner(d.UserUUID) {
		return true, u.finishTask(ctx, task, constants.COMPANY_TASK_STATUS_FAILED, constants.COMPANY_TASK_ERROR_NO_ACCOUNT)
	}

	if len(task.CfgData) == 0 || common.GetCfgHash(task.CfgData) == 0 {
		return true, u.finishTask(ctx, task, constants.COMPANY_TASK_STATUS_FAILED, "invalid configuration payload")
	}

	if task.Attempts >= config.CompanyTaskMaxAttempts {
		return true, u.finishTask(ctx, task, constants.COMPANY_TASK_STATUS_FAILED, "configuration confirmation timeout")
	}

	return true, u.SendTaskConfiguration(ctx, session, d, task)
}

func (u *CompanyTaskUseCase) SendTaskConfiguration(ctx context.Context, session *memory.SessionMemory, d *models.Device_DeviceConf_Sync, task *models.CompanyTaskWithCompany) error {
	token := uuid.NewString()
	leaseUntil := task.DatabaseNow.Add(config.CompanyTaskLeaseTimeout)

	changed, err := u.Store.Companies.Update_CompanyTaskDelivery(ctx, task,
		&models.CompanyTaskDelivery{
			Status:       constants.COMPANY_TASK_STATUS_SENDING,
			AttemptToken: &token,
			LeaseUntil:   &leaseUntil,
			Claim:        true,
		},
	)

	if err != nil || !changed {
		return err
	}

	cfgHash := common.GetCfgHash(task.CfgData)

	current, err := u.Store.Companies.Get_ActiveCompanyTaskByDeviceId(ctx, d.ID)
	if err != nil {
		return err
	}

	if !ownsTaskAttempt(current, task.ID, token) {
		return nil
	}

	if u.isTaskExpired(current) {
		return u.finishTask(ctx, current, constants.COMPANY_TASK_STATUS_EXPIRED, constants.COMPANY_TASK_ERROR_EXPIRED)
	}

	imei := util.PrepareImei(d.IMEI)

	logger.Info("BulkConfiguration imei={%s}: task_id={%d} attempt={%d} pushing configuration bytes={%d}", imei, task.ID, current.Attempts, len(task.CfgData))

	sendErr := session.TransitData(imei, constants.ADM_RC_TYPE_SET_CFG, task.CfgData)

	if sendErr == nil {
		session.MarkConfigurationPushed(d.ID, task.CfgData, cfgHash, imei)

		d.CfgHash, d.CfgData, d.CfgSyncStatus = cfgHash, task.CfgData, constants.CFG_SYNC_STATUS_PENDING

		d.CfgPushedHash = cfgHash
	}

	return u.resolveSend(ctx, task, token, sendErr)
}

func ownsTaskAttempt(task *models.CompanyTaskWithCompany, taskID uint64, token string) bool {
	return task != nil && task.ID == taskID && (task.Status == constants.COMPANY_TASK_STATUS_SENDING || task.Status == constants.COMPANY_TASK_STATUS_PENDING) && task.AttemptToken != nil && *task.AttemptToken == token
}

func (u *CompanyTaskUseCase) resolveSend(ctx context.Context, task *models.CompanyTaskWithCompany, token string, sendErr error) error {
	current, err := u.Store.Companies.Get_ActiveCompanyTaskByDeviceId(ctx, task.DeviceID)

	if err != nil || !ownsTaskAttempt(current, task.ID, token) {
		return err
	}

	if u.isTaskExpired(current) {
		return u.finishTask(ctx, current, constants.COMPANY_TASK_STATUS_EXPIRED, constants.COMPANY_TASK_ERROR_EXPIRED)
	}

	if sendErr != nil {
		logger.Error("BulkConfiguration task_id={%d}: send failed: %s", task.ID, sendErr.Error())
		return u.retryTask(ctx, current, sendErr.Error())
	}

	deadline := current.DatabaseNow.Add(config.CompanyTaskConfirmationTimeout)

	_, err = u.Store.Companies.Update_CompanyTaskDelivery(ctx, current,
		&models.CompanyTaskDelivery{
			Status:               constants.COMPANY_TASK_STATUS_AWAITING_CONFIRMATION,
			ConfirmationDeadline: &deadline,
			Sent:                 true,
		},
	)

	return err
}

func (u *CompanyTaskUseCase) ProcessTaskDeadline(ctx context.Context, task *models.CompanyTaskWithCompany) error {
	if u.isTaskExpired(task) {
		return u.finishTask(ctx, task, constants.COMPANY_TASK_STATUS_EXPIRED, constants.COMPANY_TASK_ERROR_EXPIRED)
	}

	now := task.DatabaseNow

	if task.LeaseUntil != nil && !now.Before(*task.LeaseUntil) {
		return u.retryTask(ctx, task, "configuration delivery interrupted")
	}

	if task.ConfirmationDeadline != nil &&
		!now.Before(*task.ConfirmationDeadline) {
		return u.retryTask(ctx, task, "configuration confirmation timeout")
	}

	return nil
}

func (u *CompanyTaskUseCase) retryTask(ctx context.Context, task *models.CompanyTaskWithCompany, message string) error {
	if task.Attempts >= config.CompanyTaskMaxAttempts {
		return u.finishTask(ctx, task, constants.COMPANY_TASK_STATUS_FAILED, message)
	}

	status := constants.COMPANY_TASK_STATUS_QUEUED

	if task.Status == constants.COMPANY_TASK_STATUS_PENDING {
		status = constants.COMPANY_TASK_STATUS_PENDING
	}

	nextAttempt := task.DatabaseNow.Add(config.CompanyTaskRetryDelay)

	_, err := u.Store.Companies.Update_CompanyTaskDelivery(ctx, task,
		&models.CompanyTaskDelivery{
			Status:        status,
			ErrorMessage:  &message,
			NextAttemptAt: &nextAttempt,
		},
	)

	return err
}

func (u *CompanyTaskUseCase) isTaskExpired(task *models.CompanyTaskWithCompany) bool {
	return !task.DatabaseNow.Before(task.ExpiresAt)
}

func (u *CompanyTaskUseCase) finishTask(ctx context.Context, task *models.CompanyTaskWithCompany, status, message string) error {
	var errMsg *string

	if message != "" {
		errMsg = &message
	}

	_, err := u.Store.Companies.Update_CompanyTaskDelivery(ctx, task,
		&models.CompanyTaskDelivery{
			Status:       status,
			ErrorMessage: errMsg,
		},
	)

	return err
}
