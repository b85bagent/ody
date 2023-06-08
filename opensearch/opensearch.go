package opensearch

import (
	"Agent/model"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

const IndexName = "go-test-index1"

//建立opensearch Client端
func New() *opensearch.Client {

	var urls = []string{"https://10.11.233.102:9200"} // 多urls請用逗號隔開

	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: urls,
		Username:  "admin",
		Password:  "systex123!",
	})

	if err != nil {
		fmt.Println("無法建立 OpenSearch 客戶端:", err)
		return nil
	}

	// Print OpenSearch version information on console.
	fmt.Println(client.Info())

	return client
}

//建立Index
func CreateIndex(client *opensearch.Client, IndexName string) error {
	//設定Index
	settings := strings.NewReader(`{
	'settings': {
	    'index': {
	        'number_of_shards': 1,
	        'number_of_replicas': 2
	        }
	    }
	}`)

	res := opensearchapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  settings,
	}

	createIndexResponse, errCreateIndex := res.Do(context.Background(), client)
	if errCreateIndex != nil {
		fmt.Println("failed to create index ", errCreateIndex)
		return errCreateIndex
	}
	defer createIndexResponse.Body.Close()

	fmt.Println(createIndexResponse)
	return nil
}

//單一插入
func SingleInsert(client *opensearch.Client, document string) error {
	req := opensearchapi.IndexRequest{
		Index: IndexName,
		Body:  strings.NewReader(document),
	}
	insertResponse, err := req.Do(context.Background(), client.Transport)
	if err != nil {
		fmt.Println("failed to insert document: ", err)
		return err
	}
	defer insertResponse.Body.Close()

	fmt.Println(insertResponse)
	fmt.Println("add success")
	return nil
}

//Search something
func Search(client *opensearch.Client, key, value string) (result model.SearchResponse, err error) {

	s := map[string]interface{}{
		// "size": 5,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				key: value,
			},
		},
	}

	content, errMarshal := json.Marshal(s)
	if errMarshal != nil {
		fmt.Println(errMarshal)
		return result, errMarshal
	}

	search := opensearchapi.SearchRequest{
		Body: bytes.NewReader(content),
	}

	searchResponse, errSearch := search.Do(context.Background(), client)
	if errSearch != nil {
		fmt.Println("failed to search document ", errSearch)
		return result, errSearch
	}

	defer searchResponse.Body.Close()

	json.NewDecoder(searchResponse.Body).Decode(&result)

	return result, nil
}

//Bulk Insert
func BulkInsert(client *opensearch.Client, documents string) error {

	blk, errBulk := client.Bulk(strings.NewReader(documents))

	if errBulk != nil {
		fmt.Println("failed to perform bulk operations", errBulk)
		return errBulk
	}
	defer blk.Body.Close()

	fmt.Println("Performing bulk operations")
	fmt.Println(blk)

	if blk.IsError() {
		var errBulk model.BulkError

		json.NewDecoder(blk.Body).Decode(&errBulk)

		errBody := errors.New(errBulk.Error.Reason)
		return errBody
	}

	body, errReadAll := io.ReadAll(blk.Body)
	if errReadAll != nil {
		log.Printf("error occurred: [%s]", errReadAll.Error())
		return errReadAll
	}

	var response model.BulkCreateResponse
	if errUnmarshal := json.Unmarshal(body, &response); errUnmarshal != nil {
		log.Printf("error Unmarshal blkResponse: [%s]", errUnmarshal.Error())
		return errUnmarshal
	}

	for _, item := range response.Items {
		if item.Create.Status > 299 {
			log.Printf("error occurred: [%s]", item.Create.Result)
		} else {
			log.Printf("success: [%s]", item.Create.Index)
		}
	}

	return nil
}

//delete Index
func DeleteIndex(client *opensearch.Client, index []string) error {
	deleteIndex := opensearchapi.IndicesDeleteRequest{
		Index: index,
	}

	deleteIndexResponse, errDeleteIndex := deleteIndex.Do(context.Background(), client.Transport)
	if errDeleteIndex != nil {
		fmt.Println("failed to delete index ", errDeleteIndex)
		return errDeleteIndex
	}
	defer deleteIndexResponse.Body.Close()

	fmt.Println(deleteIndexResponse)

	return nil

}
