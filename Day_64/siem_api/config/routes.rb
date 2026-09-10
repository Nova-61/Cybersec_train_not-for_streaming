Rails.application.routes.draw do
  namespace :api do
    namespace :v1 do
      resources :events, only: [:create]
    end
  end
  get "/health", to: proc { [200, {}, ["OK"]] }
end
