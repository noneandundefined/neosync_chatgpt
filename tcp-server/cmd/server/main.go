package main

import (
	"database/sql"
	"fmt"
	"neomatica/neosync-tcp/cmd/rabbitmq"
	"neomatica/neosync-tcp/cmd/server/handlers"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/infra/store/memory"
	"neomatica/neosync-tcp/infra/store/postgres"
	"neomatica/neosync-tcp/infra/store/postgres/store"
	"neomatica/neosync-tcp/infra/store/redis"
	"neomatica/neosync-tcp/types"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var (
	Version          = "1.1.x"
	maxHelloQueues   = 10000
	maxCommandQueues = 15000
	maxPool          = 255
)

type tcpServer struct {
	db      *sql.DB
	cron    *cron.Cron
	store   store.Storage
	cache   *memory.Cache
	rmq     *rabbitmq.RabbitMQ
	session *memory.SessionMemory
	handler *handlers.BasePackEventHandler
}

func main() {
	/* .env - .env.production */
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	/* Initial logger */
	logger.InitLogger()

	/* Initial connection db */
	db, err := postgres.New(os.Getenv("DB_ADDR"), 1000, 80, "7m")
	if err != nil {
		logger.Error("Error connect to database: %s", err.Error())
		return
	}
	defer db.Close()

	/* Initial redis */
	if err := redis.NewRedisDb(); err != nil {
		logger.Error("Error connect to redis: %s", err.Error())
		return
	}

	/* Initial RabbitMQ */
	rmq, err := rabbitmq.NewRabbitMQ(fmt.Sprintf("amqp://%s:%s@%s/", os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST")))
	if err != nil {
		logger.Error("Error connect to RabbitMQ: %s", err.Error())
		return
	}
	defer rmq.Close()

	/* Create limiter DB/Redis */
	limiterDb := make(chan struct{}, maxPool)

	/* Store for postgres */
	store := store.NewStorage(db)
	/* Memory session */
	session := memory.NewSessionMemory(store, rmq, limiterDb)
	/* Memory cache */
	cache := memory.NewCache()

	go func() {
		_ = rmq.Consume(func(body []byte) {
			rabbitmqTransit, err := rmq.Decoder(body)
			if err != nil {
				logger.Error("Failed decode: %s", err.Error())
				return
			}

			if err := session.NewBundleSession(rabbitmqTransit); err != nil {
				logger.Error("Failed bundle session req={%s} imei={%s}: %s", rabbitmqTransit.RequestID, rabbitmqTransit.Imei, err.Error())

				/* RabbitMQ Error */
				_ = rmq.SendToRabbitAsync(types.RabbitMQ_TransitBinary{
					Type:      0x00,
					RequestID: rabbitmqTransit.RequestID,
					Imei:      rabbitmqTransit.Imei,
					StatCode:  http.StatusBadRequest,
					Data:      []byte(err.Error()),
				})

				return
			}
		})
	}()

	/* Initial base structure */
	server := &tcpServer{
		db:      db,
		store:   store,
		rmq:     rmq,
		cache:   cache,
		session: session,
	}
	server.cron = cron.New()

	/* Handlers inside tcpServer structure */
	handler := &handlers.BasePackEventHandler{
		Session: session,
		RMQ:     rmq,
		Cache:   cache,
		Store:   store,
		/* Limiter */
		LimiterDb: limiterDb,
	}
	handler.QueueHelloEvent = make(chan handlers.JobHelloEvent, maxHelloQueues)
	handler.QueueCommandEvent = make(chan handlers.JobCommandEvent, maxCommandQueues)

	/* Initial hello events workers */
	handler.HelloPackEventWorkers(maxPool)
	/* Initial command events workers */
	handler.CommandPackEventWorkers(maxPool)

	server.handler = handler

	/* Cron */
	server.startCronStatusNotActive()
	server.startCronCompanyTasks()
	/* Clear old logs, every day in 00:00 */
	server.startCronCleanupLogs()

	server.cron.Start()

	/* Start TCP server */
	if err := server.tcpStart(); err != nil {
		logger.Error("Failed start TCP server: %s", err.Error())
	}
}
