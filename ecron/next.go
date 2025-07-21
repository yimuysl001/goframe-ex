package ecron

import (
	"github.com/gogf/gf/v2/container/gtype"
	_ "github.com/gogf/gf/v2/os/gcron"
	"time"
	_ "unsafe"
)

type cronSchedule struct {
	// todo 可以不要
	createTimestamp int64            // Created timestamp in seconds.
	everySeconds    int64            // Running interval in seconds.
	pattern         string           // The raw cron pattern string that is passed in cron job creation.
	ignoreSeconds   bool             // Mark the pattern is standard 5 parts crontab pattern instead 6 parts pattern.
	secondMap       map[int]struct{} // Job can run in these second numbers.
	minuteMap       map[int]struct{} // Job can run in these minute numbers.
	hourMap         map[int]struct{} // Job can run in these hour numbers.
	dayMap          map[int]struct{} // Job can run in these day numbers.
	weekMap         map[int]struct{} // Job can run in these week numbers.
	monthMap        map[int]struct{} // Job can run in these moth numbers.
	// This field stores the timestamp that meets schedule latest.
	lastMeetTimestamp *gtype.Int64
	// Last timestamp number, for timestamp fix in some latency.
	lastCheckTimestamp *gtype.Int64
}

//go:linkname newSchedule  github.com/gogf/gf/v2/os/gcron.newSchedule
func newSchedule(string) (*cronSchedule, error)

//go:linkname next  github.com/gogf/gf/v2/os/gcron.(*cronSchedule).Next
func next(*cronSchedule, time.Time) time.Time

func GetNext(pattern string, ti time.Time, t ...int) ([]time.Time, error) {
	var ts = 5
	if len(t) > 0 && t[0] > 0 {
		ts = t[0]
	}
	var times = make([]time.Time, ts)
	schedule, err := newSchedule(pattern)
	if err != nil {
		return nil, err
	}

	for i := 0; i < ts; i++ {
		ti = next(schedule, ti)
		times[i] = ti
	}

	return times, nil
}
