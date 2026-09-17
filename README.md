# Neosync

Конфигуратор удаленной настройки терминала от компании [Neomatica](https://neomatica.com/).

## Разработка Neosync

- [Neosync клиент](./client/src/docs/README.md)
- [Neosync HTTP сервер](./http-server/README.md)
- [Деплой Neosync сервиса](./scripts/README.md)

## Быстрый запуск Neosync

Для запуска сервиса Neosync на локальной машине потребуется - `Docker`.

> Все 'environment' в docker-compose.yml ТЕСТОВЫЕ, они не будут использоваться в PRE/PRODUCTION.

1. Сервис Neosync требует подключение к СУБД PostgreSQL. Если PostgreSQL нету на локальной машине, измените строку подключения в `docker-compose.yml` во ВСЕХ местах:

было:

```bash
DB_ADDR=postgres://postgres:password@host.docker.internal:5432/neosync?sslmode=disable
```

стало:

```bash
DB_ADDR=postgres://postgres:password@db:5432/neosync?sslmode=disable
```

2. Запуск через docker-compose

```bash
docker compose -f docker-compose.yml up --build
```

3. Проверка сервиса

- После запуска можно перейти на `localhost:5173` - откроется сервис Neosync.
- Для подключения терминалов/эмуляторов используйте - localhost:12346.

## Деплой сервиса

1. Подключение к серверу по SSH.

```bash
ssh -p <port> <user>@<ip>
```

2. Чтобы использовать исходный код, потребуется сервис `Tailscale`.

Установка Tailscale на Linux:

```bash
curl -fsSL https://tailscale.com/install.sh | sh
```

Для настройки Tailscale и авторизации устройства:

```bash
tailscale up
```

3. Клонирование репозитория.

```bash
git clone ssh://git@<ip>:222/neomatica/neosync.git
```

> После клонирования репозитория, добавьте переменные окружения: `.env.production` OR `.env.dev`.

4. Запуск сервиса Neosync.

```bash
cd ./neosync

git stash
git checkout main
git pull origin main

cd ./scripts

bash pre_deploy
bash systemd
bash setup_nginx

bash set_technical_work false
```
