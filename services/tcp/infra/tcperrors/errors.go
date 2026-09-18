package tcperrors

import (
	"context"
	"errors"
	"fmt"
	"neomatica/neosync-tcp/infra/logger"
	"net"

	"github.com/lib/pq"
)

var (
	Err_SendPacketToDevice       = errors.New("ошибка отправки пакета данных на трекер")
	Err_DeviceNotConnected       = errors.New("трекер не подключен к серверу")
	Err_UnknownTypeOfCommand     = errors.New("неправильный или неизвестный тип команды")
	Err_FirmwareOutdatedNotFound = errors.New("сервер не поддерживает текущую версию или прошивку терминала")
	Err_PacketIsNull             = errors.New("трекер отпарвил пустой пакет")
	Err_SendCommandToDevice      = errors.New("ошибка при отправке команды на устройство")
	Err_IncorrectPackageSize     = errors.New("некорректный размер пакета")
	Err_InternalServer           = errors.New("ошибка на стороне сервера")
	Err_ContextDeadlineExceeded  = errors.New("долгое время ожидания ответа от базы данных, повторите попытку позже")
	Err_ContextCanceled          = errors.New("операция была отменена")
	Err_UniqueViolation          = errors.New("введенные данные уже существует")
	Err_DbTimeout                = errors.New("соединение с базой данных превысило время ожидания")
	Err_DbNetworkTemporary       = errors.New("временная проблема с сетью. попробуйте ещё раз")
	Err_DbNetwork                = errors.New("сетевая ошибка при подключении к базе данных")
	Err_NotDeleted               = errors.New("ошибка, ни одна запись не была удалена")
	Err_NotUpdated               = errors.New("ошибка, ни одна запись не была обновлена")
)

func HandleSQLError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		logger.Warning("HandleSQLError: %s", err.Error())
		return Err_ContextDeadlineExceeded
	}

	if errors.Is(err, context.Canceled) {
		logger.Warning("HandleSQLError: %s", err.Error())
		return Err_ContextCanceled
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			logger.Warning("HandleSQLError: %s", err.Error())
			return Err_DbTimeout
		}

		logger.Error("HandleSQLError: %s", err.Error())
		return Err_DbNetwork
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": /* Error unnique */
			logger.Warning("HandleSQLError: %s", err.Error())
			return Err_UniqueViolation
		case "57014": /* query_canceled */
			logger.Warning("HandleSQLError: %s", err.Error())
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return Err_ContextDeadlineExceeded
			}
			return Err_ContextCanceled
		}
	}

	return fmt.Errorf("не удалось выполнить операцию с базой данных")
}
