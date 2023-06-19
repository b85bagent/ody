package model

import (
	"Agent/server"
	"fmt"
	"log"

	os "github.com/b85bagent/opensearch"
	"github.com/opensearch-project/opensearch-go"
)

func DataInsert(data map[string]interface{}, osIndex string) error {
	client := DBinit()
	var Setting os.BulkPreviousUse
	Setting.Create.Data = data
	Setting.Create.Index = osIndex
	result, err := os.BulkPrevious(client, "create", Setting)
	if err != nil {
		log.Println("Bulk Insert error: ", err)
		return err
	}
	log.Println(result)

	return nil
}

func DBinit() *opensearch.Client {

	client, ok := server.GetServerInstance().GetOpensearch()["One"]
	if !ok {
		fmt.Println("Opensearch Client Get failed")
		return nil
	}

	return client
}
