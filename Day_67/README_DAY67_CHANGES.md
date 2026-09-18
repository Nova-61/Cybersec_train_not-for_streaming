# День 67 — Bootstrap и UI

Это твоя папка `siem_api` (со всеми предыдущими днями), плюс Bootstrap
поверх существующих HTML-страниц. Логика (контроллеры, policy, Devise)
не менялась вообще — только вёрстка.

## Изменённые файлы

- `app/views/layouts/application.html.erb`
  - Добавлен `stylesheet_link_tag` на Bootstrap CSS (CDN) — рядом с твоим
    существующим `stylesheet_link_tag :app`, не вместо него.
  - `<nav>` переделан в Bootstrap-навбар (тёмный, с email и меткой
    `(admin)`, если пользователь админ).
  - `<%= yield %>` обёрнут в `<div class="container">`, flash-сообщения
    теперь `alert alert-success` / `alert alert-danger`.
  - Перед `</body>` добавлен `<script>` с Bootstrap JS (bundle-версия,
    нужна для будущих дропдаунов/тогглов, хотя сейчас не используется).

- `app/views/events/index.html.erb` — таблица `table table-striped
  table-hover`, уровень события показан цветным `badge`.
- `app/views/events/show.html.erb` — карточка `card` вместо голых `<p>`.
- `app/views/events/new.html.erb`, `edit.html.erb` — Bootstrap-классы на
  полях формы (`form-control`, `form-select`, `form-label`), + блок вывода
  ошибок валидации (`alert alert-danger` со списком).

## Новый файл

- `app/helpers/events_helper.rb` — метод `level_badge_color(level)`,
  сопоставляет уровень события (`INFO`/`WARNING`/`ERROR`/`CRITICAL`/`DEBUG`)
  с цветом Bootstrap-бейджа.

## Что НЕ трогал

- `config/initializers/content_security_policy.rb` — весь файл
  закомментирован по умолчанию в твоём проекте, значит CSP не активен
  и CDN-ссылки на Bootstrap не заблокируются. Если ты его когда-нибудь
  раскомментируешь — не забудь добавить `cdn.jsdelivr.net` в
  `style_src`/`script_src`.
- Контроллеры, policy, миграции, Devise — без изменений.

## Порядок запуска (PowerShell, из папки проекта)

Новых гемов/миграций в этом дне нет — просто:

```powershell
rails server
```

## Проверка

Зайди на `http://localhost:3000/events` — должен появиться тёмный навбар,
кнопка "+ Добавить событие", таблица с цветными бейджами уровней.

Открой DevTools (F12) → вкладка Network → обнови страницу → найди
`bootstrap.min.css` — должен быть статус `200`. Если красный/заблокирован —
смотри раздел "Типичные ошибки" в объяснении Дня 67 (CSP).

Зайди на `/events/new` — поля формы должны выглядеть как настоящие
Bootstrap input'ы (со скруглением, паддингами), а не голый HTML.
