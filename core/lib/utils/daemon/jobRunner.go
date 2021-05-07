package daemon

import (
	"fmt"
	"runtime/debug"
	"time"
)

// RunJobAtInterval runs a job every d duration, but will not start the next job unless the previous job has finished.
func RunJobAtInterval(d time.Duration, f func(int, time.Time)) {
	jobCounter := 0
	for t := range time.Tick(d) {
		runner(jobCounter, t, f)
		jobCounter++
	}
}

func runner(jobCounter int, t time.Time, fn func(int, time.Time)) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Job failed: %d \n %s \n %s", jobCounter, r, string(debug.Stack()))
		}
	}()

	fn(jobCounter, t)
}
