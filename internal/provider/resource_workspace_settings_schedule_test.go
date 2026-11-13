package provider

import "testing"

func TestFeedScheduleValue(t *testing.T) {
	cases := map[FeedSummarySchedule]string{
		"":       "daily",
		"daily":  "daily",
		"weekly": "weekly",
	}

	for schedule, want := range cases {
		if got := feedScheduleValue(schedule); got != want {
			t.Errorf("feedScheduleValue(%q) = %q, want %q", schedule, got, want)
		}
	}
}
