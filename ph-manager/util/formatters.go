package util

import (
	"strconv"
	"time"
)

func FormatDate(t time.Time) string {
	return t.Format("15:04 02.01.2006")
}

func Mul(a, b float32) float32 {
	return a * b
}

func FormatFloat(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', 1, 32)
}
