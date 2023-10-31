package autoload

import (
	"context"
	"newProject/handler"
	"newProject/model"
	"os"
	"strings"
	"sync"
	"time"

	"crypto/tls"
	"log"
	"net/http"
	"newProject/pkg/tool"
	"newProject/server"

	httpserver "newProject/http_server"

	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go"
)

var wg sync.WaitGroup

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

	insertInterval := server.GetServerInstance().Constant["insert_interval"].(int)

	flushInterval := time.Duration(insertInterval) * time.Second

	bufferSize := server.GetServerInstance().Constant["bufferSize"].(int)

	HandlerServer := handler.New(serverInstance)

	debugMode := getDebugSetting()

	logger := tool.NewLogger(debugMode)

	HandlerServer.ServerStruct.SetLogger(logger)

	enableCheck := false

	if len(config.Opensearch.Opensearch) > 0 {
		HandlerServer.ServerStruct.OpensearchIndex = config.Opensearch.Index
		logger.Println("Auto loading opensearch")
		opensearch, err := initOpensearch(config.Opensearch.Opensearch)
		if err != nil {
			log.Println(err)
			panic("initOpensearch fail")
		}

		HandlerServer.ServerStruct.SetOpensearch(opensearch)
		enableCheck = true
	}

	if len(config.RabbitMQ.RabbitMQ) > 0 {

		for _, v := range config.RabbitMQ.RabbitMQ {

			if v.Enable {
				handlerServer.ServerStruct.RabbitMQConfig.Host = v.Host
				handlerServer.ServerStruct.RabbitMQConfig.Username = v.Username
				handlerServer.ServerStruct.RabbitMQConfig.Password = v.Password
				handlerServer.ServerStruct.RabbitMQConfig.RabbitMQExchange = v.RabbitMQExchange
				handlerServer.ServerStruct.RabbitMQConfig.RabbitMQRoutingKey = v.RabbitMQRoutingKey
				handlerServer.ServerStruct.RabbitMQConfig.Enable = v.Enable
				logger.Println("Auto loading RabbitMQ")
				enableCheck = true
			}

		}

	}

	// run http
	var httpSrv *http.Server
	httpServerPort, ginOK := serverInstance.Constant["http_server_port"]
	if ginOK {
		addr := ":" + httpServerPort.(string)
		gin.SetMode("release")
		r := gin.New()
		r.Use(gin.Recovery())
		r, initRouterErr := httpserver.InitRouter(r, HandlerServer)
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

	if !enableCheck {
		log.Println("沒有執行任何操作設定，請確認config 內enable部分是否有設定 true")
		os.Exit(0)
	}

	// Main context
	ctx := tool.WaitShutdown(func() {})
	handlerServer.ServerStruct.SetGracefulCtx(&ctx)

	logger.Println("AutoLoader Success")

	// background
	bgCtx, bgCancel := context.WithCancel(ctx)
	defer bgCancel()

	wg.Add(1) // 加1
	go func() {
		defer wg.Done() // 當 goroutine 結束時減1
		handlerServer.Background(bgCtx)
	}()

	go func() {
		ticker := time.NewTicker(flushInterval)
		defer ticker.Stop()

		buffer := make([]string, 0, bufferSize)

		for {
			select {
			case item := <-HandlerServer.BufferChan:
				logger.Println("收到item 寫入buffer")
				buffer = append(buffer, item)

				if len(buffer) >= bufferSize {
					log.Println("***buffer > bufferSize ，提前開始執行寫入清空***")

					err := model.BulkInsert(strings.Join(buffer, "\n"))
					if err != nil {
						log.Println("Error performing bulk insert: ", err)
					}

					buffer = nil // clear the batch
				}

			case <-ticker.C:
				if len(buffer) > 0 {
					err := model.BulkInsert(strings.Join(buffer, "\n"))
					if err != nil {
						log.Println("Error performing bulk insert: ", err)
					}
					logger.Println("*********", len(buffer), "筆資料入opensearch 寫入成功*********")

					buffer = nil // clear the batch
				}
			}
		}
	}()

	select {
	case s := <-ctx.Done():
		logger.Printf("shutdownObserver:", s)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	wg.Wait()

	var countdownTime = 5
	for t := countdownTime; t > 0; t-- {
		log.Printf("%d秒後退出", t)
		time.Sleep(time.Second * 1)
	}

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
