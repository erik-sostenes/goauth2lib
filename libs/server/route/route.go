package route

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	GET    METHOD = "GET"
	POST   METHOD = "POST"
	PUT    METHOD = "PUT"
	DELETE METHOD = "DELETE"
	PATCH  METHOD = "PATCH"
)

type (
	METHOD         string
	MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc
	Route          map[METHOD]http.HandlerFunc
)

func (route *RouteCollection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	index := strings.LastIndex(path, "/")
	rootPath := path[:index]
	subPath := path[index:]

	handlerFunc := (*route)[rootPath][subPath][METHOD(r.Method)]

	if handlerFunc == nil {
		http.Error(w, fmt.Sprintf("method %s is not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	handlerFunc(w, r)
}

type (
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

func (r *RouteGroup) GET(subPath string, handler http.HandlerFunc, m ...MiddlewareFunc) {
	r.RouteCollectionExists()
	r.RouteCollection[r.RootPath][subPath] = r.Route(GET, handler, m...)
}

func (r *RouteGroup) POST(subPath string, handler http.HandlerFunc, m ...MiddlewareFunc) {
	r.RouteCollectionExists()
	r.RouteCollection[r.RootPath][subPath] = r.Route(POST, handler, m...)
}

func (r *RouteGroup) RouteCollectionExists() {
	if r.RouteCollection[r.RootPath] == nil {
		r.RouteCollection[r.RootPath] = make(map[string]Route)
	}
}

func (r *RouteGroup) Route(method METHOD, handler http.HandlerFunc, m ...MiddlewareFunc) Route {
	if handler == nil {
		panic("")
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
