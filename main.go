package main

import (
	"Agent/pkg"
	"Agent/server"
	"fmt"
	"log"
	"time"

	"github.com/b85bagent/opensearch"
)

const osIndex = "lex-test14"

var data opensearch.BulkPreviousUse

func bulkCreate() {

	os, ok := server.GetServerInstance().GetOpensearch()["One"]
	if !ok {
		fmt.Println("No OK")
		return
	}

	data.Create.Data = map[string]interface{}{
		"host": "10.11.22.333",
		"http": map[string]interface{}{
			"method":  "POST",
			"request": 3369,
			"version": "HTTP/1.1",
		},
		"url": map[string]interface{}{
			"domain": "10.42.11.255",
			"path":   "/",
			"port":   8888,
		},
		"timestamp": time.Now(),
	}

	data.Create.Index = osIndex

	result, err := opensearch.BulkPrevious(os, "create", data)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(result)

}

func bulkUpdate() {

	os, ok := server.GetServerInstance().GetOpensearch()["One"]
	if !ok {
		fmt.Println("No OK")
		return
	}

	var data opensearch.BulkPreviousUse
	data.Update.Data = opensearch.InsertData{
		Data: map[string]interface{}{"host": "10.20.30.40"},
	}
	data.Update.Index = osIndex
	data.Update.Id = "LZYIuIgBM-XHgcOmeujd"

	result, err := opensearch.BulkPrevious(os, "update", data)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)

}

func bulkDelete() {

	os, ok := server.GetServerInstance().GetOpensearch()["One"]
	if !ok {
		fmt.Println("No OK")
		return
	}

	data.Delete = map[string]string{
		"lex-test14": "L5YOuIgBM-XHgcOmpujw",
	}

	result, err := opensearch.BulkPrevious(os, "delete", data)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
}

func main() {
	// 載入Server
	pkg.AutoLoader()
	search()
	// bulkCreate()
	// bulkUpdate()
	// bulkDelete()

}

func search() {
	os, ok := server.GetServerInstance().GetOpensearch()["One"]
	if !ok {
		fmt.Println("No OK")
		return
	}
	//Search key
	result, err := opensearch.Search(os, "lex-test14", "", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println(result.Hits)
}
