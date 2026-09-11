# Работа с базой данных в Rails

Этот файл описывает основные команды Active Record и Rails для PostgreSQL в
проекте `siem-dashboard`.

## 1. Подготовка PostgreSQL

В `config/database.yml` проект использует PostgreSQL на `localhost:5432`:

- пользователь: `postgres`;
- пароль для локальной разработки: `123`;
- база разработки: `siem_dashboard_development`;
- тестовая база: `siem_dashboard_test`.

После запуска PostgreSQL установите зависимости и создайте базы:

```powershell
bundle install
bin\rails db:create
bin\rails db:migrate
```

Команда `db:prepare` удобна для первого запуска: она создаёт базу, если её
нет, и применяет миграции.

```powershell
bin\rails db:prepare
```

> Для реального проекта пароль не стоит хранить в `database.yml`. Лучше передать
> его через переменную окружения или файл credentials.

## 2. Основные команды Rails для базы

Посмотреть все доступные задачи, связанные с базой:

```powershell
bin\rails --tasks | Select-String "db:"
```

Популярные команды:

| Команда | Назначение |
| --- | --- |
| `bin\rails db:create` | создать базы development и test |
| `bin\rails db:drop` | удалить базы |
| `bin\rails db:reset` | удалить, создать и заполнить базу заново |
| `bin\rails db:migrate` | применить новые миграции |
| `bin\rails db:rollback` | отменить последнюю миграцию |
| `bin\rails db:rollback STEP=2` | отменить две последние миграции |
| `bin\rails db:migrate:status` | показать статус миграций |
| `bin\rails db:seed` | выполнить `db/seeds.rb` |
| `bin\rails db:prepare` | создать базу и применить миграции |
| `bin\rails db:schema:load` | загрузить структуру из `db/schema.rb` |
| `bin\rails db:version` | показать версию схемы |

Команды `db:drop` и `db:reset` удаляют данные. Используйте их только когда
понимаете последствия.

## 3. Миграции

Миграция описывает изменение структуры базы данных. Создать модель события и
таблицу можно одной командой:

```powershell
bin\rails generate model Event source:string severity:string message:text occurred_at:datetime
```

Rails создаст:

- модель `app/models/event.rb`;
- миграцию в `db/migrate/`;
- тест модели.

Пример миграции:

```ruby
class CreateEvents < ActiveRecord::Migration[8.1]
  def change
    create_table :events do |table|
      table.string :source, null: false
      table.string :severity, null: false, default: "info"
      table.text :message
      table.datetime :occurred_at, null: false
      table.timestamps
    end

    add_index :events, :severity
    add_index :events, :occurred_at
  end
end
```

Применить миграцию:

```powershell
bin\rails db:migrate
```

Создать новую миграцию для изменения таблицы:

```powershell
bin\rails generate migration AddIpAddressToEvents ip_address:string
bin\rails db:migrate
```

Отменить последнюю миграцию:

```powershell
bin\rails db:rollback
```

Откатить конкретную миграцию по версии можно так:

```powershell
bin\rails db:migrate:down VERSION=20260906104530
```

Вернуть её обратно:

```powershell
bin\rails db:migrate:up VERSION=20260906104530
```

## 4. Модель Active Record

Модель `Event` представляет строку таблицы `events`. Добавьте в модель
валидации, чтобы некорректные события не сохранялись:

```ruby
class Event < ApplicationRecord
  validates :source, :severity, :message, presence: true
  validates :severity, inclusion: { in: %w[info low medium high critical] }
end
```

После изменения модели миграция обычно не нужна. Миграция нужна, когда меняется
сама структура таблицы: столбцы, индексы или ограничения.

## 5. CRUD через Rails Console

Открыть консоль Rails:

```powershell
bin\rails console
```

Создать событие:

```ruby
event = Event.create!(
  source: "web-server-01",
  severity: "high",
  message: "Много неудачных попыток входа",
  occurred_at: Time.current
)
```

Найти записи:

```ruby
Event.all
Event.first
Event.find(1)
Event.where(severity: "high")
Event.where("occurred_at >= ?", 1.hour.ago)
Event.order(occurred_at: :desc).limit(10)
```

Изменить запись:

```ruby
event.update!(severity: "critical")
```

Удалить запись:

```ruby
event.destroy!
```

Выйти из консоли:

```ruby
exit
```

`create!` и `update!` выбрасывают ошибку, если валидация не пройдена. Это
удобно во время разработки, потому что причина ошибки сразу видна.

## 6. Seeds

Файл `db/seeds.rb` нужен для тестовых данных:

```ruby
Event.create!(
  source: "firewall-01",
  severity: "medium",
  message: "Подозрительный входящий трафик",
  occurred_at: 15.minutes.ago
)
```

Запустить seeds:

```powershell
bin\rails db:seed
```

Чтобы пересоздать базу и снова загрузить seeds:

```powershell
bin\rails db:reset
```

## 7. Запросы и статистика для Dashboard

Получить критические события:

```ruby
Event.where(severity: "critical").order(occurred_at: :desc)
```

Посчитать количество событий по уровням опасности в PostgreSQL:

```ruby
Event.group(:severity).count
```

Получить события за последние сутки:

```ruby
Event.where(occurred_at: 24.hours.ago..Time.current)
```

Количество событий по источникам:

```ruby
Event.group(:source).count
```

## 8. Проверка базы и тесты

Показать статус миграций:

```powershell
bin\rails db:migrate:status
```

Запустить тесты:

```powershell
bin\rails test
```

Проверить только тесты моделей:

```powershell
bin\rails test test/models
```

Полезно проверять миграции на тестовой базе перед изменениями в production:

```powershell
bin\rails db:prepare RAILS_ENV=test
bin\rails test
```

## 9. Рекомендуемый порядок работы

1. Изменить модель или создать миграцию.
2. Проверить файл миграции.
3. Выполнить `bin\rails db:migrate`.
4. Добавить или обновить валидации модели.
5. Проверить данные через `bin\rails console`.
6. Запустить `bin\rails test`.
7. Перед деплоем проверить статус миграций и сделать резервную копию базы.
