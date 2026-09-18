# Code review NeoSync — 2026-09-18

Ветка: `chatgpt/code-review-fixes-20260918`  
База: `main`

## Исправлено

### TCP sessions / reconnect

- Старое TCP-соединение при завершении могло удалить уже созданную новую сессию того же IMEI.
  Добавлено удаление с проверкой принадлежности сессии конкретному `ADMDevice`.
- Очередь команд старой сессии проверяла только наличие IMEI в `sessions`, поэтому после reconnect могла продолжать работать с уже закрытым соединением.
  Теперь очередь проверяет identity текущей сессии.
- `AnswerCh` закрывался одновременно с возможной отправкой ответа из другого goroutine, что создавало риск `send on closed channel`.
  Для отмены добавлен отдельный `DoneCh`; канал ответов больше не закрывается конкурентно.
- Ответы команд, ошибки и пакеты конфигурации от старого соединения теперь не принимаются новой сессией того же IMEI.
- Для GET_CFG ошибка отправки в терминал раньше оставляла `RequestId` занятым до timeout.
  Теперь состояние запроса очищается сразу.

### RabbitMQ

- `Publish` мог взять AMQP channel из старого пула, а после reconnect вернуть его через уже новый `r.poolChannel`.
  Пул и connection теперь фиксируются на одну generation на время publish.
- Закрытие RabbitMQ-клиента больше не закрывает Go channel пула под активным publisher.
- Добавлен явный `closed` state, чтобы reconnect-loop не поднимал соединение после штатного `Close()`.
- При ошибке первичной инициализации каналов AMQP connection теперь закрывается.

### HTTP / auth

- Telemetry endpoints при отсутствии `AccessConfigurationRead` возвращали пустой успешный ответ.
  Теперь возвращается `403 Forbidden`.
- Process-level map `analyticsUserOnlineLastSeen` рос без очистки.
  Добавлена периодическая очистка устаревших записей.
- Асинхронная запись online-аналитики использовала request context, который мог быть отменён сразу после ответа.
  Для фоновой записи используется отдельный context с timeout.
- Получение raw/section configuration теперь корректно обрабатывает отсутствие строки конфигурации вместо возможного nil dereference.

### Redis drafts

- `HasDraft` использовал Redis `KEYS`, который блокирует Redis на больших keyspace.
  Заменено на cursor-based `SCAN`.
- `GetMergedDraft` выбирал `cfg_hash` в зависимости от недетерминированного порядка `SCAN`.
  Теперь metadata берётся из наиболее свежего draft.
- Защищено обновление legacy/corrupted draft с `changes: null`.

### Frontend

- Старые SSE callbacks могли обновить состояние уже после переключения IMEI/section/template.
  Добавлены active guards и корректная отмена async draft request.
- `onClose` SSE использовал stale значение `loading` из closure.
- `cooldownUntil` в `RequestControlUtils` хранил уникальные action keys бессрочно.
  Истёкшие записи теперь удаляются.
- При экспорте configuration не освобождался `ObjectURL`.
- Повреждённый JSON в localStorage для сортировки CompactTable мог уронить рендер.
  Добавлены parsing/shape validation и сброс некорректного cache.
- React Query Devtools больше не отображается в production и не открывается автоматически.

## Что стоит сделать отдельным рефакторингом

1. `useTemplateExitGuard` изменяет `UNSAFE_NavigationContext.navigator.push/replace` и вручную управляет history через `pushState/go(-2)`. Это хрупко при обновлении React Router и при сложной истории навигации. Лучше перейти на поддерживаемый blocker API после обновления router.
2. `GetDeviceSession` возвращает mutable `*TCPSession` после снятия lock. Это позволяет легко читать поля без `session.Mu`. Надёжнее постепенно заменить доступ к полям методами SessionMemory/TCPSession.
3. В UI существуют два похожих механизма template draft: `TemplateDraftUtils.ts` и `TemplateDraftStorageUtils.ts`. Перед удалением старого варианта нужно проверить все runtime imports.
4. Текущий `.forgejo/workflows/ci.yml` фактически является deploy workflow и не запускает отдельные lint/test checks для feature branches. Стоит добавить независимый CI: frontend `npm ci && npm run build && npm run lint`, backend `go test ./...` и отдельный race-check для TCP-кода.

## Проверка

Изменения сделаны только в отдельной ветке от `main`. `main` не изменялся.
