package opensearch

import (
	"Agent/model"
)

func dataMix(data []interface{}, Action any, ContentDetail any) []interface{} {
	data = append(data, Action)
	data = append(data, ContentDetail)

	return data
}

func actionCreate(index string) *model.ActionCreate {
	return &model.ActionCreate{Create: &model.IndexDetail{Index: index}}
}

func contentDetailCreate(data map[string]interface{}) *model.InsertData {
	return &model.InsertData{Data: data}
}

func actionDelete(index, id string) *model.ActionDelete {
	return &model.ActionDelete{Delete: &model.IndexAndIDDetail{Index: index, Id: id}}
}

func actionUpdate(index, id string) *model.ActionUpdate {
	return &model.ActionUpdate{Update: &model.IndexAndIDDetail{Index: index, Id: id}}
}

func contentDetailUpdate(data model.InsertData) *model.UpdateData {
	return &model.UpdateData{Doc: data}
}
