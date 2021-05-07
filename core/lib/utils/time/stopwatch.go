package time

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// StopWatch is a multi-step timer for measuring execution time of functions, utils.
type StopWatch struct {
	Total     *StopWatchStep
	timers    []*StopWatchStep
	timerKeys map[string]int
	curStep   *StopWatchStep
}

// NewStopWatch creates a new StopWatch
func NewStopWatch() *StopWatch {
	return &StopWatch{
		Total:     NewStopWatchStep("Total"),
		timers:    []*StopWatchStep{},
		timerKeys: map[string]int{},
	}
}

// Step creates a new step inside stopWatchMgr
func (t *StopWatch) Step(name string) {
	newTimer := NewStopWatchStep(name)
	t.timerKeys[name] = len(t.timers)
	t.timers = append(t.timers, newTimer)
	t.curStep = newTimer
}

// FinishStep finishes the current step (if exists)
func (t *StopWatch) FinishStep() {
	if t.curStep == nil {
		log.Fatal("No current step to finish")
	}

	t.curStep.Finish()
}

// CurrentStep returns the current step
func (t *StopWatch) CurrentStep() *StopWatchStep {
	return t.curStep
}

// Finish finishes the stop watch
func (t *StopWatch) Finish() {
	t.Total.Finish()
}

// Print prints stats to stdout
func (t *StopWatch) Print() {

	// func (t *StopWatchMgr) String() string {

	maxNameLen := 0
	maxDurLen := 0

	for _, timer := range t.timers {

		if len(timer.name) > maxNameLen {
			maxNameLen = len(timer.name)
		}

		if len(timer.du.String()) > maxDurLen {
			maxDurLen = len(timer.du.String())
		}
	}

	maxNameLen = maxNameLen + 2
	fmt.Printf("Total Operations: %d\n", len(t.timers))
	// str += strings.Repeat("-", maxNameLen+maxDurLen+4))

	for _, timer := range t.timers {
		// padding := maxNameLen - len(timer.name)
		fmt.Printf("  %-"+strconv.Itoa(maxNameLen)+"s: %s\n", timer.name, timer.du)
		// fmt.Println(timer.String())
	}

	fmt.Println(strings.Repeat("-", maxNameLen+maxDurLen+4))
	fmt.Println(fmt.Sprintf("  %-"+strconv.Itoa(maxNameLen)+"s: %s", t.Total.name, t.Total.du))
}

// StopWatchStep is a timer container
type StopWatchStep struct {
	name string
	tm   time.Time
	du   time.Duration
}

// NewStopWatchStep returns a new stopWatch step
func NewStopWatchStep(name string) *StopWatchStep {
	t := &StopWatchStep{
		name: name,
	}
	t.Start()
	return t
}

// Start starts the stop watch
func (t *StopWatchStep) Start() {
	t.tm = time.Now()
}

// Finish ends the stop watch
func (t *StopWatchStep) Finish() {
	t.du = time.Since(t.tm)
}

// String() returns a human readable result of the stop watch
func (t *StopWatchStep) String() string {
	return fmt.Sprintf("%s: %s", t.name, t.du)
}

// Milliseconds returns the step time in Milliseconds
func (t *StopWatchStep) Milliseconds() float64 {
	return float64(t.du / time.Millisecond)
}
