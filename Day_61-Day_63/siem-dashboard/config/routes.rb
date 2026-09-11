Rails.application.routes.draw do
  # Все 7 стандартных RESTful-маршрутов для Event:
  # index, show, new, create, edit, update, destroy
  resources :events

  # Можно сделать events стартовой страницей сайта (необязательно):
  # root "events#index"
end
