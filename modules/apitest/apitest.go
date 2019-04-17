package apitest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz_1234567890")

// APITest tests your API
type APITest struct {
	url        string
	authKey    string
	t          *testing.T
	stringVals map[string]string
	intVals    map[string]int64
}

// NewAPITest returns a new APITest instance
func NewAPITest(t *testing.T, url string) *APITest {
	rand.Seed(time.Now().UnixNano())
	return &APITest{
		url:        url,
		t:          t,
		stringVals: map[string]string{},
		intVals:    map[string]int64{},
	}
}

// RandString returns a random string based on a set of runes
func (a *APITest) RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// SetRandString generates a random string of length `length`,
// sets it to the string cache and returns the string
func (a *APITest) SetRandString(name string, length int) string {
	randString := a.RandString(length)
	a.SetString(name, randString)
	return randString
}

// RandInt generates a random number where 0 <= n <= max
func (a *APITest) RandInt(max int) int64 {
	return int64(rand.Intn(max))
}

// SetRandInt generates a random number where 0 <= n <= max,
// sets it to the number cache and returns the number
func (a *APITest) SetRandInt(name string, max int) int64 {
	randNum := a.RandInt(max)
	a.SetInt(name, randNum)
	return randNum
}

// SetAuthKey sets the auth key
func (a *APITest) SetAuthKey(authKey string) {
	log.Printf("Setting auth key to %s", authKey)
	a.authKey = authKey
}

// GetAuthKey returns the authentication key
func (a *APITest) GetAuthKey() string {
	return a.authKey
}

// SetStringVal sets a string val identified by `name`
func (a *APITest) SetString(name string, v string) {
	a.stringVals[name] = v
}

// GetString gets the string value identified by `name`
func (a *APITest) GetString(name string) string {
	val, ok := a.stringVals[name]
	if !ok {
		return ""
	}

	return val
}

// SetInt sets an int value identified by `name`
func (a *APITest) SetInt(name string, v int64) {
	a.intVals[name] = v
}

// GetInt gets the int value identified by `name`
func (a *APITest) GetInt(name string) int64 {
	val, ok := a.intVals[name]
	if !ok {
		return -1
	}

	return val
}

// Increment sets an integer in the number cache if it doesn't exist
// And then increments it by 1
// The resulting value is returned
func (a *APITest) Increment(name string) int64 {
	val, ok := a.intVals[name]
	if !ok {
		a.intVals[name] = 1
	} else {
		a.intVals[name] = val + 1
	}

	return a.intVals[name]
}

// Decrement sets an integer in the number cache if it doesn't exist
// And then decrements it by 1
// The resulting value is returned
func (a *APITest) Decrement(name string) int64 {
	val, ok := a.intVals[name]
	if !ok {
		a.intVals[name] = 0
	} else {
		a.intVals[name] = val - 1
	}

	return a.intVals[name]
}

// Post does a post request
func (a *APITest) Post(path string, body interface{}, authenticated bool) (responseBody []byte, statusCode int) {

	var bodyBytes []byte
	var e error

	if bodyBytes, e = json.Marshal(body); e != nil {
		a.t.Fatal(e)
		return
	}

	var request *http.Request
	var response *http.Response

	requestURL := a.url + "/" + path
	log.Printf("REQUEST: POST %s >> %s", requestURL, string(bodyBytes))
	request, e = http.NewRequest("POST", requestURL, bytes.NewBuffer(bodyBytes))
	if e != nil {
		a.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if authenticated {
		request.Header.Set("Authorization", a.authKey)
	}
	client := &http.Client{}
	response, e = client.Do(request)

	// assert.NotNil(t, e)
	if e != nil {
		log.Println("#### Hello?")
		log.Println(e.Error())
		// log.Println(response.Status)
		a.t.Fatal(e.Error())
		return
	}

	statusCode = response.StatusCode

	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	log.Printf("RESPONSE: %d %s", statusCode, responseBody)

	return
}

// Put does a PUT request
func (a *APITest) Put(path string, body interface{}, authenticated bool) (responseBody []byte, statusCode int) {

	var bodyBytes []byte
	var e error

	if bodyBytes, e = json.Marshal(body); e != nil {
		a.t.Fatal(e)
		return
	}

	var request *http.Request
	var response *http.Response

	requestURL := a.url + "/" + path
	log.Printf("REQUEST: PUT %s >> %s", requestURL, string(bodyBytes))
	request, e = http.NewRequest("PUT", requestURL, bytes.NewBuffer(bodyBytes))
	if e != nil {
		a.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if authenticated {
		request.Header.Set("Authorization", a.authKey)
	}
	client := &http.Client{}
	response, e = client.Do(request)
	statusCode = response.StatusCode

	// assert.NotNil(t, e)
	if e != nil {
		log.Println(response.Status)
		a.t.Fatal(e)
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	log.Printf("RESPONSE: %d %s", statusCode, responseBody)
	return
}

// Get does a GET request
func (a *APITest) Get(path string, withAuth bool) (responseBody []byte, statusCode int) {

	var request *http.Request
	var response *http.Response
	var e error

	requestURL := a.url + "/" + path
	log.Printf("REQUEST: GET %s", requestURL)
	request, e = http.NewRequest("GET", requestURL, nil)
	if e != nil {
		a.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if withAuth {
		request.Header.Set("Authorization", a.authKey)
	}
	client := &http.Client{}
	response, e = client.Do(request)
	statusCode = response.StatusCode

	// assert.NotNil(t, e)
	if e != nil {
		log.Println(response.Status)
		a.t.Fatal(e)
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	log.Printf("RESPONSE: %d %s", statusCode, responseBody)

	return
}

// Delete does a DELETE request
func (a *APITest) Delete(path string, withAuth bool) (responseBody []byte, statusCode int) {

	var request *http.Request
	var response *http.Response
	var e error

	requestURL := a.url + "/" + path
	log.Printf("REQUEST: DELETE %s", requestURL)
	request, e = http.NewRequest("DELETE", requestURL, nil)
	if e != nil {
		a.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if withAuth {
		request.Header.Set("Authorization", a.authKey)
	}
	client := &http.Client{}
	response, e = client.Do(request)
	statusCode = response.StatusCode

	// assert.NotNil(t, e)
	if e != nil {
		log.Println(response.Status)
		a.t.Fatal(e)
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	log.Printf("RESPONSE: %d %s", statusCode, responseBody)

	return
}
