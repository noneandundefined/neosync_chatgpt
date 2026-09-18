# NeoSync scripts

Эксплуатационные скрипты репозитория.

- `pre_deploy` — проверяет исходники, собирает HTTP/TCP/Web перед deploy.
- `deploy` — запускает полный deploy pipeline.
- `systemd` — устанавливает и перезапускает unit-файлы из `deploy/systemd`.
- `setup_nginx` — собирает web и публикует `dist` в nginx.
- `set_technical_work` — включает/выключает технический режим.
- `psql_backup`, `psql_backup_cron` — резервные копии PostgreSQL.
- `analyze_imei_traffic` — диагностический анализ трафика терминалов.
- `lw` — работа с логами.
- `cache` — операции с кэшем.
- `share` — вспомогательный sharing script.

Скрипты рассчитывают путь к корню через Git, поэтому запускаются из checkout NeoSync и используют актуальные каталоги `apps/`, `services/` и `deploy/`.

Для обычной разработки deploy-скрипты не нужны — используйте команды из корневого README.
