module Api
  module V1
    class EventsController < ApplicationController
      # Этот контроллер обслуживает API версии v1.
      # Маршрут для создания события: POST /api/v1/events.

      # Go-агент отправляет обычный HTTP-запрос, а не HTML-форму.
      # Поэтому у него нет CSRF-токена, который Rails обычно требует
      # от браузерных форм. Для этого API-пути CSRF-проверка отключена.
      skip_before_action :verify_authenticity_token

      # Перед выполнением create Rails вызывает authenticate_api_key!.
      # Go должен передать ключ в HTTP-заголовке X-API-Key.
      before_action :authenticate_api_key!

      # Принимает событие от Go-агента и сохраняет его в базу данных.
      def create
        # event_params содержит только разрешенные параметры из JSON-запроса.
        @event = Event.new(event_params)

        if @event.save
          # Событие успешно сохранено.
          # HTTP 201 означает Created, а в JSON возвращается созданная запись.
          render json: @event, status: :created
        else
          # Данные не прошли валидацию модели Event.
          # HTTP 422 означает, что запрос понятен, но данные некорректны.
          render json: { errors: @event.errors.full_messages }, status: :unprocessable_entity
        end
      end

      private

      # Разрешает параметры, которые можно записать в Event.
      # Из Go нужно отправлять JSON с вложенным объектом event, например:
      # {
      #   "event": {
      #     "title": "Ошибка",
      #     "description": "Описание ошибки",
      #     "event_time": "2026-09-12T10:30:00Z"
      #   }
      # }
      #
      # params.require(:event) требует наличие объекта event.
      # params.permit(...) защищает приложение от записи лишних полей.
      def event_params
        params.require(:event).permit(:title, :description, :event_time)
      end

      # Проверяет API-ключ Go-агента.
      def authenticate_api_key!
        # Секрет читается из config/credentials.yml.enc.
        # В credentials должен быть записан ключ верхнего уровня:
        # api_key: my_super_secret_key_12345
        expected_key = Rails.application.credentials.api_key

        # Ключ из заголовка должен полностью совпадать с ключом Rails.
        # При несовпадении запрос получает HTTP 401 Unauthorized.
        unless request.headers["X-API-Key"] == expected_key
          render json: { error: "Неверный API-ключ" }, status: :unauthorized
        end
      end
    end
  end
end
