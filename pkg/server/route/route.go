package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/erik-sostenes/goauth2lib/pkg/server/response"
)

// Constants that represent the http methods
const (
	GET    METHOD = "GET"
	POST   METHOD = "POST"
	PUT    METHOD = "PUT"
	DELETE METHOD = "DELETE"
	PATCH  METHOD = "PATCH"
)

type (
	// METHOD is a string type that represents an http method
	METHOD string
	// MiddlewareFunc is a type of function that represents a decorator for http handlers
	MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc
	// The route is a map containing the method key and value of an http controller
	Route map[METHOD]http.HandlerFunc
)

// RouteCollection implements the ServerHTTP method of the Handler interface to customize the http request
func (route *RouteCollection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	index := strings.LastIndex(path, "/")
	rootPath := path[:index]
	subPath := path[index:]

	if (*route)[rootPath] == nil || (*route)[rootPath][subPath] == nil {
		response.JSON(w, http.StatusNotFound, response.Response{
			Message: fmt.Sprintf("the route was not found %s ", path),
		})
		return
	}

	handlerFunc := (*route)[rootPath][subPath][METHOD(r.Method)]

	if handlerFunc == nil {
		response.JSON(w, http.StatusMethodNotAllowed, response.Response{
			Message: fmt.Sprintf("method %s is not allowed", r.Method),
		})
		return
	}

	handlerFunc(w, r)
}

type (
	// RouteCollection is a map containing all http handlers, implementing the Handler interface
	RouteCollection map[string]map[string]Route
	RouteGroup      struct {
		RootPath           string
		RouteCollection    RouteCollection
		DefaultMiddlewares []MiddlewareFunc
	}
)

func NewGroup(rootPath string, m ...MiddlewareFunc) *RouteGroup {
	return &RouteGroup{
		RootPath:           rootPath,
		RouteCollection:    make(map[string]map[string]Route),
		DefaultMiddlewares: m,
	}
}

// GET represents a route with its http handler, you will only be able to make http requests of type GET
func (r *RouteGroup) GET(subPath string, handler http.HandlerFunc, m ...MiddlewareFunc) {
	r.RouteCollectionExists()
	r.RouteCollection[r.RootPath][subPath] = r.Route(GET, handler, m...)
}

// POST represents a route with its http handler, you will only be able to make http requests of type POST
func (r *RouteGroup) POST(subPath string, handler http.HandlerFunc, m ...MiddlewareFunc) {
	r.RouteCollectionExists()
	r.RouteCollection[r.RootPath][subPath] = r.Route(GET, handler, m...)
}

// PUT represents a route with its http handler, you will only be able to make http requests of type PUT
func (r *RouteGroup) PUT(subPath string, handler http.HandlerFunc, m ...MiddlewareFunc) {
	r.RouteCollectionExists()
	r.RouteCollection[r.RootPath][subPath] = r.Route(PUT, handler, m...)
}

// DELETE represents a route with its http handler, you will only be able to make http requests of type DELETE
func (r *RouteGroup) DELETE(subPath string, handler http.HandlerFunc, m ...MiddlewareFunc) {
	r.RouteCollectionExists()
	r.RouteCollection[r.RootPath][subPath] = r.Route(DELETE, handler, m...)
}

// PATH represents a route with its http handler, you will only be able to make http requests of type PATH
func (r *RouteGroup) PATH(subPath string, handler http.HandlerFunc, m ...MiddlewareFunc) {
	r.RouteCollectionExists()
	r.RouteCollection[r.RootPath][subPath] = r.Route(POST, handler, m...)
}

// RouteCollectionExists checks if the root path exists if it doesn't an instance is created
func (r *RouteGroup) RouteCollectionExists() {
	if r.RouteCollection[r.RootPath] == nil {
		r.RouteCollection[r.RootPath] = make(map[string]Route)
	}
}

// Route creates a new route to the http handler with its route, http method and middlewares
func (r *RouteGroup) Route(method METHOD, handler http.HandlerFunc, m ...MiddlewareFunc) Route {
	if handler == nil {
		panic(fmt.Sprintf("missing http handler for the root path %s with %s method", r.RootPath, method))
	}

	route := make(Route)
	route[method] = r.applyMiddlewares(handler, m...)
	return route
}

func (r *RouteGroup) applyMiddlewares(handler http.HandlerFunc, m ...MiddlewareFunc) http.HandlerFunc {
	if m == nil {
		return handler
	}

	for _, middleware := range append(r.DefaultMiddlewares, m...) {
		handler = middleware(handler)
	}

	return handler
}
