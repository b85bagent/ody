package model

import (
	"agent/pkg/tool"
	"agent/server"
	"log"

	os "github.com/b85bagent/opensearch"
	"github.com/opensearch-project/opensearch-go"
)

type opensearchConfig struct {
	client *opensearch.Client
	index  string
}

var (
	l *tool.Logger
)

func DataInsert(data map[string]interface{}) error {
	l = server.GetServerInstance().GetLogger()
	client := DBinit()

	var Setting os.BulkPreviousUse
	Setting.Create.Data = data
	Setting.Create.Index = client.index

	result, err := os.BulkPrevious(client.client, "create", Setting)
	if err != nil {
		log.Println("Bulk Insert error: ", err)
		return err
	}

	l.Println(result)

	return nil
}

func DBinit() (r opensearchConfig) {

	client, ok := server.GetServerInstance().GetOpensearch()["One"]
	if !ok {
		log.Println("Opensearch Client Get failed")
		return r
	}

	index, ok := server.GetServerInstance().GetConst()["index"]
	if !ok {
		log.Println("timeoutSetting Get failed")
		return r
	}

	r.client = client
	r.index = index.(string)

	return r
}
