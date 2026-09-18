# Deploy

Production-конфигурация вынесена из исходников приложений в один каталог.

```text
deploy/
├── nginx/       reverse proxy и web-serving
├── postgres/    PostgreSQL configuration
└── systemd/     HTTP/TCP unit files
```

Изменения здесь относятся к окружению развёртывания, а не к application logic.

Установку и обновление конфигурации выполняют скрипты из `scripts/`. Не копируйте deploy-логику внутрь `services/http`, `services/tcp` или `apps/web`.
