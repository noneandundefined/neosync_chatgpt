package worker

import (
	"context"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/postgres/models"
	"neomatica/neosync/infra/store/postgres/store"
	"neomatica/neosync/pkg/cfgmerge"
	"sync"
)

type CompanyTasksJob struct {
	CompanyID uint64
	DeviceIDs []uint64
	UserUUID  string
}

type CompanyTasksWorker struct {
	store  store.Storage
	mu     sync.Mutex
	queues map[string]chan CompanyTasksJob
}

func NewCompanyTasksWorker(storage store.Storage) *CompanyTasksWorker {
	return &CompanyTasksWorker{
		store:  storage,
		queues: make(map[string]chan CompanyTasksJob),
	}
}

func (w *CompanyTasksWorker) Enqueue(job CompanyTasksJob) {
	w.mu.Lock()

	queue, exists := w.queues[job.UserUUID]
	if !exists {
		queue = make(chan CompanyTasksJob, 64)
		w.queues[job.UserUUID] = queue

		go w.runUserQueue(job.UserUUID, queue)
	}

	w.mu.Unlock()

	queue <- job
}

func (w *CompanyTasksWorker) runUserQueue(userUUID string, queue <-chan CompanyTasksJob) {
	for job := range queue {
		w.processJob(userUUID, job)
	}
}

func (w *CompanyTasksWorker) processJob(userUUID string, job CompanyTasksJob) {
	ctx := context.Background()
	tr := locale.NewTranslator("en")

	company, err := w.store.Companies.Get_CompanyForProcessing(ctx, job.CompanyID, userUUID)
	if err != nil {
		logger.Error("CompanyTasksWorker user={%s} company_id={%d}: failed to load company: %s", userUUID, job.CompanyID, err.Error())
		return
	}

	if company == nil {
		logger.Error("CompanyTasksWorker user={%s} company_id={%d}: company not found", userUUID, job.CompanyID)
		return
	}

	createdTasks := 0

	for _, deviceID := range job.DeviceIDs {
		var deviceConfiguration *models.Configuration

		if company.ConfigurationSource == "template_sources" && cfgmerge.IsTemplateChangesData(company.ConfigurationSourceData) {
			deviceConfiguration, err = w.store.Configurations.Get_ConfigurationByDeviceId(ctx, deviceID)
			if err != nil {
				logger.Error("CompanyTasksWorker user={%s} company_id={%d} device_id={%d}: failed to load device configuration: %s", userUUID, job.CompanyID, deviceID, err.Error())
				continue
			}
		}

		cfgData, err := cfgmerge.BuildTaskCfgData(tr, company.ConfigurationSource, company.ConfigurationSourceData, deviceConfiguration)
		if err != nil {
			logger.Error("CompanyTasksWorker user={%s} company_id={%d} device_id={%d}: failed to build cfg_data: %s", userUUID, job.CompanyID, deviceID, err.Error())
			continue
		}

		if err := w.store.Companies.Create_CompanyTask(ctx, job.CompanyID, deviceID, cfgData); err != nil {
			logger.Error("CompanyTasksWorker user={%s} company_id={%d} device_id={%d}: failed to create company task: %s", userUUID, job.CompanyID, deviceID, err.Error())
			continue
		}

		createdTasks++
	}

	if createdTasks == 0 {
		if err := w.store.Companies.Update_CompanyStatus(ctx, job.CompanyID, "failed"); err != nil {
			logger.Error("CompanyTasksWorker user={%s} company_id={%d}: failed to mark company failed: %s", userUUID, job.CompanyID, err.Error())
		}

		return
	}

	if err := w.store.Companies.Refresh_CompanyStatus(ctx, job.CompanyID); err != nil {
		logger.Error("CompanyTasksWorker user={%s} company_id={%d}: failed to update company status: %s", userUUID, job.CompanyID, err.Error())
	}
}
