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

// RandString returns a random string based on a set of runes
func RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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
func (a *APITest) SetStringVal(name string, v string) {
	a.stringVals[name] = v
}

// SetIntVal sets an int value identified by `name`
func (a *APITest) SetIntVal(name string, v int64) {
	a.intVals[name] = v
}

// GetStringVal gets the string value identified by `name`
func (a *APITest) GetStringVal(name string) string {
	val, ok := a.stringVals[name]
	if !ok {
		return ""
	}

	return val
}

// GetIntVal gets the int value identified by `name`
func (a *APITest) GetIntVal(name string) int64 {
	val, ok := a.intVals[name]
	if !ok {
		return -1
	}

	return val
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
		request.Header.Set("Authorization", "Bearer "+a.authKey)
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
		request.Header.Set("Authorization", "Bearer "+a.authKey)
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
		request.Header.Set("Authorization", "Bearer "+a.authKey)
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
