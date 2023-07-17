package autoload

import (
	"newProject/handler"

	"crypto/tls"
	"log"
	"net/http"
	"newProject/pkg/tool"
	"newProject/server"

	httpserver "newProject/http_server"

	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go"
)

func AutoLoader(configFile string) {

	config, err := configInit(configFile)

	if err != nil {
		log.Println(err)
		panic("config init fail")
	}

	serverInstance, err := server.NewServer()
	if err != nil {
		log.Println(err)
		panic("autoload fail")
	}

	serverInstance.Constant = config.Const

	handlerServer := &handler.Server{
		ServerStruct: serverInstance,
	}

	debugMode := getDebugSetting()

	logger := tool.NewLogger(debugMode)

	handlerServer.ServerStruct.SetLogger(logger)

	if len(config.Opensearch.Opensearch) > 0 {
		handlerServer.ServerStruct.OpensearchIndex = config.Opensearch.Index
		logger.Println("Auto loading opensearch")
		opensearch, err := initOpensearch(config.Opensearch.Opensearch)
		if err != nil {
			log.Println(err)
			panic("initOpensearch fail")
		}

		handlerServer.ServerStruct.SetOpensearch(opensearch)
	}

	// run http
	var httpSrv *http.Server
	httpServerPort, ginOK := serverInstance.Constant["http_server_port"]
	if ginOK {
		addr := ":" + httpServerPort.(string)
		gin.SetMode("release")
		r := gin.New()
		r.Use(gin.Recovery())
		r, initRouterErr := httpserver.InitRouter(r)
		if initRouterErr != nil {
			log.Printf("initRouterErr: %v ", err)
			panic("autoload fail")
		}
		httpSrv = &http.Server{
			Addr:    addr,
			Handler: r,
		}
		go func() {
			if serverRunErr := httpSrv.ListenAndServe(); serverRunErr != nil && serverRunErr != http.ErrServerClosed {
				log.Printf("serverRunErr: %v", serverRunErr)
				panic("autoload fail")
			}
		}()
	} else {
		log.Println("http_server_port is not set")
		return
	}

	logger.Println("AutoLoader Success")

	select {}

}

func initOpensearch(setting map[string]OpensearchConfig) (map[string]*opensearch.Client, error) {

	opensearchClient := make(map[string]*opensearch.Client)

	for key, v := range setting {

		client, err := opensearch.NewClient(opensearch.Config{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Addresses: v.Host,
			Username:  v.Username,
			Password:  v.Password,
		})

		if err != nil {
			log.Println("無法建立 OpenSearch 客戶端:", err)
			return nil, err
		}
		// log.Println(client.Info())
		opensearchClient[key] = client
	}

	// Print OpenSearch version information on console.

	return opensearchClient, nil
}

func getDebugSetting() bool {

	debugSetting, ok := server.GetServerInstance().GetConst()["debug"]
	if !ok {
		log.Println("DebugSetting Get failed")
		return false
	}

	return debugSetting.(bool)
}
