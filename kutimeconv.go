package kutimeconv

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// parseUptime parses the uptime from /proc/uptime
func parseUptime(buf []byte) (uptime uint64, err error) {
	uptimeStr := strings.Split(string(buf), " ")[0]
	uptimeSeconds, err := strconv.ParseFloat(uptimeStr, 64)
	if err != nil {
		return 0, err
	}
	return uint64(uptimeSeconds * 1e9), nil
}

// kernel uptime in nanoseconds
func GetKernelUptime() (uptime uint64, err error) {
	buf, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	uptime, err = parseUptime(buf)
	if err != nil {
		return 0, err
	}
	return uptime, nil
}

// UptimeToTime converts uptime to time
func UptimeToTime(uptime uint64) (t time.Time, err error) {
	nowUptime, err := GetKernelUptime()
	if err != nil {
		return t, err
	}
	return uptimeToTime(nowUptime, uptime, time.Now())
}

// uptimeをtime.Timeに変換
func uptimeToTime(nowUptime, uptime uint64, now time.Time) (t time.Time, err error) {
	uptimeDiff := uptimeDiff(nowUptime, uptime)
	if err != nil {
		return t, err
	}

	// 現在時刻からuptimeDiffを引いた時刻を返す
	return now.Add(uptimeDiff), nil
}

// uptimeの差を計算
func uptimeDiff(uptime1, uptime2 uint64) time.Duration {
	return time.Duration(-(int64(uptime1) - int64(uptime2)))
}

// TimeToUptime converts time to uptime
func TimeToUptime(t time.Time) (uptime uint64, err error) {
	nowUptime, err := GetKernelUptime()
	if err != nil {
		return 0, err
	}
	return timeToUptime(t, time.Now(), nowUptime)
}

// time.Timeをuptimeに変換
func timeToUptime(t, now time.Time, nowUptime uint64) (uptime uint64, err error) {
	uptimeDiff := t.Sub(now)

	// uptimeDurationをnanosecondsに変換
	diff := uptimeDiff.Nanoseconds()

	// 現在のuptimeからtを引いたuptimeを返す
	return uint64(int64(nowUptime) + diff), nil
}
