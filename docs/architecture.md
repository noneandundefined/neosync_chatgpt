# Архитектура NeoSync

## Обзор

NeoSync состоит из web-клиента, HTTP API и TCP-сервера терминалов. PostgreSQL хранит постоянные данные, Redis — сессии/кэш/оперативное состояние, RabbitMQ связывает HTTP и TCP части.

```text
Web
 │ HTTP
 ▼
HTTP service ───── PostgreSQL
 │  │
 │  ├──────────── Redis
 │  │
 │  └─ RabbitMQ ─────────┐
 │                       ▼
 └────────────────── TCP service ── TCP ── Terminals
                         │
                         ├──────── PostgreSQL
                         └──────── Redis
```

## Правило зависимостей backend

```text
cmd/server + handlers
          │
          ▼
   internal/usecase
          │
          ▼
 repository/store
          │
          ▼
        infra
 PostgreSQL / Redis / RabbitMQ
```

### cmd

Точки входа и composition root. Переиспользуемые клиенты и бизнес-логика здесь не хранятся.

### handler / TCP handlers

Transport layer: разобрать вход, вызвать use case, сформировать ответ. Сложные бизнес-правила сюда не добавляются.

### internal/usecase

Application layer: сценарии системы и координация операций. Use case не является частью PostgreSQL, поэтому этот слой отделён от `infra/store/postgres`.

### infra

Adapters внешних систем:
- `store/postgres` — модели и SQL;
- `store/redis` — Redis;
- `messaging/rabbitmq` — RabbitMQ;
- logger, locale, analytics и другие integrations.

### pkg

Переиспользуемые технические пакеты: ADM parsing, HTTP helpers, permissions, protocol codecs.

## HTTP service

```text
services/http/
├── cmd/server/
├── config/
├── handler/v1/
├── internal/usecase/
├── infra/
│   ├── messaging/rabbitmq/
│   └── store/
├── middleware/
├── migrations/
├── pkg/
├── tests/
└── util/
```

## TCP service

```text
services/tcp/
├── cmd/server/
├── config/
├── emulators/
├── internal/usecase/
├── infra/
│   ├── messaging/rabbitmq/
│   └── store/
├── pkg/protocol/
├── types/
└── util/
```

## Web

Web-приложение находится в `apps/web`. Это самостоятельный Vite-проект, поэтому оно не вложено в Java-подобный путь `client/src/main/resources`.

## Deploy

nginx, PostgreSQL и systemd — deployment-конфигурация, поэтому они находятся в явном каталоге `deploy/`, а не в скрытых корневых папках.

## Добавление новой функции

1. Handler принимает и валидирует вход.
2. Handler вызывает use case.
3. Use case выполняет бизнес-правила.
4. Use case обращается к store/integration.
5. Infra выполняет SQL/Redis/RabbitMQ операцию.
6. Handler переводит результат в transport response.

Если код нельзя однозначно отнести к одному из этих шагов, сначала стоит проверить ответственность пакета.
