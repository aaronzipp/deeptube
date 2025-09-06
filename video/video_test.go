package video

import (
	"testing"
	"time"
)

func TestTimeSincePublished(t *testing.T) {
	testData := []struct {
		input  time.Time
		output string
	}{
		{
			input: time.Date(
				time.Now().Year()-1,
				time.Now().Month(),
				time.Now().Day() - 1,
				0,
				time.Now().Minute(),
				time.Now().Second(),
				time.Now().Nanosecond(),
				time.UTC,
			),
			output: "1 year ago",
		},
		{
			input: time.Date(
				time.Now().Year(),
				time.Now().Month()-2,
				time.Now().Day(),
				time.Now().Hour(),
				time.Now().Minute(),
				time.Now().Second(),
				time.Now().Nanosecond(),
				time.UTC,
			),
			output: "2 months ago",
		},
	}

	for _, tt := range testData {
		t.Run(tt.input.String(), func(t *testing.T) {
			exampleVideo := Video{
				PublishedAt: tt.input,
			}
			got := exampleVideo.TimeSincePublished()

			if got != tt.output {
				t.Errorf("Want %q, got %q", tt.output, got)
			}
		})
	}
}
