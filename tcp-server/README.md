/cmd/
  server/              # Точка входа TCP сервера (main.go)
/pkg/
  config/              # Загрузка конфигурации (yaml/env)
  logger/              # Настройка логгера (zap/logrus)
  tracker/
    adm/
      parser.go        # Парсинг бинарных пакетов ADM
      encoder.go       # Сборка ответов ADM
      consts.go        # Константы протокола ADM
      types.go         # Структуры пакетов ADM
    dispatcher.go      # Маршрутизация данных по типу трекеров
  connection/
    client.go          # Работа с одним подключением (TCP client)
    manager.go         # Управление всеми соединениями (map IMEI → Client)
  handler/
    session.go         # Управление сессией (анализ handshake, ACK, команды)
    processor.go       # Обработка приходящих данных
  storage/
    repository.go      # Интерфейсы хранилища
    memory/            # In-memory реализация
    postgres/          # База данных PostgreSQL (pgx/sqlx)
    redis/             # Redis для очередей/кеша
  command/
    sender.go          # Отправка команд трекерам (асинхронно/синхронно)
    queue.go           # Очередь команд, если используется Redis или другая очередь
  utils/
    crc.go             # Проверка контрольных сумм
    hex.go             # Утилиты для работы с hex/binary
    time.go            # Парсинг времени из пакетов
/docs/
  protocol/            # Документация по протоколу ADM
  diagrams/            # Диаграммы архитектуры
/test/
  fixtures/            # Тестовые пакеты (raw binary samples)
/scripts/
  migrate.sh           # Миграции базы, утилиты
go.mod
go.sum
README.md


/pkg/
  ...
  protocol/
    factory.go            # Определяет версию и возвращает соответствующий парсер/кодер
    admv1/
      parser.go           # Парсер пакетов ADM v1
      encoder.go          # Сборщик пакетов ADM v1
      types.go            # Структуры пакетов ADM v1
      consts.go
    admv2/
      parser.go           # Парсер пакетов ADM v2
      encoder.go
      types.go
      consts.go
    common/
      interface.go        # Общие интерфейсы PacketParser, Packet
      utils.go            # Общие утилиты (CRC, время и т.д.)
