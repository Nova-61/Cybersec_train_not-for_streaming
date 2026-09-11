module Api
  module V1
    class EventsController < ApplicationController
      skip_before_action :verify_authenticity_token
      before_action :authenticate_api_key

      def create
        if params[:entries].present?
          create_batch
        else
          create_single
        end
      end

      private

      # Реальный Go-агент шлёт именно это: {"entries": [...], "count": N, "source": "siem-agent"}
      def create_batch
        created_ids = []
        failed = []

        entries_params.each do |entry|
          event = Event.new(entry)
          if event.save
            created_ids << event.id
          else
            failed << { data: entry.to_h, errors: event.errors.full_messages }
          end
        end

        status = failed.empty? ? :created : :multi_status
        render json: {
          status: failed.empty? ? "ok" : "partial",
          created: created_ids.size,
          failed: failed.size,
          errors: failed
        }, status: status
      end

      # Оставляю и одиночный режим — вдруг понадобится для ручного теста через curl
      def create_single
        event = Event.new(event_params)
        if event.save
          render json: { status: "ok", id: event.id, message: "Event saved" }, status: :created
        else
          render json: { status: "error", errors: event.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def authenticate_api_key
        api_key = request.headers["X-API-Key"]
        expected_key = ENV.fetch("SIEM_API_KEY", "dev-secret-key-12345")

        unless api_key == expected_key
          render json: { status: "error", message: "Unauthorized" }, status: :unauthorized
        end
      end

      def event_params
        params.require(:event).permit(:level, :message, :source, :user_id)
      end

      def entries_params
        params.require(:entries).map { |e| e.permit(:level, :message, :source) }
      end
    end
  end
end