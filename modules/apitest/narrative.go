package apitest

import (
	"fmt"
	"strings"
)

/**
 * Logging Stories
 */

func (a *APITest) incrementStory() {
	a.storyCounters[a.storyLevel]++
}

func (a *APITest) levelUp() {
	if a.storyLevel == 0 {
		return
	}
	a.storyLevel--
	a.logger.UnIndent()

	// for k, c := range a.storyCounters {
	// 	if k <= a.storyLevel {
	// 		continue
	// 	}

	// 	a.storyCounters[]
	// }
}

func (a *APITest) levelDown() {

	a.storyLevel++

	if len(a.storyCounters) < (a.storyLevel + 1) {
		a.storyCounters = append(a.storyCounters, 0)
	}

	a.storyCounters[a.storyLevel] = 0
	a.logger.Indent()
}

func (a *APITest) countTitle() string {
	titleCounts := []string{}
	for k, storyCount := range a.storyCounters {
		if k > a.storyLevel {
			break
		}

		titleCounts = append(titleCounts, fmt.Sprintf("%d", storyCount))
	}

	return strings.Join(titleCounts, ".")
}

// Start starts a new story
func (a *APITest) Start(title string) {
	a.incrementStory()
	a.logger.Log(fmt.Sprintf("%s. %s", a.countTitle(), title))
	a.levelDown()
}

func (a *APITest) Finish() {
	a.levelUp()
}

// Note prints a string to log with an asterisk prefix
func (a *APITest) Note(note string) {
	a.logger.Log(fmt.Sprintf("* %s", note))
}

// Notef call Note() with formatting
func (a *APITest) Notef(note string, args ...interface{}) {
	note = fmt.Sprintf(note, args...)
	a.logger.Log("* " + note)
}
