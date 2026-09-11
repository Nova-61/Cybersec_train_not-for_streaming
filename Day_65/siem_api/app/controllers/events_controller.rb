class EventsController < ApplicationController
  def index
    @events = Event.order(created_at: :desc).limit(50)
  end
end
