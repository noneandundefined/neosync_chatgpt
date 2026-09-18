# NeoSync

NeoSync — монорепозиторий сервиса удалённой настройки терминалов Neomatica.

## С чего начать

Если вы впервые открыли проект, сначала прочитайте [архитектуру](./docs/architecture.md).

```text
apps/
  web/              React/Vite приложение

services/
  http/             REST API
  tcp/              TCP-сервер терминалов

deploy/
  nginx/            nginx-конфигурация
  postgres/         PostgreSQL-конфигурация
  systemd/          production systemd units

docs/               документация
scripts/            эксплуатационные/deploy-скрипты
docker-compose.yml  локальное окружение
```

Главное правило backend: transport вызывает `internal/usecase`, а внешние системы находятся в `infra`.

## Быстрый запуск

```bash
cp services/http/.env.example services/http/.env
cp services/tcp/.env.example services/tcp/.env
cp apps/web/.env.example apps/web/.env

docker compose up --build
```

После запуска:
- Web: http://localhost:5173
- HTTP API: http://localhost:8080
- TCP: localhost:12346
- RabbitMQ UI: http://localhost:15672

## Разработка

```bash
cd services/http && make test
cd services/tcp && make test
cd apps/web && pnpm install && pnpm build
```

Миграции принадлежат HTTP-сервису и лежат в `services/http/migrations`:

```bash
cd services/http
make migrate-up
```

## Где добавлять код

- REST endpoint → `services/http/handler/v1/<domain>_handler_v1`
- бизнес-операция → `services/http/internal/usecase` или `services/tcp/internal/usecase`
- PostgreSQL → `services/*/infra/store/postgres`
- Redis → `services/*/infra/store/redis`
- RabbitMQ → `services/*/infra/messaging/rabbitmq`
- TCP protocol → `services/tcp/pkg/protocol`
- frontend page → `apps/web/src/pages`
- frontend API client → `apps/web/src/rest`
- production config → `deploy`

## Деплой

CI вызывает скрипты из `scripts/`. Ручной deploy:

```bash
cd scripts
bash pre_deploy
bash systemd
bash setup_nginx
```
