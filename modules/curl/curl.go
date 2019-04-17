package curl

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

// Curl handles network requests
type Curl struct {
	url        string
	authKey    string
	t          *testing.T
	stringVals map[string]string
	intVals    map[string]int64
}

// NewCurl returns a new Curl instance
func NewCurl(url string) *Curl {
	return &Curl{
		url: url,
	}
}

// SetAuthKey sets the auth key
func (a *Curl) SetAuthKey(authKey string) {
	log.Printf("Setting auth key to %s", authKey)
	a.authKey = authKey
}

// GetAuthKey returns the authentication key
func (a *Curl) GetAuthKey() string {
	return a.authKey
}

// Post does a post request
func (a *Curl) Post(path string, body interface{}, authenticated bool) (responseBody []byte, statusCode int) {

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
		return
	}

	statusCode = response.StatusCode

	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	log.Printf("RESPONSE: %d %s", statusCode, responseBody)

	return
}

// Put does a PUT request
func (a *Curl) Put(path string, body interface{}, authenticated bool) (responseBody []byte, statusCode int) {

	var bodyBytes []byte
	var e error

	if bodyBytes, e = json.Marshal(body); e != nil {
		return
	}

	var request *http.Request
	var response *http.Response

	requestURL := a.url + "/" + path
	log.Printf("REQUEST: PUT %s >> %s", requestURL, string(bodyBytes))
	request, e = http.NewRequest("PUT", requestURL, bytes.NewBuffer(bodyBytes))
	if e != nil {
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
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	log.Printf("RESPONSE: %d %s", statusCode, responseBody)
	return
}

// Get does a GET request
func (a *Curl) Get(path string, withAuth bool) (responseBody []byte, statusCode int) {

	var request *http.Request
	var response *http.Response
	var e error

	requestURL := a.url + "/" + path
	log.Printf("REQUEST: GET %s", requestURL)
	request, e = http.NewRequest("GET", requestURL, nil)
	if e != nil {
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
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	log.Printf("RESPONSE: %d %s", statusCode, responseBody)

	return
}

// Delete does a DELETE request
func (a *Curl) Delete(path string, withAuth bool) (responseBody []byte, statusCode int) {

	var request *http.Request
	var response *http.Response
	var e error

	requestURL := a.url + "/" + path
	log.Printf("REQUEST: DELETE %s", requestURL)
	request, e = http.NewRequest("DELETE", requestURL, nil)
	if e != nil {
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
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	log.Printf("RESPONSE: %d %s", statusCode, responseBody)

	return
}
