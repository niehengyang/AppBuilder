package utils

import (
	"time"
)

const (
	TimeStrTemplate1 = "2006-01-02 15:04:05"
	TimeStrTemplate2 = "2006/01/02 15:04:05"
	TimeStrTemplate3 = "2006-01-02"
	TimeStrTemplate4 = "15:04:05"
	TimeStrTemplate5 = "2006-01-02 15:04"
)

//
// TimestampFormat
//  @Description: 时间戳转字符串（秒）
//  @param timestamp	时间戳
//  @param formatTemplate	时间字符串格式
//  @return timestr	时间字符串
//
func TimestampFormat(timestamp int64, formatTemplate string) (timestr string) {
	timestr = time.Unix(timestamp, 0).Format(formatTemplate)
	return timestr
}

/*
将一段时间转换为hh:mm:ss 格式的字符串，比如将3523秒转换为00:58:43，
@duration 时间长度，以s为单位
*/
func DurationFormat(duration int) string {
	// 将持续时间转换为 time.Duration 类型
	durationInDuration := time.Duration(duration) * time.Second

	// 将持续时间转换为 hh:mm:ss 格式的字符串
	durationFormatted := time.Time{}.Add(durationInDuration).Format("15:04:05")

	return durationFormatted
}
