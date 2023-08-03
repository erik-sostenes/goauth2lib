package server

import (
	"net/http"

	"github.com/erik-sostenes/goauth2lib/libs/server/route"
)

type (
	HandlerFunc func(http.HandlerFunc) http.HandlerFunc
	Server      struct {
		http.Handler
		route.RouteGroup
	}
)

func New(groups ...route.RouteGroup) *Server {
	s := &Server{}

	routes := make(route.RouteCollection, len(groups))

	for _, group := range groups {
		for key, value := range group.RouteCollection {
			routes[key] = value
		}
	}

	s.Handler = &routes

	return s
}

func (s *Server) Start() error {
	return http.ListenAndServe("localhost:8080", s.Handler)
}
