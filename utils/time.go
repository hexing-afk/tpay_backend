package utils

import (
	"log"
	"time"
)

var DefaultLocation *time.Location

// 设置时区
func SetTimezone(timezone string) {
	local, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatal(err)
	}
	// 设置time包中的默认时区
	time.Local = local

	DefaultLocation = local
}

func Unixtime2Time(timeIn interface{}, location *time.Location) time.Time {
	tmp := ToInt64(timeIn)
	return time.Unix(tmp, 0).In(location)
}

func TodaySec(location *time.Location) int64 {
	t := time.Now().In(location)
	return int64(t.Hour()*60*60 + t.Minute()*60 + t.Second())
}
