Rails.application.routes.draw do
  # Аутентификация (День 65). Даёт маршруты /users/sign_in, /users/sign_up,
  # /users/sign_out, /users/password/new и т.д.
  devise_for :users

  # Devise после логина/логаута редиректит на root_path.
  root "events#index"

  # HTML-интерфейс для человека (День 63)
  resources :events

  # JSON-API для Go-агента (День 64)
  namespace :api do
    namespace :v1 do
      resources :events, only: [:create]
    end
  end

  # root "events#index"
end
