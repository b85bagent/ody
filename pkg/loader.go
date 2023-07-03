package pkg

import (
	"Agent/handler"
	"Agent/server"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/opensearch-project/opensearch-go"
)

func AutoLoader(configFile,snmpFile string) {
	config, err := configInit(configFile)

	if err != nil {
		log.Println(err)
	}

	serverInstance, err := server.NewServer()
	if err != nil {
		log.Println("NewServer:", err)
		panic("autoload fail")
	}

	serverInstance.Constant = config.Const

	handlerServer := &handler.Server{
		ServerStruct: serverInstance,
	}

	if len(config.Opensearch) > 0 {
		log.Println("Auto loading opensearch")
		opensearch, err := initOpensearch(config.Opensearch)
		if err != nil {
			log.Println(err)
		}

		handlerServer.ServerStruct.SetOpensearch(opensearch)
	}

	handler.BlackboxProcess(snmpFile)

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
