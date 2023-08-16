package server

import (
	"log"
	"net/http"
	"os"

	"github.com/erik-sostenes/goauth2lib/pkg/server/route"
)

const defaultPort = "8080"

type (
	// Server contains all the settings for the server
	Server struct {
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
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Server is running on port '%s'\n", port)
	return http.ListenAndServe("localhost:"+port, s.Handler)
}
