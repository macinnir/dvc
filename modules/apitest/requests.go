package apitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

// Requests sends network requests using `profile` as context
type Requests struct {
	profile *UserProfile
	logger  *Logger
	t       *testing.T
}

// NewRequests returns a new Requests instance
func NewRequests(profile *UserProfile, logger *Logger, t *testing.T) *Requests {
	return &Requests{profile, logger, t}
}

// Post does a post request
func (r *Requests) Post(path string, body interface{}, authenticated bool) (responseBody []byte, statusCode int) {

	var bodyBytes []byte
	var e error

	if bodyBytes, e = json.Marshal(body); e != nil {
		r.logger.Error(e.Error())
		return
	}

	var request *http.Request
	var response *http.Response

	requestURL := r.profile.baseURL + "/" + path

	r.logger.Info(fmt.Sprintf("REQUEST (Profile %d: %s): POST %s", r.profile.ID, r.profile.Name, requestURL))
	r.logger.Debug(fmt.Sprintf("%s", string(bodyBytes)))

	request, e = http.NewRequest("POST", requestURL, bytes.NewBuffer(bodyBytes))
	if e != nil {
		r.logger.Error(e.Error())
		r.t.Fatal(e.Error())
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if authenticated {
		r.logger.Debug(fmt.Sprintf("AuthKey: %s", r.profile.GetString(AuthKey)))
		request.Header.Set("Authorization", r.profile.GetString(AuthKey))
	}
	client := &http.Client{}
	response, e = client.Do(request)

	// assert.NotNil(t, e)
	if e != nil {
		r.logger.Error(e.Error())
		r.t.Fatal(e.Error())
		return
	}

	statusCode = response.StatusCode

	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	r.logger.Info(fmt.Sprintf("RESPONSE: %d", statusCode))
	r.logger.Debug(fmt.Sprintf("%s", responseBody))

	return
}

// PostDebug does a post request while temporarily setting the log level to debug
func (r *Requests) PostDebug(path string, body interface{}, authenticated bool) (responseBody []byte, statusCode int) {
	oldLogLevel := r.logger.logLevel
	r.logger.logLevel = LogLevelDebug
	responseBody, statusCode = r.Post(path, body, authenticated)
	r.logger.logLevel = oldLogLevel
	return
}

// Put does a PUT request
func (r *Requests) Put(path string, body interface{}, authenticated bool) (responseBody []byte, statusCode int) {

	var bodyBytes []byte
	var e error

	if bodyBytes, e = json.Marshal(body); e != nil {
		r.t.Fatal(e)
		return
	}

	var request *http.Request
	var response *http.Response

	requestURL := r.profile.baseURL + "/" + path

	r.logger.Info(fmt.Sprintf("REQUEST (Profile %d: %s): PUT %s", r.profile.ID, r.profile.Name, requestURL))
	r.logger.Debug(fmt.Sprintf("%s", string(bodyBytes)))

	request, e = http.NewRequest("PUT", requestURL, bytes.NewBuffer(bodyBytes))
	if e != nil {
		r.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if authenticated {
		r.logger.Debug(fmt.Sprintf("AuthKey: %s", r.profile.GetString(AuthKey)))
		request.Header.Set("Authorization", r.profile.GetString(AuthKey))
	}
	client := &http.Client{}
	response, e = client.Do(request)
	statusCode = response.StatusCode

	// assert.NotNil(t, e)
	if e != nil {
		r.logger.Info(response.Status)
		r.t.Fatal(e)
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	r.logger.Info(fmt.Sprintf("RESPONSE: %d", statusCode))
	r.logger.Debug(fmt.Sprintf("%s", responseBody))
	return
}

// PutDebug does a put request while temporarily setting the log level to debug
func (r *Requests) PutDebug(path string, body interface{}, authenticated bool) (responseBody []byte, statusCode int) {
	oldLogLevel := r.logger.logLevel
	r.logger.logLevel = LogLevelDebug
	responseBody, statusCode = r.Put(path, body, authenticated)
	r.logger.logLevel = oldLogLevel
	return
}

// Get does a GET request
func (r *Requests) Get(path string, withAuth bool) (responseBody []byte, statusCode int) {

	var request *http.Request
	var response *http.Response
	var e error

	requestURL := r.profile.baseURL + "/" + path

	r.logger.Info(fmt.Sprintf("REQUEST (Profile %d: %s): GET %s", r.profile.ID, r.profile.Name, requestURL))
	request, e = http.NewRequest("GET", requestURL, nil)
	if e != nil {
		r.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if withAuth {
		r.logger.Debug(fmt.Sprintf("AuthKey: %s", r.profile.GetString(AuthKey)))
		request.Header.Set("Authorization", r.profile.GetString(AuthKey))
	}
	client := &http.Client{}
	response, e = client.Do(request)
	statusCode = response.StatusCode

	// assert.NotNil(t, e)
	if e != nil {
		r.logger.Error(response.Status)
		r.t.Fatal(e)
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	r.logger.Info(fmt.Sprintf("RESPONSE: %d", statusCode))
	r.logger.Debug(fmt.Sprintf("%s", responseBody))

	return
}

// GetDebug does a get request while temporarily setting the log level to debug
func (r *Requests) GetDebug(path string, withAuth bool) (responseBody []byte, statusCode int) {
	oldLogLevel := r.logger.logLevel
	r.logger.logLevel = LogLevelDebug
	responseBody, statusCode = r.Get(path, withAuth)
	r.logger.logLevel = oldLogLevel
	return
}

// Delete does a DELETE request
func (r *Requests) Delete(path string, withAuth bool) (responseBody []byte, statusCode int) {

	var request *http.Request
	var response *http.Response
	var e error

	requestURL := r.profile.baseURL + "/" + path
	r.logger.Info(fmt.Sprintf("REQUEST (Profile %d: %s): DELETE %s", r.profile.ID, r.profile.Name, requestURL))
	request, e = http.NewRequest("DELETE", requestURL, nil)
	if e != nil {
		r.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if withAuth {
		r.logger.Debug(fmt.Sprintf("AuthKey: %s", r.profile.GetString(AuthKey)))
		request.Header.Set("Authorization", r.profile.GetString(AuthKey))
	}
	client := &http.Client{}
	response, e = client.Do(request)
	statusCode = response.StatusCode

	// assert.NotNil(t, e)
	if e != nil {
		log.Println(response.Status)
		r.t.Fatal(e)
		return
	}
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	r.logger.Info(fmt.Sprintf("RESPONSE: %d", statusCode))
	r.logger.Debug(fmt.Sprintf("%s", responseBody))

	return
}

// DeleteDebug does a DELETE request while temporarily setting the log level to debug
func (r *Requests) DeleteDebug(path string, withAuth bool) (responseBody []byte, statusCode int) {

	oldLogLevel := r.logger.logLevel
	r.logger.logLevel = LogLevelDebug
	responseBody, statusCode = r.Delete(path, withAuth)
	r.logger.logLevel = oldLogLevel
	return

}
