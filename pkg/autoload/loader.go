package autoload

import (
	"agent/handler"
	"agent/pkg/tool"
	"agent/server"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/opensearch-project/opensearch-go"
)

func AutoLoader(configFile, targetFile, blackboxFile string) {

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

	if len(config.Opensearch) > 0 {
		logger.Println("Auto loading opensearch")
		opensearch, err := initOpensearch(config.Opensearch)
		if err != nil {
			log.Println(err)
			panic("initOpensearch fail")
		}

		handlerServer.ServerStruct.SetOpensearch(opensearch)
	}

	logger.Println("AutoLoader Success")

	handler.BlackboxProcess(targetFile, blackboxFile)

	// server.SetServerInstance(serverInstance)

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
