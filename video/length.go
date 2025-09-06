package video

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Length struct {
	Hours   int
	Minutes int
	Seconds int
}

func (l Length) String() string {
	if l.Hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", l.Hours, l.Minutes, l.Seconds)
	}
	if l.Minutes > 0 {
		return fmt.Sprintf("%dm %ds", l.Minutes, l.Seconds)
	}
	return fmt.Sprintf("%ds", l.Seconds)
}

func LengthFromString(lengthText string) (Length, error) {
	if lengthText == "P0D" {
		return Length{0, 0, 0}, nil
	}
	re := regexp.MustCompile("^PT([0-9]+H)?([0-9]+M)?([0-9]+S)?$")

	matches := re.FindStringSubmatch(lengthText)

	hours, err := parseInt(matches[1])
	if err != nil {
		return Length{}, err
	}
	minutes, err := parseInt(matches[2])
	if err != nil {
		return Length{}, err
	}
	seconds, err := parseInt(matches[3])
	if err != nil {
		return Length{}, err
	}

	return Length{hours, minutes, seconds}, nil

}

func parseInt(text string) (int, error) {
	if len(text) > 1 {
		text = text[:len(text)-1]
	} else {
		if len(text) == 1 {
			return 0, errors.New("Got time substring of length 1. This should not be possible.")
		}
		text = "0"
	}
	return strconv.Atoi(text)
}
