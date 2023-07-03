package model

import (
	"Agent/server"
	"log"

	os "github.com/b85bagent/opensearch"
	"github.com/opensearch-project/opensearch-go"
)

type opensearchConfig struct {
	client *opensearch.Client
	index  string
}

func DataInsert(data map[string]interface{}) error {
	client := DBinit()
	log.Println(client)

	var Setting os.BulkPreviousUse
	Setting.Create.Data = data
	Setting.Create.Index = client.index

	result, err := os.BulkPrevious(client.client, "create", Setting)
	if err != nil {
		log.Println("Bulk Insert error: ", err)
		return err
	}

	log.Println(result)

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
