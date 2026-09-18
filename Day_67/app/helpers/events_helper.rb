module EventsHelper
  def level_badge_color(level)
    {
      "INFO" => "info",
      "WARNING" => "warning",
      "ERROR" => "danger",
      "CRITICAL" => "dark",
      "DEBUG" => "secondary"
    }.fetch(level, "secondary")
  end
end
