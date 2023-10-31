package handler

import (
	"newProject/server"
)

type Server struct {
	ServerStruct   *server.Server
	BufferChan     chan string

}

func New(s *server.Server) *Server {

	bufferSize := server.GetServerInstance().Constant["bufferSize"].(int)

	return &Server{
		ServerStruct: s,
		BufferChan:   make(chan string, bufferSize),
	}
}

