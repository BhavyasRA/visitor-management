package helpers

import "time"

func NowPointer() *time.Time {
	now := time.Now()
	return &now
}

func AddMinutesPointer(minutes int) *time.Time {
	t := time.Now().Add(time.Duration(minutes) * time.Minute)
	return &t
}