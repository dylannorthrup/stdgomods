package stdgomods

import (
	"fmt"
	"time"
)

type durationTimer struct {
	Start time.Time
}

// Ideally this would be a map or array so we could have multiple timers.
// TODO: Do that thing mentioned in the previous line
var (
	Start time.Time // When we want to start our timer
	End   time.Time // Track if we togged debug output or not in ToggleDebugIfNeeded()
)

func (d durationTimer) Elapsed() string {
	return fmt.Sprintf("%v", time.Since(d.Start).Truncate(time.Second).String())
}

func (d durationTimer) PrintElapsed(prefix ...string) {
	fmt.Printf("%s%v has elapsed since this operation began.\n", prefix[0], time.Since(d.Start).Truncate(time.Second).String())
}

func (d durationTimer) Reset() {
	d.Start = time.Now()
	fmt.Printf("Reset DurationTimer to %v.\n", d.Start)
}

func NewDurationTimer() *durationTimer {
	d := new(durationTimer)
	d.Start = time.Now()
	return d
}
