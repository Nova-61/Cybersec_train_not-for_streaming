class Event < ApplicationRecord
  validates :level, presence: true, inclusion: { in: %w[INFO WARNING ERROR DEBUG CRITICAL] }
  validates :message, presence: true
  validates :source, presence: true
end
