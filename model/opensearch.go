package model

import (
	"agent/pkg/tool"
	"agent/server"
	"errors"
	"io/ioutil"
	"log"
	"sync"

	os "github.com/b85bagent/opensearch"
	"github.com/opensearch-project/opensearch-go"
)

type opensearchConfig struct {
	client *opensearch.Client
	index  string
}

var (
	l               *tool.Logger
	dataInsertMutex sync.Mutex
)

func DataInsert(data map[string]interface{}) error {
	dataInsertMutex.Lock()
	defer dataInsertMutex.Unlock()

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
		log.Println("index Get failed")
		return r
	}

	r.client = client
	r.index = index.(string)

	return r
}

//把data 轉成 字串
func DataCompression(data map[string]interface{}, r string) string {

	index, ok := server.GetServerInstance().GetConst()["index"]
	if !ok {
		log.Println("index Get failed")
		return ""
	}

	result, err := os.BulkCreate(index.(string), data)
	if err != nil {
		log.Println("Bulk Create error: ", err)
		return ""
	}

	r = r + result + "\n"

	return result
}

func BulkInsert(data string) error {
	dataInsertMutex.Lock()
	defer dataInsertMutex.Unlock()

	l = server.GetServerInstance().GetLogger()
	client := DBinit()

	result, err := os.BulkExecute(client.client, data)
	if err != nil {
		return err
	}

	if result.IsError() {
		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return err
		}
		// log.Println("Bulk Insert error: ", result.Body)
		return errors.New(string(body))

	}

	return nil
}
