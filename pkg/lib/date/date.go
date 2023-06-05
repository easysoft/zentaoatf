package dateUtils

import (
	"time"

	"github.com/easysoft/zentaoatf/pkg/consts"
)

func DateStr(tm time.Time) string {
	return tm.Format("2006-01-02")
}
func DateStrShort(tm time.Time) string {
	return tm.Format("06-01-02")
}

func TimeStr(tm time.Time) string {
	return tm.Format("15:04:05")
}

func DateTimeStrFmt(tm time.Time, fm string) string {
	return tm.Format(fm)
}

func DateTimeStr(tm time.Time) string {
	return tm.Format(consts.DateTimeFormat)
}

func DateTimeStrLong(tm time.Time) string {
	return tm.Format("20060102150405")
}

func DateStrToTimestamp(str string) (int64, error) {
	layout := "20060102"

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0, err
	}

	time, err := time.ParseInLocation(layout, str, loc)
	if err != nil {
		return 0, err
	}

	return time.Unix(), nil
}

func TimeStrToTimestamp(str string) int64 {
	layout := "2006-01-02 15:04:05"

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0
	}

	time, err := time.ParseInLocation(layout, str, loc)
	if err != nil {
		return 0
	}

	return time.Unix()
}

func UnitToDate(unit int64) (date time.Time, err error) {
	timeStr := time.Unix(unit, 0).Format(consts.DateTimeFormat)

	date, _ = time.ParseInLocation(consts.DateTimeFormat, timeStr, time.Local)

	return
}
