Rails.application.routes.draw do
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
