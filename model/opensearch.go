package model

import (
	"errors"
	"io/ioutil"
	"log"
	"newProject/pkg/tool"
	"newProject/server"
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

	index := server.GetServerInstance().GetOpensearchIndex()

	r.client = client
	r.index = index

	return r
}

//把data 轉成 字串
func DataCompression(data map[string]interface{}, r string) string {

	index := server.GetServerInstance().GetOpensearchIndex()

	result, err := os.BulkCreate(index, data)
	if err != nil {
		log.Println("Bulk Create error: ", err)
		return ""
	}

	r = r + result + "\n"

	return result
}

type Response struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	Items  []struct {
		Create struct {
			Index string `json:"_index"`
			ID    string `json:"_id"`
			// define other fields if needed
		} `json:"create"`
	} `json:"items"`
}

func BulkInsert(data string) error {

	dataInsertMutex.Lock()
	defer dataInsertMutex.Unlock()

	l = server.GetServerInstance().GetLogger()
	client := DBinit()

	result, err := os.BulkExecute(client.client, data)
	if err != nil {
		log.Println("BulkExecute error: ", err)
		return err
	}

	defer result.Body.Close()

	if result.IsError() {
		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			log.Println("result.IsError ReadAll error: ", err)

			return err
		}
		// log.Println("Bulk Insert error: ", result.Body)
		return errors.New(string(body))

	}

	return nil
}
