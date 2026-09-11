require "test_helper"

class EventsControllerTest < ActionDispatch::IntegrationTest
  test "should get index" do
    get events_url
    assert_response :success
  end

  test "should get stats" do
    get "/events/stats"
    assert_response :success
  end

  test "should get show" do
    event = Event.create!(level: "INFO", message: "test message", source: "test-source")

    get event_url(event)
    assert_response :success
  end
end
