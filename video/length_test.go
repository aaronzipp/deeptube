package video

import (
	"reflect"
	"testing"
)

func TestLengthFromString(t *testing.T) {
	testData := []struct {
		input  string
		output Length
	}{
		{input: "PT1H2M3S", output: Length{Hours: 1, Minutes: 2, Seconds: 3}},
		{input: "PT10M", output: Length{Hours: 0, Minutes: 10, Seconds: 0}},
		{input: "PT10H10M10S", output: Length{Hours: 10, Minutes: 10, Seconds: 10}},
	}

	for _, tt := range testData {
		t.Run(tt.input, func(t *testing.T) {
			got, err := LengthFromString(tt.input)

			if err != nil {
				t.Fatalf("Got an unexpected error: %q", err)
			}

			if !reflect.DeepEqual(tt.output, got) {
				t.Errorf("Got %+v, want %+v", got, tt.output)
			}
		})
	}
}

func TestString(t *testing.T) {
	length := Length{1, 2, 3}

	got := length.String()
	want := "1h 2m 3s"

	if got != want {
		t.Errorf("Got %q, want %q", got, want)
	}
}
