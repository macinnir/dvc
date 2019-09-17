package utils

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type contextKey int

const (
	// ContextUserID key used to fetch the user id from the request context
	ContextUserID contextKey = iota

	// iota
)

// NoAuthRequest represents a request that requires no authentication
type NoAuthRequest struct {
	Method string
	Path   string
}

// NewAuthHandler returns a new AuthHandler object
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		anonymousRoutes: map[string]string{},
	}
}

// AuthHandler validates the incoming Authentication header
type AuthHandler struct {
	anonymousRoutes map[string]string
	authCallback    func(userID int64) (interface{}, error)
}

// AddAnonymousRoute adds an anonymous route to the allowed anonymous routes map
func (a *AuthHandler) AddAnonymousRoute(method string, path string) {

	a.anonymousRoutes[method+"_"+path] = method + "_" + path
}

// IsAnonymousRoute verifies whether or not the route being used is in fact anonymous
func (a *AuthHandler) IsAnonymousRoute(r *http.Request) bool {

	if r.RequestURI == "*" {
		return true
	}

	requestPath := r.RequestURI
	requestPathLen := len(requestPath)

	if requestPathLen == 0 {
		return true
	}

	if strings.Contains(requestPath, "?") {
		requestPathParts := strings.Split(requestPath, "?")
		requestPath = requestPathParts[0]
	}

	fullRouteName := r.Method + "_" + requestPath

	// log.Printf("FullRouteName: %s\n", r.Method+"_"+requestPath)

	_, ok := a.anonymousRoutes[fullRouteName]
	return ok
}

// LogRoute logs a route event to the logger
func (a *AuthHandler) LogRoute(r *http.Request, userID int64, anonymous bool) {

	userIDString := strconv.Itoa(int(userID))
	anonymousString := "AUTH " + userIDString
	if anonymous {
		anonymousString = "ANON " + userIDString
	}
	log.Println("INF HTTP > " + r.Method + " " + r.RequestURI + " [" + anonymousString + "]")
}

// SetAuthCallback sets a callback method to be used for authentication
func (a *AuthHandler) SetAuthCallback(cb func(userID int64) (interface{}, error)) {
	a.authCallback = cb
}

// CreateRouteAuthHandler creates an authentication handler for all routes
func (a *AuthHandler) CreateRouteAuthHandler(h http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var userID int64
		var e error
		var user interface{}

		if a.IsAnonymousRoute(r) {
			a.LogRoute(r, 0, false)
			h.ServeHTTP(w, r)
			return
		}

		if userID, e = GetUserIDFromAuthHeader(r); e != nil {
			a.LogRoute(r, 0, false)
			Unauthorized(r, w)
			return
		}

		if user, e = a.authCallback(userID); e != nil {
			a.LogRoute(r, -1, false)
			Unauthorized(r, w)
			return
		}

		a.LogRoute(r, userID, false)

		// add the user object to the global context for this request
		h.ServeHTTP(
			w,
			r.WithContext(context.WithValue(r.Context(), ContextUserID, user)),
		)
	}
}
