package server

import (
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Run(port string, handler http.Handler) error {
	server.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return server.httpServer.ListenAndServe()
}

func (server *Server) Shutdown(context context.Context) error {
	return server.httpServer.Shutdown(context)
}
