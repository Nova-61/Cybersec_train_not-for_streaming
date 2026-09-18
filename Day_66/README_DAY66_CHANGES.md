# День 66 — Авторизация (Pundit)

Это твоя папка `siem_api` (уже с исправлениями из Дня 65: devise,
исправленный Gemfile, HTML events-контроллер), плюс изменения Дня 66.

## Новые файлы

- `db/migrate/20260914000000_add_role_to_users.rb` — добавляет колонку
  `role` (string, по умолчанию `"user"`) в таблицу `users`.
- `app/policies/application_policy.rb` — базовый класс Pundit (стандартный
  шаблон из `pundit:install`).
- `app/policies/event_policy.rb` — правила доступа к `Event`: свой объект
  или роль `admin`.

## Изменённые файлы

- `app/models/user.rb` — добавлен метод `admin?` (`role == "admin"`).
- `Gemfile` — добавлена строка `gem "pundit"`.
- `app/controllers/application_controller.rb` — подключён
  `Pundit::Authorization`, добавлен `rescue_from Pundit::NotAuthorizedError`
  (без доступа → редирект на `/events` с сообщением, а не страница 500).
- `app/controllers/events_controller.rb` — `index` теперь через
  `policy_scope(Event)` (обычный юзер видит только свои события, админ —
  все), остальные action'ы обёрнуты в `authorize @event`.
- `app/views/events/show.html.erb` — ссылки "Редактировать"/"Удалить"
  показываются только если `policy(@event)` это разрешает.

## Порядок запуска (PowerShell, из папки проекта)

```powershell
bundle install
rails generate pundit:install
```

Генератор `pundit:install` создаст `app/policies/application_policy.rb`
заново своим шаблоном — если он спросит про перезапись существующего
файла, ответь `n` (наш файл уже такой же по смыслу, просто не давай ему
затирать `event_policy.rb`, если он вдруг попробует).

```powershell
rails db:migrate
```

Применит миграцию `role`.

## Делаем себя админом

```powershell
rails console
```
```ruby
User.first.update(role: "admin")
exit
```

## Проверка

```powershell
rails server
```

1. Зарегистрируй **второго** пользователя (обычного, не админа) —
   через `/users/sign_up` в режиме инкогнито, чтобы не разлогинить первого.
2. Под обычным пользователем создай событие.
3. Разлогинься, зайди под админом (`User.first`) — должен видеть **все**
   события в `/events`, включая чужие.
4. Под обычным юзером попробуй открыть `/events/<id чужого события>` —
   на странице не будет кнопок "Редактировать"/"Удалить". Если зайти на
   `/events/<id>/edit` напрямую по URL — должен сработать редирект на
   `/events` с сообщением "У вас нет доступа к этому действию." — это и
   есть серверная защита, а не просто спрятанная кнопка.

```powershell
rails runner "puts User.pluck(:email, :role)"
```
Покажет всех пользователей и их роли — удобно проверить, что миграция
и обновление роли применились.
