package opensearch

import "Agent/model"


func DataMix(data []interface{}, Action *model.Action, ContentDetail *model.ContentDetail) []interface{} {
	data = append(data, Action)
	data = append(data, ContentDetail)

	return data
}

func ActionCreate(index string) *model.Action {
	return &model.Action{Create: &model.CreateDetail{Index: index}}
}

func ContentDetailCreate(title, director, year string) *model.ContentDetail {
	return &model.ContentDetail{Title: title, Director: director, Year: year}
}
