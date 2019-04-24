package apitest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

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

	profile := a.GetActiveProfile()

	requestURL := a.baseURL + "/" + path
	log.Printf("REQUEST (Profile %d: %s): POST %s \n>> %s", profile.id, profile.name, requestURL, string(bodyBytes))
	request, e = http.NewRequest("POST", requestURL, bytes.NewBuffer(bodyBytes))
	if e != nil {
		a.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if authenticated {
		request.Header.Set("Authorization", profile.authKey)
	}
	client := &http.Client{}
	response, e = client.Do(request)

	// assert.NotNil(t, e)
	if e != nil {
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

	requestURL := a.baseURL + "/" + path

	profile := a.GetActiveProfile()
	log.Printf("REQUEST (Profile %d: %s): PUT %s \n>> %s", profile.id, profile.name, requestURL, string(bodyBytes))
	request, e = http.NewRequest("PUT", requestURL, bytes.NewBuffer(bodyBytes))
	if e != nil {
		a.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if authenticated {
		request.Header.Set("Authorization", profile.authKey)
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

	requestURL := a.baseURL + "/" + path

	profile := a.GetActiveProfile()
	log.Printf("REQUEST (Profile %d: %s): GET %s", profile.id, profile.name, requestURL)

	request, e = http.NewRequest("GET", requestURL, nil)
	if e != nil {
		a.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if withAuth {
		request.Header.Set("Authorization", profile.authKey)
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

	requestURL := a.baseURL + "/" + path
	profile := a.GetActiveProfile()
	log.Printf("REQUEST (Profile %d: %s): DELETE %s", profile.id, profile.name, requestURL)
	request, e = http.NewRequest("DELETE", requestURL, nil)
	if e != nil {
		a.t.Fatal(e)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	if withAuth {
		request.Header.Set("Authorization", profile.authKey)
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
