package utils

import (
	"context"
	"log"
	"net/http"
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

func (a *AuthHandler) AddAnonymousRoute(method string, path string) {

	a.anonymousRoutes[method+"_"+path] = method + "_" + path
}

func (a *AuthHandler) IsAnonymousRoute(r *http.Request) bool {

	if r.RequestURI == "*" {
		return true
	}

	requestPath := r.RequestURI
	requestPathLen := len(requestPath)

	if requestPathLen == 0 {
		return true
	}

	fullRouteName := r.Method + "_" + r.RequestURI

	_, ok := a.anonymousRoutes[fullRouteName]
	return ok
}

func (a *AuthHandler) LogRoute(r *http.Request, anonymous bool) {

	anonymousString := "Auth"
	if anonymous {
		anonymousString = "Anon"
	}
	log.Println(r.Method + " " + r.RequestURI + " [" + anonymousString + "]")
}

func (a *AuthHandler) SetAuthCallback(cb func(userID int64) (interface{}, error)) {
	a.authCallback = cb
}

func (a *AuthHandler) CreateRouteAuthHandler(h http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var userID int64
		var e error
		var user interface{}

		if a.IsAnonymousRoute(r) {
			a.LogRoute(r, true)
			h.ServeHTTP(w, r)
			return
		}

		if userID, e = GetUserIDFromAuthHeader(r); e != nil {
			a.LogRoute(r, false)
			Unauthorized(w)
			return
		}

		if user, e = a.authCallback(userID); e != nil {
			a.LogRoute(r, false)
			Unauthorized(w)
			return
		}

		a.LogRoute(r, false)

		// add the user object to the global context for this request
		h.ServeHTTP(
			w,
			r.WithContext(context.WithValue(r.Context(), ContextUserID, user)),
		)
	}
}
