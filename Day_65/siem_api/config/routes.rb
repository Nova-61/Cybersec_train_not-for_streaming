Rails.application.routes.draw do
  devise_for :users

  root "events#index"
  resources :events, only: [:index]

  namespace :api do
    namespace :v1 do
      resources :events, only: [:create]
    end
  end

  get "/health", to: proc { [200, {}, ["OK"]] }
end
