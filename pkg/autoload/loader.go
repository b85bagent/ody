package autoload

import (
	"agent/handler"
	"agent/pkg/tool"
	"agent/server"
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"sync"

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
	reload := make(chan bool, 1)
	newReload := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())

	go handler.BlackboxProcess(ctx, targetFile, blackboxFile)
	var reloadMutex sync.Mutex
	go func() {
		for {
			select {
			case <-reload:
				reloadMutex.Lock()
				cancel()
				log.Println("启动新的handler.BlackboxProcess")
				ctx2, cancel := context.WithCancel(context.Background())
				handler.BlackboxProcess(ctx2, "target999.yaml", blackboxFile)

				defer cancel()
				reloadMutex.Unlock()
			}
			reloadMutex.Lock()
			reload = newReload
			reloadMutex.Unlock()
		}
	}()

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		reloadMutex.Lock()
		// 当接收到GET请求时，发送一个信号到reload管道
		log.Println("当接收到GET请求时，发送一个信号到reload管道")

		reload <- true
		w.Write([]byte("Reload signal sent!"))
		reloadMutex.Unlock()
	})

	http.ListenAndServe(":8080", nil)

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
