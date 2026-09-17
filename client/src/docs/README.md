# Neosync клиент

Это клиентская часть (frontend) системы управления терминалами Neomatica. Приложение реализовано на React с использованием TypeScript, React Router, React Query, Context API и поддержкой i18n (многоязычность). Система предназначена для управления устройствами, пользователями, отправки команд, конфигурирования терминалов и т.д.

## Быстрый старт

- Установка зависимостей

```bash
pnpm install
```

- Запуск в режиме разработки

```bash
pnpm dev
```

- Сборка для продакшена

```bash
pnpm build
```

## Архитектура и структура

├── components/ - UI-компоненты (Header SideBar, таблицы, модальные окна и др.)<br>
├── constants/ - Константы (роуты, роли, статусы и др.)<br>
├── context/ - Контексты React (например, RoleContext)<br>
├── hooks/ - Кастомные хуки<br>
├── interface/ - TypeScript-интерфейсы<br>
├── locale/ - Локализации (en, ru, es)<br>
├── pages/ - Страницы приложения (SPA)<br>
├── rest/ - API-клиенты (axios, sse, fetch)<br>
├── router/ - Роутинг (React Router)<br>
├── styles/ - Стили (SASS, CSS)<br>
├── utils/ - Утилиты<br>

## Аутентификация и роли

- Аутентификация реализована через API (`/auth/signin`, `/auth/signup`, `/auth/signout`).
- Для управления доступом используется **RoleContext** (см. `context/RoleContext/useRoleContext.tsx`).
- Роли определены в `constants/Roles.constant.ts`:
    - SUPERADMIN
    - SUPPORT
    - ADMINL2
    - USER

> Доступ к страницам и действиям определяется ролью пользователя.

## Работа с API

API-клиенты реализованы через **axios** в папке `rest/`:

## Локализация

- Поддержка русского, английского и испанского языка.
- Файлы переводов: `locale/languages/en-us.json`, `locale/languages/ru-ru.json`, `locale/languages/en-es.json`.
- Смена языка — через меню пользователя в Header.

## Стилизация

- Используется SASS и Tailwind CSS.
- Готовые компоненты стилизованы через классы и отдельные файлы в `styles/components/`.

## Советы для разработки

- Для добавления новых страниц — регистрируйте их в `router/config.ts` и добавляйте компонент в `pages/`.
- Для новых API — создавайте отдельный файл в `rest/` и используйте через хуки/React Query(см. `hooks/Server/useHandleServer.ts`).
- Для новых ролей — добавьте их в `Roles.constant.ts` и настройте доступ в компонентах/роутинге.
- Для новых языков — добавьте json-файл в `locale/languages/` и подключите в i18n.
