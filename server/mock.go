package server

import (
	"time"
)

func nowIsh() time.Time {
	t := time.Now()
	return t.Round(time.Hour)
}
