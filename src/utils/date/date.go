package dateUtils

import (
	"time"
)

func DateStr(tm time.Time) string {
	return tm.Format("2006-01-02")
}

func TimeStr(tm time.Time) string {
	return tm.Format("15:04:05")
}

func DateTimeStrFmt(tm time.Time, fm string) string {
	return tm.Format(fm)
}

func DateTimeStr(tm time.Time) string {
	return tm.Format("2006-01-02 15:04:05")
}

func DateTimeStrLong(tm time.Time) string {
	return tm.Format("20060102150405")
}
