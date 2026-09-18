# NeoSync HTTP service

REST API NeoSync.

```text
cmd/server/          запуск и composition root
config/              конфигурация
handler/v1/          HTTP transport
internal/usecase/    application flows
infra/store/         PostgreSQL, Redis, memory
infra/messaging/     RabbitMQ
middleware/          HTTP middleware
migrations/          Goose SQL migrations
pkg/                 reusable packages
tests/               tests
```

Use case не должен находиться внутри PostgreSQL-пакета: transport → usecase → infrastructure.

```bash
make test
make lint
make build-linux
make migrate-up
```

Swagger: `/api/v1/doc.api/swagger/gui`.
