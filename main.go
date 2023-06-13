package main

import (
	"Agent/model"
	"Agent/opensearch"
	"Agent/pkg"
	"Agent/server"
	"fmt"
	"log"
	"time"
)

const osIndex = "lex-test66"

func bulkPrevious() {
	var data model.BulkPrevious
	// data.Delete = map[string]string{
	// 	"lex-test66": "EJazs4gBM-XHgcOmDej3",
	// }

	if err := opensearch.BulkPrevious("delete", data); err != nil {
		log.Println(err)
	}

}

func main() {
	pkg.AutoLoader()
	bulkPrevious()
	return

	a := model.InsertData{
		Data: map[string]interface{}{"host": "10.40.192.277"},
	}

	id := "8ZY8s4gBM-XHgcOmCOcK"

	result, err := opensearch.BulkUpdate(osIndex, id, a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	if err := opensearch.BulkExecute(result); err != nil {
		fmt.Println(err)
	}

}

func search() {
	os, ok := server.GetServerInstance().GetOpensearch()["One"]
	if !ok {
		fmt.Println("No OK")
		return
	}
	//Search key
	result, err := opensearch.Search(os, "lex-test", "*", "8080")
	if err != nil {
		fmt.Println("No OK")
		return
	}

	log.Println(result.Hits.Hits[0].Source)
}

func create() {
	a := map[string]interface{}{
		"host": "10.40.192.213",
		"http": map[string]interface{}{
			"method":  "POST",
			"request": 1669,
			"version": "HTTP/1.1",
		},
		"url": map[string]interface{}{
			"domain": "10.11.233.11",
			"path":   "/",
			"port":   8080,
		},
		"timestamp": time.Now(),
	}

	result, err := opensearch.BulkCreate(osIndex, a)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("1 result: ", result)

}

func delete() {
	os, ok := server.GetServerInstance().GetOpensearch()["One"]
	if !ok {
		fmt.Println("No OK")
		return
	}

	deleteIndex := []string{"lex-test66"}
	opensearch.SingleDeleteIndex(os, deleteIndex)
}
