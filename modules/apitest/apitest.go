package apitest

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// APITest tests your API
type APITest struct {
	// The injected testing object
	t *testing.T
	// The base URL used for all API calls
	baseURL string
	// UserProfiles
	userProfiles map[string]*UserProfile
	// The name of the active UserProfile
	activeProfile string
	// Map of global string values accessible across all profiles
	stringVals map[string]string
	// Map of global int64 values accessible across all profiles
	intVals map[string]int64
	// Map of global object values accessible across all profiles
	objectVals map[string]interface{}
	// Unique key for the current testing session
	sessionKey    string
	logger        *Logger
	storyLevel    int
	storyCounters []int
}

// NewAPITest returns a new APITest instance
func NewAPITest(t *testing.T, baseURL string, logLevel LogLevel) *APITest {

	apiTest := &APITest{
		t:             t,
		baseURL:       baseURL,
		userProfiles:  map[string]*UserProfile{},
		stringVals:    map[string]string{},
		intVals:       map[string]int64{},
		objectVals:    map[string]interface{}{},
		logger:        InitLogger(logLevel),
		storyLevel:    0,
		storyCounters: []int{0},
	}

	apiTest.init()

	return apiTest
}

// Parallel runs parallel tests
// https://gist.github.com/thatisuday/6612075cf0b7e04ef232717f4e8815a3
// https://medium.com/statuscode/pipeline-patterns-in-go-a37bb3a7e61d
func (a *APITest) Parallel(description string, doAsFns ...func()) {
	a.logger.Info(fmt.Sprintf("Parallel: %s", description))
	var wg sync.WaitGroup
	wg.Add(len(doAsFns))
	for _, doAsFn := range doAsFns {
		go func(doAsFn func()) {
			defer wg.Done()
			doAsFn()
		}(doAsFn)
	}
	wg.Wait()
}

func (a *APITest) init() {
	rand.Seed(time.Now().UnixNano())
	a.sessionKey = RandString(10)

	// Create the first user profile
	a.NewProfile(DefaultProfileName)
	a.SetActiveProfile(DefaultProfileName)
}
