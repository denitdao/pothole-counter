package util

import "time"

func FormatDate(t time.Time) string {
	return t.Format("15:04 02.01.2006")
}
