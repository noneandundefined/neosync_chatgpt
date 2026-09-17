package main

import (
	"database/sql"
	"fmt"
	"neomatica/neosync/cmd/rabbitmq"
	"neomatica/neosync/infra/analytics"
	"neomatica/neosync/infra/locale"
	"neomatica/neosync/infra/logger"
	"neomatica/neosync/infra/store/memory"
	"neomatica/neosync/infra/store/postgres"
	"neomatica/neosync/infra/store/postgres/store"
	"neomatica/neosync/infra/store/postgres/usecase"
	"neomatica/neosync/infra/store/redis"
	"neomatica/neosync/pkg/clientip"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

type httpServer struct {
	db        *sql.DB
	cron      *cron.Cron
	store     store.Storage
	usecase   usecase.UseCase
	rmq       *rabbitmq.RabbitMQ
	session   *memory.SessionService
	analytics *analytics.Collector
}

func main() {
	/* .env - .env.production */
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	if err := clientip.Configure(os.Getenv("TRUSTED_PROXIES")); err != nil {
		panic(err)
	}

	/* Inital logger */
	logger.InitLogger()

	/* Localization en/es/ru */
	locale.InitI18n()

	/* Initial connect db */
	db, err := postgres.New(os.Getenv("DB_ADDR"), 150, 25, "7m")
	if err != nil {
		logger.Error("Failed connect to database: %s", err.Error())
		return
	}
	defer db.Close()

	/* Store for postgres */
	store := store.NewStorage(db)

	/* Background product analytics collector */
	analyticsCollector := analytics.NewCollector(db)
	defer analyticsCollector.Close()

	/* Usecase for postgres */
	usecase := usecase.NewUseCase(db, store)

	/* Memory session */
	session := memory.NewSessionService()

	/* Initial connect redis */
	if err := redis.NewRedisDb(); err != nil {
		logger.Error("Failed connect to redis: %s", err.Error())
		return
	}

	/* RabbitMQ actions */
	rmq, err := rabbitmq.NewRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s/", os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST")))
	if err != nil {
		logger.Error("Failed connect to RabbitMQ: %s", err.Error())
		return
	}
	defer rmq.Close()

	_ = rmq.Consume(func(body []byte) {
		session.Dispatch(body)
	})

	server := &httpServer{
		db:        db,
		store:     store,
		usecase:   usecase,
		rmq:       rmq,
		session:   session,
		analytics: analyticsCollector,
	}
	server.cron = cron.New()

	/* Terminal firmware updates, every day at 00:00 */
	server.startCronFirmwareUpdater()
	/* Clear old logs, every day in 00:00 */
	server.startCronCleanupLogs()

	server.cron.Start()

	/* Started HTTPx server */
	if err := server.httpStart(); err != nil {
		logger.Error("Failed start TCP server: %s", err.Error())
	}
}
