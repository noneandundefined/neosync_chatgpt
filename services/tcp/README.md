# NeoSync TCP service

TCP-сервис принимает соединения терминалов, разбирает ADM protocol, поддерживает online-сессии и выполняет команды/синхронизацию.

```text
cmd/server/          listener, event wiring, handlers
config/              конфигурация
emulators/           локальные ADM эмуляторы
internal/usecase/    application flows
infra/store/         PostgreSQL, Redis, memory
infra/messaging/     RabbitMQ
pkg/protocol/        ADM protocol codecs/parsers
types/               transport/session types
util/                небольшие helpers
```

RabbitMQ является infrastructure adapter и не хранится в `cmd`. Use cases также не являются частью PostgreSQL.

```bash
make test
make lint
make build-linux
```
