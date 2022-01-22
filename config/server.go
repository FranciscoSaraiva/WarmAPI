package config

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
*	Public Functions
 */

// ServerAPI type that represents the Server hosting the API.
// It overwrites a pointer of httprouter.Router
type ServerAPI struct {
	*httprouter.Router
}

// ServeHTTP function that sets the header content to application/json
// format every time there's a HTTP request coming through.
func (s *ServerAPI) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	s.Router.ServeHTTP(resp, req)
}

// NewServerAPI function that returns a new instance
// of httprouter with httprouter.New().
func NewServerAPI() *ServerAPI {
	return &ServerAPI{
		httprouter.New(),
	}
}
