package store

import (
	"context"
	"database/sql"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/pkg/adm"
	"neomatica/neosync/types"
)

type Storage struct {
	Users interface { //nolint
		Create_UserCore(ctx context.Context, tx *sql.Tx, user *models.UserCore) error
		Create_UserContact(ctx context.Context, tx *sql.Tx, user *models.UserContact) error
		Create_UserRole(ctx context.Context, tx *sql.Tx, user *models.UserRole) error
		Create_UserAccess(ctx context.Context, tx *sql.Tx, user *models.UserAccess) error

		Get_ColumnTableUserCore(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_ColumnTableUserRole(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_ColumnTableUserContacts(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_Users(ctx context.Context) ([]models.UserToOwner, error)
		Get_UserAuthByUuid(ctx context.Context, userUuid string) (*models.UserAuth, error)
		Get_UsersWithParams(ctx context.Context, limit, offset int, search, role, columnSortKey, columnSortDir, userUuid string) ([]models.User, int, int, error)
		Get_UsersToOwner(ctx context.Context, role, userUuid string, parentUuid sql.NullString, search string, limit, offset int, all bool) ([]models.UserToOwner, int, error)
		Get_UsersByParentUuid(ctx context.Context, userUuid string) ([]models.User, error)
		Get_UsersForTree(ctx context.Context, roleCode, viewerUUID string, dealerUUID *string, search string, limit, offset int, all bool) ([]models.UserTree, int, error)
		Get_UsersByParentUuidWithParams(ctx context.Context, limit, offset int, search, role, columnSortKey, columnSortDir, dealerUuid, currentUserUuid string) ([]models.User, int, int, error)
		Get_UsersByParentUuidToOwner(ctx context.Context, userUuid string) ([]models.UserToOwner, error)
		Get_UserByUuid(ctx context.Context, userUuid string) (*models.UserAuth, error)
		Get_UserLoginWithRoleCodeByUuid(ctx context.Context, userUuid string) (*models.UserLoginWithRoleCode, error)
		Get_UserWithPasswordByUuid(ctx context.Context, userUuid string) (*models.UserAuth, error)
		Get_UserCoreByUuid(ctx context.Context, userUuid string) (*models.UserCore, error)
		Get_UserCoreByRefreshToken(ctx context.Context, token string) (*models.UserCore, error)
		Get_UserAccessByUuid(ctx context.Context, userUuid string) (*models.UserAccess, error)
		Get_UserCoreByEmail(ctx context.Context, email string) (*models.UserCore, error)
		Get_UserRoleByUuid(ctx context.Context, userUuid string) (*models.UserRole, error)

		Update_UserCoreRefreshToken(ctx context.Context, token, userUuid string) error
		Update_UserPasswordByUuid(ctx context.Context, password, userUuid string) error
		Update_UserContactOnesByUuid(ctx context.Context, tx *sql.Tx, userUuid string, user *models.UserUpdate) error
		Update_UserAccessOnesByUuid(ctx context.Context, tx *sql.Tx, userUuid string, user *models.UserUpdate) error
		Update_UserCoreOnesByUuid(ctx context.Context, tx *sql.Tx, userUuid string, user *models.UserUpdate) error
		Update_CanViewChildGroupsByUuid(ctx context.Context, userUuid string, canViewChildGroups bool) error
		Update_ForceNeosyncConfigurationPriority(ctx context.Context, userUuid string, enabled bool) error

		Delete_UsersByUuid(ctx context.Context, userUuids []string) error
		Delete_UserByUuid(ctx context.Context, userUuid string) error
	}
	Devices interface { //nolint
		Create_Device(ctx context.Context, tx *sql.Tx, device *models.Device) (uint64, error)
		Create_DevicesBatch(ctx context.Context, tx *sql.Tx, devices []*models.Device) ([]*models.Device, error)
		Create_DeviceConf(ctx context.Context, tx *sql.Tx, device *models.DeviceConf) error
		Create_DeviceConfBatch(ctx context.Context, tx *sql.Tx, devices []*models.DeviceConf) error

		Get_Devices(ctx context.Context) ([]models.DeviceShort, error)
		Get_ColumnTableDevice(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_DeviceConfByDeviceId(ctx context.Context, deviceId uint64) (*models.DeviceConf, error)
		Get_DevicesWithParams(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, userUuid string) ([]models.Device_DeviceConf_Sync, int, int, error)
		Get_DevicesByUuidsWithParams(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, roleCode string, userUuid, ownerUuid *string) ([]models.Device_DeviceConf_Sync, int, int, error)
		// Get_DevicesByUuid(ctx context.Context, authToken *types.AuthToken) ([]models.Device_DeviceConf_Sync, error)
		// Get_DevicesByAuthTokenWithParams(ctx context.Context, limit, offset int, search string, authToken *types.AuthToken) ([]models.Device_DeviceConf_Sync, int, error)
		Get_DeviceStatusByImei(ctx context.Context, imei string) (*models.DeviceStatus, error)
		Get_DeviceDetailedStatusByImeis(ctx context.Context, imeis []string) ([]models.DeviceDetailedStatus, error)
		Get_DevicesByImeis(ctx context.Context, imeis []string) ([]models.Device, error)
		Get_DeviceByImei(ctx context.Context, imei string) (*models.Device, error)
		Get_DeviceAndConfByImei(ctx context.Context, imei string) (*models.Device_DeviceConf_Sync, error)

		Update_DeviceStatusByImei(ctx context.Context, imei string, status bool) error
		Update_DeviceOwnerUuidByImei(ctx context.Context, imei string, ownerUuid *string) error
		Update_DeviceUuidByImei(ctx context.Context, imei string, userUuid, ownerUuid *string) error
		Update_DeviceUuidByDeviceId(ctx context.Context, deviceId uint64, userUuid *string) error
		Update_DeviceUuidByDeviceIdTx(ctx context.Context, tx *sql.Tx, deviceId uint64, userUuid, ownerUuid *string) error
		Update_DeviceByImei(ctx context.Context, tx *sql.Tx, imei string, device *models.UpdateDevice) error
		Update_DevicesUuidByImeiToNull(ctx context.Context, imeis []string) error
		Update_DeviceUnlinkFromAccountByImei(ctx context.Context, imei, userUuid string) error
		Update_DevicesUnlinkFromAccountByImei(ctx context.Context, imeis []string, userUuid string) error
		Update_DeviceConfByDeviceId(ctx context.Context, tx *sql.Tx, deviceId uint64, device *models.UpdateDevice) error

		Delete_DevicesByImei(ctx context.Context, imeis []string) error
		Delete_DevicesByImeiAndUserUuid(ctx context.Context, imeis []string, userUuid string) error
		Delete_DeviceByImeiAndUserUuid(ctx context.Context, imei, userUuid string) error
		Delete_DeviceByImei(ctx context.Context, imei string) error
	}
	DeviceCommands interface { //nolint
		Create_DeviceCommand(ctx context.Context, tx *sql.Tx, command *models.DeviceCommand) (uint64, error)
		Create_Batch_DeviceCommandExecution(ctx context.Context, tx *sql.Tx, commandId uint64, imeis []string) error

		Get_DeviceCommandWithExecutionsById(ctx context.Context, userUuid string, commandId uint64) (*models.DeviceCommandWithExecutions, error)
		Get_DeviceCommandWithExecutionsByUuid(ctx context.Context, userUuid string) ([]models.DeviceCommand, error)
		Get_DeviceCommandWithExecutionsByUuidAndSessId(ctx context.Context, userUuid, sessId string) ([]models.DeviceCommandWithExecutions, error)
		Get_DevicesStateByUserUuid(ctx context.Context, authToken *types.AuthToken, search string) ([]models.DevicesForSendCommand, error)
		Get_DevicesState(ctx context.Context, userUuid string, search string) ([]models.DevicesForSendCommand, error)
		Get_DevicesStateDevicesByUserUuid(ctx context.Context, authToken *types.AuthToken, groupID uint64, search string, limit, offset int, all bool) ([]models.DeviceShort, int, error)
		Get_DevicesStateDevices(ctx context.Context, userUuid string, groupID uint64, search string, limit, offset int, all bool) ([]models.DeviceShort, int, error)

		Update_DeviceCommandStatusProgressByImeiAndCommandId(ctx context.Context, imeis []string, commandId uint64) error
		Update_DeviceCommandCancelByCommandId(ctx context.Context, userUuid string, commandId uint64) error

		Delete_DeviceCommandByUuidAndId(ctx context.Context, userUuid string, command_id uint64) error
	}
	Configurations interface { //nolint
		Create_Configuration(ctx context.Context, tx *sql.Tx, cfg *models.Configuration) error
		Create_ConfigurationsBatch(ctx context.Context, tx *sql.Tx, cfgs []*models.Configuration) error

		Get_ConfigurationByDeviceId(ctx context.Context, id uint64) (*models.Configuration, error)
		Get_ReferenceConfigurationByModelAndUserUuid(ctx context.Context, model, userUuid string) (*models.Configuration, error)
		Get_ColumnTableConfiguration(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_ConfigurationHistoriesByDeviceId(ctx context.Context, id uint64) ([]models.ConfigurationHistory, error)
		Get_ConfigurationHistoryByCfgHash(ctx context.Context, cfgHash uint32) (*models.ConfigurationHistory, error)
		Get_ColumnTableConfigurationProf(ctx context.Context) ([]models.DatabaseSchemaColumn, error)

		Replace_ConfigurationProf(ctx context.Context, deviceID uint64, cfgHash uint32, prof *adm.ConfigurationProfValues) error
		Update_ConfigurationByDeviceId(ctx context.Context, id uint64, cfg []byte) error
		Update_ConfigurationSyncStatusByDeviceId(ctx context.Context, deviceID uint64, status string, syncError *string) error
	}
	Syncs interface { //nolint
		Create_Sync(ctx context.Context, tx *sql.Tx, sync *models.Sync) error
		Create_SyncsBatch(ctx context.Context, tx *sql.Tx, syncs []*models.Sync) error

		Get_ColumnTableSync(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_SyncByDeviceId(ctx context.Context, id uint64) (*models.Sync, error)

		Update_SyncByDeviceId(ctx context.Context, tx *sql.Tx, id uint64, cfg []byte) error
		Update_SyncFirmwareVersionUpd(ctx context.Context, firmware uint16, model string) error
		Update_FirmwareUpdateStatusByDeviceId(ctx context.Context, deviceID uint64, status string) error
	}
	Groups interface { //nolint
		Create_Group(ctx context.Context, group *models.Group) error

		Get_Groups(ctx context.Context) ([]models.GroupName, error)
		Get_ColumnTableGroup(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_GroupById(ctx context.Context, groupId uint64) (*models.Group, error)
		Get_GroupsByUuid(ctx context.Context, userUuid string) ([]models.Group, error)
		Get_GroupsByUuidWithParams(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, userUuid string) ([]models.Group, int, int, error)
		Get_GroupWithDevices(ctx context.Context, groupId uint64, authToken *types.AuthToken) (*models.GroupWithDevices, error)
		Get_Assigned_And_Available_Devices(ctx context.Context, groupId uint64, authToken *types.AuthToken) ([]models.GDevice, []models.GDevice, error)
		Get_GroupNamesByUuid(ctx context.Context, userUuid string) ([]models.GroupName, error)
		Get_ListDeviceIDsInGroup(ctx context.Context, groupID uint64) ([]uint64, error)

		Update_AddDevicesToGroup(ctx context.Context, device *models.GroupDevices) error
		Update_AddDevicesToGroupTx(ctx context.Context, tx *sql.Tx, device *models.GroupDevices) error
		Update_RemoveDevicesFromGroup(ctx context.Context, device *models.GroupDevices) error
		Update_GroupById(ctx context.Context, groupId uint64, group *models.Group) error

		Delete_GroupsById(ctx context.Context, groupIds []uint64) error
		Delete_GroupById(ctx context.Context, groupId uint64) error
	}
	GroupMembers interface { //nolint
		AggregatedDevicePermissions(ctx context.Context, memberUserUUID string, deviceID uint64) (canSendCommands, canReadConfig, canEditConfig bool, err error)

		Get_ColumnTableGroupMember(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_GroupMember(ctx context.Context, groupID uint64, memberUserUUID string) (*models.GroupMember, error)
		Get_GroupMembers(ctx context.Context, groupID uint64) ([]models.GroupMember, error)
		Get_UserHasGroupAccess(ctx context.Context, groupID uint64, userUUID string) (bool, error)

		Update_GroupShareDefaults(ctx context.Context, groupID uint64, canEditGroup, canManageDevices, canReadConfig, canEditConfig, canSendCommands bool) error
		Upsert_GroupMember(ctx context.Context, groupID uint64, memberUserUUID string) error
		Upsert_GroupMembers(ctx context.Context, groupID uint64, memberUUIDs []string) error

		Delete_GroupMember(ctx context.Context, groupID uint64, memberUserUUID string) error
		Delete_GroupMembers(ctx context.Context, groupID uint64, memberUUIDs []string) error
	}
	Companies interface { //nolint
		Create_Company(ctx context.Context, tx *sql.Tx, company *models.Company) (uint64, error)
		Create_CompanyTask(ctx context.Context, companyID, deviceID uint64, cfgData []byte) error
		Get_CompanyForProcessing(ctx context.Context, companyID uint64, userUUID string) (*models.Company, error)
		Get_PendingCompanyTaskDeviceIds(ctx context.Context, deviceIDs []uint64) ([]uint64, error)
		Get_CompaniesByUserUuid(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, userUuid string) ([]models.Company, int, int, error)
		Get_CompanyByIdAndUserUuid(ctx context.Context, companyID uint64, userUuid string) (*models.CompanyWithTasks, error)
		Refresh_CompanyStatus(ctx context.Context, companyID uint64) error
		Retry_FailedCompanyTasksByCompanyId(ctx context.Context, companyID uint64) (int64, error)
		Cancel_PendingCompanyTasksByCompanyId(ctx context.Context, companyID uint64) (int64, error)
		Cancel_PendingCompanyTasksByDeviceIds(ctx context.Context, tx *sql.Tx, deviceIDs []uint64) error
		Update_CompanyStatus(ctx context.Context, companyID uint64, status string) error
	}
	DeviceModelSources interface { //nolint
		Get_Sources(ctx context.Context) ([]models.DeviceModelSource, error)
		Get_Models(ctx context.Context) ([]string, error)
		Get_FirmwareUrl(ctx context.Context) ([]string, error)
	}
	Meta interface { //nolint
		Get_MaintenanceIsActive(ctx context.Context) (*models.Maintenance, error)
		Update_MaintenanceIsActive(ctx context.Context, isActive bool) error
	}
	ConfigurationTemplates interface { //nolint
		Create_ConfigurationTemplate(ctx context.Context, template *models.ConfigurationTemplate, cfgData []byte) (uint64, error)
		Get_ConfigurationTemplatesByUserUuid(ctx context.Context, limit, offset int, search, columnSortKey, columnSortDir, userUuid string) ([]models.ConfigurationTemplate, int, int, error)
		Get_ConfigurationTemplateByIdAndUserUuid(ctx context.Context, id uint64, userUuid string) (*models.ConfigurationTemplate, error)
		Get_ConfigurationTemplateCfgDataByIdAndUserUuid(ctx context.Context, id uint64, userUuid string) ([]byte, error)
		Delete_ConfigurationTemplateByIdAndUserUuid(ctx context.Context, id uint64, userUuid string) error
		Update_ConfigurationTemplateByIdAndUserUuid(ctx context.Context, id uint64, userUuid, name, typeOfSaving string, cfgData []byte) error
	}
	Analytics interface { //nolint
		Create_AnalyticsRecord(ctx context.Context, analytic *models.AnalyticRecord) error
		// Create_AnalyticsConfigProfile(ctx context.Context, profile *models.AnalyticsConfigProfile) error
		Create_AnalyticsUser(ctx context.Context, analytic *models.AnalyticsUser) error

		Get_ColumnTableAnalyticsUsers(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_ColumnTableAnalyticRecords(ctx context.Context) ([]models.DatabaseSchemaColumn, error)
		Get_ColumnTableProductAnalyticsEvents(ctx context.Context) ([]models.DatabaseSchemaColumn, error)

		Exec_SelectQuery(ctx context.Context, query string, args ...any) (*models.AnalyticQueryResult, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:                  &UserStore{db},
		Devices:                &DeviceStore{db},
		DeviceCommands:         &DeviceCommandStore{db},
		Configurations:         &ConfigurationStore{db},
		Syncs:                  &SyncStore{db},
		Groups:                 &GroupStore{db},
		GroupMembers:           &GroupMembersStore{db},
		Companies:              &CompanyStore{db},
		ConfigurationTemplates: &ConfigurationTemplateStore{db},
		DeviceModelSources:     &DeviceModelSources{db},
		Meta:                   &MetaStore{db},
		Analytics:              &AnalyticStore{db},
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
