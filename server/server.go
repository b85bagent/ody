package server

import (
	"log"
	"remote_write/pkg/tool"
	"sync"

	"github.com/opensearch-project/opensearch-go"
)

var (
	serverObject *Server
	once         sync.Once
)

type Server struct {
	Constant map[string]interface{}
	// redisCacher   map[string]*redis.Cacher
	// goworker      map[string]*goworker.Pool
	// grpcClient    map[string]*grpc.ClientConn
	// httpClient    map[string]httpClient.Methods
	// gracefulCtx   *context.Context
	opensearchClient map[string]*opensearch.Client
	logger           *tool.Logger
}

func NewServer() (newServerObject *Server, err error) {
	initServer()
	newServerObject = serverObject
	return
}

func initServer() {
	serverObject = &Server{}
	serverObject.opensearchClient = make(map[string]*opensearch.Client)
}

func GetServerInstance() *Server {
	if serverObject == nil {
		log.Fatal("Server instance has not been initialized. Please call NewServer first.")
	}
	return serverObject
}

func (s *Server) SetOpensearch(opensearchClient map[string]*opensearch.Client) {
	s.opensearchClient = opensearchClient
}

func (s *Server) GetOpensearch() map[string]*opensearch.Client {
	return s.opensearchClient
}

func (s *Server) GetOpensearchIndex() string {
	return s.GetOpensearchIndex()
}

func (s *Server) SetConst(Const map[string]interface{}) {
	s.Constant = Const
}

func (s *Server) GetConst() (Const map[string]interface{}) {
	return s.Constant
}

func (s *Server) SetLogger(logger *tool.Logger) {
	s.logger = logger
}

func (s *Server) GetLogger() *tool.Logger {
	return s.logger
}
