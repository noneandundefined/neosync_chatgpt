## Neomatica WebRC

> Это клиентская часть (frontend) системы управления терминалами Neomatica. Приложение реализовано на React с использованием TypeScript, React Router, React Query, Context API и поддержкой i18n (многоязычность). Система предназначена для управления устройствами, пользователями, отправки команд, конфигурирования терминалов и т.д.

### Быстрый старт

- Установка зависимостей

```bash
pnpm install
# или
npm install
```

- Запуск в режиме разработки

```bash
pnpm run dev
# или
npm run dev
```

- Сборка для продакшена

```bash
pnpm run build
# или
npm run build
```

### Архитектура и структура

├── components/ # UI-компоненты (Header SideBar, таблицы, модальные окна и др.)
├── constants/ # Константы (роуты, роли, статусы и др.)
├── context/ # Контексты React (например, RoleContext)
├── hooks/ # Кастомные хуки
├── interface/ # TypeScript-интерфейсы
├── locale/ # Локализации (en, ru)
├── pages/ # Страницы приложения (SPA)
├── rest/ # API-клиенты (axios)
├── router/ # Роутинг (React Router)
├── styles/ # Стили (SASS, CSS)
├── utils/ # Утилиты

## Аутентификация и роли

- Аутентификация реализована через API (`/auth/signin`, `/auth/signup`, `/auth/signout`).
- После входа данные пользователя (login, uuid, role_id) сохраняются в `localStorage` (`user-auth`).
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

- Поддержка русского и английского языков.
- Файлы переводов: `locale/languages/en-us.json`, `locale/languages/ru-ru.json`.
- Смена языка — через меню пользователя в Header.

## Стилизация

- Используется SASS и Tailwind CSS.
- Готовые компоненты стилизованы через классы и отдельные файлы в `styles/components/`.

## Советы для разработки

- Для добавления новых страниц — регистрируйте их в `router/config.ts` и добавляйте компонент в `pages/`.
- Для новых API — создавайте отдельный файл в `rest/` и используйте через хуки/React Query.
- Для новых ролей — добавьте их в `Roles.constant.ts` и настройте доступ в компонентах/роутинге.
- Для новых языков — добавьте json-файл в `locale/languages/` и подключите в i18n.
