package opensearch

import (
	"Agent/model"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
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

	var urls = []string{"https://localhost:9200"} // 多urls請用逗號隔開

	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: urls,
		Username:  "admin", // For testing only. Don't store credentials in code.
		Password:  "admin",
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

	fmt.Println(createIndexResponse)
	return nil
}

//單一插入
func SingleInsert(client *opensearch.Client, document string) {
	req := opensearchapi.IndexRequest{
		Index: IndexName,
		Body:  strings.NewReader(document),
	}
	insertResponse, err := req.Do(context.Background(), client.Transport)
	if err != nil {
		fmt.Println("failed to insert document: ", err)
	}

	fmt.Println(insertResponse)
	fmt.Println("add success")
}

//Search something
func Search(client *opensearch.Client, key, value string) (result model.SearchResponse) {

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
		return
	}

	search := opensearchapi.SearchRequest{
		Body: bytes.NewReader(content),
	}

	searchResponse, err := search.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to search document ", err)
		return
	}

	fmt.Println("main: ", searchResponse)

	defer searchResponse.Body.Close()

	json.NewDecoder(searchResponse.Body).Decode(&result)

	return result
}

//Bulk Insert
func BulkInsert(client *opensearch.Client, documents string) bool {

	blk, errBulk := client.Bulk(strings.NewReader(documents))

	if errBulk != nil {
		fmt.Println("failed to perform bulk operations", errBulk)
		return true
	}

	fmt.Println("Performing bulk operations")
	fmt.Println(blk)

	if blk.IsError() {
		var errBulk model.BulkError

		json.NewDecoder(blk.Body).Decode(&errBulk)

		fmt.Println(errBulk.Error.Reason)

		return false
	}

	body, err := io.ReadAll(blk.Body)
	if err != nil {
		log.Printf("error occurred: [%s]", err.Error())
	}

	var response model.BulkCreateResponse
	if errUnmarshal := json.Unmarshal(body, &response); errUnmarshal != nil {
		log.Printf("error Unmarshal blkResponse: [%s]", errUnmarshal.Error())
	}

	for _, item := range response.Items {
		if item.Create.Status > 299 {
			log.Printf("error occurred: [%s]", item.Create.Result)
		} else {
			log.Printf("success: [%s]", item.Create.Index)
		}
	}

	return blk.IsError()
}

//delete Index
func DeleteIndex(client *opensearch.Client, index []string) {
	deleteIndex := opensearchapi.IndicesDeleteRequest{
		Index: index,
	}

	deleteIndexResponse, err := deleteIndex.Do(context.Background(), client.Transport)
	if err != nil {
		fmt.Println("failed to delete index ", err)
		return
	}

	fmt.Println("Deleting the index")
	fmt.Println(deleteIndexResponse)
	defer deleteIndexResponse.Body.Close()

}

