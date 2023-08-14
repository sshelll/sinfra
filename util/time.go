package util

import "time"

func TimeToMilliTs(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

func MilliTsToTime(ts int64) time.Time {
	return time.Unix(0, ts*1000000)
}

func NowMilliTs() int64 {
	return TimeToMilliTs(time.Now())
}
