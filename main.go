package main

func main() {
	
}

func opensearchExample() {
	//----New----
	// client := opensearch.New()

	// fmt.Println("client端完成，準備寫入")

	//----insert----

	// opensearch.SingleInsert(client, `{
	// "title": "Moneyball",
	// "director": "Bennett Miller",
	// "year": "2011"
	// }`)

	// client.Indices.Refresh()

	// fmt.Println("寫入完成，準備Search")

	//----bulk----

	// data := []interface{}{}
	// Action := opensearch.ActionCreate("go-test-index6")

	// ContentDetail1 := opensearch.ContentDetailCreate("Star", "Justin Lin1", "2023")
	// data = opensearch.DataMix(data, Action, ContentDetail1)

	// ContentDetail2 := opensearch.ContentDetailCreate("Trek", "Justin Lin2", "2024")
	// data = opensearch.DataMix(data, Action, ContentDetail2)

	// ContentDetail3 := opensearch.ContentDetailCreate("Beyond", "Justin Lin3", "2025")
	// data = opensearch.DataMix(data, Action, ContentDetail3)

	// buf := &bytes.Buffer{}
	// enc := json.NewEncoder(buf)

	// for _, v := range data {
	// 	if err := enc.Encode(v); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// fmt.Println(buf.String())

	// result := opensearch.BulkInsert(client, buf.String())
	// fmt.Println("是否成功Bulk Insert: ", result)

	//----delete----
	// deleteIndex := []string{"go-test-index6"}
	// opensearch.DeleteIndex(client, deleteIndex)

	//----search----
	// result := opensearch.Search(client, "title", "Star")

	// fmt.Println("搜尋結果: ", result)
}
