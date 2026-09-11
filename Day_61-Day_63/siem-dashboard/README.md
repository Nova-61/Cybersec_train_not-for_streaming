# Day 61: простой SIEM Dashboard

Небольшой учебный проект на Ruby on Rails. Сейчас приложение выводит главную
страницу SIEM Dashboard и показывает базовую структуру Rails-приложения:
маршрут -> контроллер -> представление.

## Что такое SIEM

SIEM (Security Information and Event Management) собирает события
безопасности, хранит их и помогает находить подозрительную активность. В этом
учебном проекте пока нет настоящего сбора логов и авторизации: реализован только
каркас будущей панели мониторинга.

## Используемые технологии

- Ruby `4.0.6`;
- Ruby on Rails `8.1.3` или новее в рамках `8.1`;
- PostgreSQL;
- Puma — веб-сервер Rails;
- встроенный Rails-тест для проверки страницы.

## Как создать такой проект с нуля

Команды ниже рассчитаны на Windows PowerShell. Сначала нужно установить Ruby,
RubyGems и Bundler, затем проверить версии:

```powershell
ruby --version
gem --version
bundle --version
```

Создание Rails-приложения:

```powershell
gem install rails
rails new siem-dashboard
cd siem-dashboard
```

Установка зависимостей и подготовка базы данных:

```powershell
bundle install
bin\rails db:prepare
```

Создание контроллера, представления и теста:

```powershell
bin\rails generate controller Pages home
```

После этой команды Rails создаёт основные файлы:

- `app/controllers/pages_controller.rb` — контроллер;
- `app/views/pages/home.html.erb` — HTML-шаблон страницы;
- `test/controllers/pages_controller_test.rb` — тест;
- маршрут `pages_home` для адреса `/pages/home`.

В представлении можно разместить простой текст панели:

```erb
<h1>Мой SIEM Dashboard</h1>
<p>Добро пожаловать в мой SIEM Dashboard!</p>
```

Запуск локального сервера:

```powershell
bin\rails server
```

Откройте в браузере <http://localhost:3000/pages/home>.

## Как устроен текущий проект

Запрос к `/pages/home` проходит по цепочке:

1. `config/routes.rb` сопоставляет URL с действием `PagesController#home`.
2. `app/controllers/pages_controller.rb` выполняет метод `home`.
3. Rails автоматически выбирает `app/views/pages/home.html.erb`.
4. Шаблон превращается в HTML и отправляется браузеру.

Проверить состояние приложения можно по адресу `/up`: Rails вернёт успешный
ответ, если приложение запускается без ошибок.

## Запуск и проверка

Запуск тестов:

```powershell
bin\rails test
```

Проверка безопасности и стиля (необязательно):

```powershell
bin\bundler-audit check --update
bin\brakeman
bin\rubocop
```

## База данных и сервисы

PostgreSQL подключается стандартными настройками Rails. Команды для миграций,
моделей и работы с данными собраны в [README_DATABASE.md](README_DATABASE.md).
Redis, очереди задач и поисковые сервисы пока не используются.