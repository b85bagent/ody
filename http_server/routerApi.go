package http_server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
)

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func write(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// 處理錯誤
		return
	}

	var req prompb.WriteRequest

	err = jsonpb.UnmarshalString(string(body), &req)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	for _, ts := range req.Timeseries {
		labels := make(map[string]string)
		for _, label := range ts.Labels {
			labels[label.Name] = label.Value
		}

		// 將需要寫入的資料轉換為 JSON 格式
		data := map[string]interface{}{
			"labels":  labels,
			"samples": ts.Samples,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(string(jsonData))

	}

	defer c.Request.Body.Close()

	c.JSON(200, gin.H{
		"message": "remote write success",
	})
}

func write2(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// 處理錯誤
		return
	}

	var req prompb.WriteRequest

	err = jsonpb.UnmarshalString(string(body), &req)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	for _, ts := range req.Timeseries {
		m := make(model.Metric, len(ts.Labels))

		for _, l := range ts.Labels {
			m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
		}
		log.Println(m)

		for _, s := range ts.Samples {
			log.Printf("\tSample: %f %d\n", s.Value, s.Timestamp)
		}

		for _, e := range ts.Exemplars {
			m := make(model.Metric, len(e.Labels))
			for _, l := range e.Labels {
				m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
			}
			log.Printf("\tExemplar: %+v %f %d\n", m, e.Value, e.Timestamp)
		}

	}

	defer c.Request.Body.Close()

	c.JSON(200, gin.H{
		"message": "remote write success",
	})
}
