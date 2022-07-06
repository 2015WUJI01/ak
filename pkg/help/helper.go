package help

import (
	"fmt"
	"time"
)

// PrintCostTimeSince 输出耗时时长，以毫秒为单位
func PrintCostTimeSince(start time.Time) time.Duration {
	return time.Since(start)
}

func SpendTimeSince(start time.Time) string {
	d := time.Since(start)
	if d.Microseconds() < 1000 {
		return fmt.Sprintf("%d us", d.Microseconds())
	}
	if d.Milliseconds() < 1000 {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	}
	if d.Seconds() < 60 {
		return fmt.Sprintf("%.2f s", d.Seconds())
	}
	return fmt.Sprintf("%.0fm%.0fs", d.Minutes(), d.Seconds())
}

func SpendSeconds(start time.Time) float64 {
	return time.Since(start).Seconds()
}

func SpendMilliseconds(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
