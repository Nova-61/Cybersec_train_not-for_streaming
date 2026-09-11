class Event < ApplicationRecord
  validates :title, presence: true
  validates :event_time, presence: true
end
