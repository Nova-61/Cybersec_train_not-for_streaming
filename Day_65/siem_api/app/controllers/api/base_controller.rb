module Api
  class BaseController < ActionController::API
    # API не требует браузерной аутентификации — сюда before_action :authenticate_user! не добавляем
  end
end
