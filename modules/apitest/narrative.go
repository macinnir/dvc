package apitest

import "fmt"

/**
 * Logging Stories
 */

// Story starts and stops a story
func (a *APITest) Story(title string, fn func()) {
	a.StartStory(title)
	fn()
	a.FinishStory()
}

// StartStory starts a new story
func (a *APITest) StartStory(title string) {
	a.storyCount++
	a.logger.Heading(fmt.Sprintf("%d. %s", a.storyCount, title))
}

// FinishStory finishes a new story
func (a *APITest) FinishStory() {
	a.logger.FinishHeading()
	a.subStoryCount = 0
}

func (a *APITest) StartSubStory(title string) {
	a.subStoryCount++
	a.logger.Log(fmt.Sprintf("%d.%d %s", a.storyCount, a.subStoryCount, title))
	a.logger.Indent()
}

func (a *APITest) FinishSubStory() {
	a.logger.UnIndent()
}
