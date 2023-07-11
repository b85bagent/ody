package handler

import (
	"newProject/server"
)

type Server struct {
	ServerStruct *server.Server
}

func New(s *server.Server) *Server {
	return &Server{
		ServerStruct: s,
	}
}
