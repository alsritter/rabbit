package timetools

import (
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	TimeDefaultLocation               = "Asia/Shanghai" // 默认时区
	TimeFormatLayout_Y_m_D_H_M_S      = "2006-01-02 15:04:05"
	TimeFormatLayout_Y_m_D_H_M_S_CST  = "2006-01-02 15:04:05 +0800 CST"
	TimeFormatLayout_YmDHMS           = "20060102150405"
	TimeFormatLayout_YmDhMS           = "20060102030405"
	TimeFormatLayout_Y_m_D            = "2006-01-02"
	TimeFormatLayout_YmD              = "20060102"
	TimeFormatChineseLayout_Y_m_D_H_M = "2006年01月02日15:04"
	TimeFormatChineseLayout_Y_m_D     = "2006年01月02日"
	TimeFormatLayout_Y_m_D_H_M        = "2006-01-02 15:04"
	TimeFormatLayout_H_M              = "15:04"
	TimeFormatLayout_H_M_S            = "15:04:05"
)

func TimestamppbToTime(t *timestamppb.Timestamp) time.Time {
	return time.Unix(t.GetSeconds(), int64(t.GetNanos()))
}

func TimeToTimestamppb(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func Now() *time.Time {
	t := time.Now()
	return &t
}

func TimeToPinter(t time.Time) *time.Time {
	return &t
}

func TimestamppbNow() *timestamppb.Timestamp {
	return TimeToTimestamppb(time.Now())
}

func PointerTimeToTimestamppb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return TimeToTimestamppb(*t)
}

func TimestampToPointerTime(t int64) *time.Time {
	if t == 0 {
		return nil
	}
	tm := time.UnixMilli(t)
	return &tm
}

func TimestampToString(t int64) string {
	return TimeToString(time.UnixMilli(t))
}

func GetStringTimeDiffer(start_time, end_time string) time.Duration {
	start, _ := time.Parse("2006-01-02 15:04:05", start_time)
	end, _ := time.Parse("2006-01-02 15:04:05", end_time)
	return end.Sub(start)
}

func GetTimeDiffer(start_time, end_time time.Time) time.Duration {
	return end_time.Sub(start_time)
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func PointerTimeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return TimeToString(*t)
}

func StringToTime(s string) *time.Time {
	layouts := []string{
		TimeFormatLayout_Y_m_D_H_M_S,
		TimeFormatLayout_Y_m_D_H_M_S_CST,
		TimeFormatLayout_YmDHMS,
		TimeFormatLayout_YmDhMS,
		TimeFormatLayout_Y_m_D,
		TimeFormatLayout_YmD,
		TimeFormatLayout_Y_m_D_H_M,
		TimeFormatLayout_H_M,
		TimeFormatLayout_H_M_S,
	}

	s = strings.TrimSpace(s)
	// 智能识别时间文本类型转换成 time
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return &t
		}
	}
	return nil
}

func TimeToTimestamp(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.Unix()
}

// TimeSetDefaultLoc 设置time为默认时区
func TimeSetDefaultLoc(timeSuk *time.Time) time.Time {
	var loc, _ = time.LoadLocation(TimeDefaultLocation)
	return timeSuk.In(loc)
}

// TimeFormat 格式化时间
func TimeFormat(timeSuk *time.Time, defaultFormat ...string) string {
	if timeSuk == nil {
		return ""
	}
	format := TimeFormatLayout_Y_m_D_H_M_S
	for _, value := range defaultFormat {
		format = value
	}
	return timeSuk.Format(format)
}

// TimeParseString 格式化时间为 time.Time
func TimeParseString(timeStr string, defaultFormat ...string) (time.Time, error) {
	format := TimeFormatLayout_YmDHMS
	for _, value := range defaultFormat {
		format = value
	}
	var cstSh, _ = time.LoadLocation(TimeDefaultLocation)
	timeUnix, err := time.ParseInLocation(format, timeStr, cstSh)
	if err != nil {
		return time.Time{}, err
	}
	return timeUnix, nil
}

// TimeCurrentTimeString 获取当前时间并转换成自定格式
func TimeCurrentTimeString(defaultFormat ...string) string {
	format := TimeFormatLayout_Y_m_D_H_M_S
	for _, value := range defaultFormat {
		format = value
	}
	loc, _ := time.LoadLocation(TimeDefaultLocation)
	return time.Now().In(loc).Format(format)
}

// TimeRestNextDawn 获取now对应到凌晨的时间(到明天凌晨零点的时间)
func TimeRestNextDawn(now time.Time) time.Duration {
	nextDay := now.AddDate(0, 0, 1)
	//明天凌晨
	nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
	return nextDay.Sub(now)
}

// TimeCurrentTimePointer 获取now时间的 *time.Time 格式
func TimeCurrentTimePointer() *time.Time {
	now := time.Now()
	now = TimeSetDefaultLoc(&now)
	return &now
}

// PointerTimeToTime 将 *time.Time 转换成 time.Time
func PointerTimeToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
