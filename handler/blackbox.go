package handler

import (
	"agent/exporter"
	"agent/model"
	"agent/pkg/tool"
	"agent/server"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	bec "agent/blackbox_exporter/config"

	logger "github.com/go-kit/log"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

var (
	sc = bec.NewSafeConfig(prometheus.DefaultRegisterer)
	l  *tool.Logger
)

//blackbox 進程
func BlackboxProcess(targetFile, blackboxFile string) {

	//讀取blackbox.yaml
	sc, err := blackboxConfig(blackboxFile)
	if err != nil {
		log.Printf("讀取blackbox配置文件錯誤: %v ，請用-h 確認指令以及符合的yaml格式", err)
		panic("blackbox config init fail")
	}

	//讀取target.yaml
	targetConfig, err := targetConfig(targetFile)
	if err != nil {
		log.Printf("讀取target配置文件錯誤: %v", err)
		panic("target config init fail")
	}

	//定時器設定
	TimeControl(targetConfig, sc)
}

//讀取Target Yaml檔轉成map
func targetConfig(targetFile string) (data map[string]interface{}, err error) {
	match, err := regexp.MatchString("^target.*\\.*", targetFile)
	if err != nil {
		e := errors.New("target regexp error : " + err.Error())
		return nil, e
	}

	if !match {
		e := errors.New("target 檔案名稱不符合要求，請用-h 確認指令以及符合的yaml格式")
		return nil, e
	}

	yamlFile, errReadFile := ioutil.ReadFile(targetFile)
	if errReadFile != nil {
		return nil, errReadFile
	}

	errUnmarshal := yaml.Unmarshal(yamlFile, &data)
	if errUnmarshal != nil {
		return nil, errReadFile
	}

	return data, nil
}

//建立定時器
func TimeControl(data map[string]interface{}, sc *bec.SafeConfig) {

	l = server.GetServerInstance().GetLogger()

	scrapeConfigs, ok := data["scrape_configs"].([]interface{})
	if !ok {
		l.Println("Invalid YAML structure: 'scrape_configs' not found or has incorrect type")
		return
	}

	for _, scrapeConfig := range scrapeConfigs {

		config, ok := scrapeConfig.(map[interface{}]interface{})
		if !ok {
			l.Println("Invalid scrape config found, skipping...")
			continue
		}

		jobName, ok := config["job_name"].(string)
		if !ok {
			l.Println("Invalid job name found, skipping...")
			continue
		}

		scrapeInterval, ok := config["scrape_interval"].(string)
		if !ok {
			l.Printf("Failed to parse scrape_interval for job '%s'", jobName)
			continue
		}

		timeControl, err := time.ParseDuration(scrapeInterval)
		if err != nil {
			l.Printf("Failed to parse scrape_interval for job '%s': %v", jobName, err)
			continue
		}

		go func(config map[interface{}]interface{}, sc *bec.SafeConfig) {

			//優先執行一次
			dataResolve(config, sc)

			// 建立定時器，定期執行工作
			ticker := time.NewTicker(timeControl)
			defer ticker.Stop()

			for range ticker.C {
				//定時任務觸發
				dataResolve(config, sc)

			}
		}(config, sc)
	}

	// 防止主程式退出
	select {}

}

//解析yaml檔後做probe
func dataResolve(config map[interface{}]interface{}, sc *bec.SafeConfig) {
	l := server.GetServerInstance().GetLogger()

	jobName, ok := config["job_name"].(string)
	if !ok {
		l.Println("Invalid job name found, skipping...")
		return
	}

	scrapeInterval, ok := config["scrape_interval"].(string)
	if !ok {
		l.Println("Invalid scrape interval found, skipping...")
		return
	}

	metricsPath, ok := config["metrics_path"].(string)
	if !ok {
		l.Println("Invalid metrics path found, skipping...")
		return
	}

	// l.Println("Job Name:", jobName)
	// l.Println("Scrape Interval:", scrapeInterval)
	// l.Println("Metrics Path:", metricsPath)

	params, ok := config["params"].(map[interface{}]interface{})
	var paramsValue interface{}
	if ok {

		for param, values := range params {
			l.Printf("----param:%s, values:%v\n", param, values)
			paramsValue = values
		}
	}

	staticConfigs, ok := config["static_configs"].([]interface{})
	if ok {
		for i, staticConfig := range staticConfigs {
			targetConfig, ok := staticConfig.(map[interface{}]interface{})
			if !ok {
				l.Println("Invalid target config found, skipping...")
				continue
			}

			targets, ok := targetConfig["targets"].([]interface{})
			if !ok {
				l.Println("Invalid targets found, skipping...")
				continue
			}

			labelsRaw, labelsOK := targetConfig["labels"].(map[interface{}]interface{})

			tag, tagOK := targetConfig["tag"].(string)

			for _, target := range targets {

				targetStr, ok := target.(string)
				if !ok {
					l.Println("Invalid target found, skipping...")
					continue
				}

				startTime := time.Now()

				doc := make(map[string]interface{})

				module := paramsValue.([]interface{})[0] //module 初步討論只會有一個，所以寫死為0

				_, errCMADP := exporter.CheckModuleAndDoProbe(module.(string), doc, targetStr, sc)
				if errCMADP != nil {
					l.Printf("第 %d 個CheckModuleAndDoProbe failed: %e", i, errCMADP)
					continue
				}

				l.Println("target: ", targetStr, "的 Process 經過時間: ", time.Since(startTime))

				// Write result to OpenSearch, considering labels and tags
				doc["target"] = targetStr

				if labelsOK {
					d := make(map[string]interface{})
					for key, value := range labelsRaw {
						strKey, ok := key.(string)
						if !ok {
							l.Println("Invalid labelsNew type found, skipping...")
							continue
						}
						d[strKey] = value
					}
					doc["labels"] = d
				}

				if tagOK {
					doc["tag"] = tag
				}

				doc["jobName"] = jobName
				doc["params"] = paramsValue
				doc["scrape_interval"] = scrapeInterval
				doc["metrics_path"] = metricsPath

				/*---如果想看資料(json格式)，請把下面註解移除---
				r, err := json.Marshal(doc)
				if err != nil {
					l.Println("json Marshal error: ", err)
				}

				l.Println("Json doc: ", string(r))
				*/

				if doc["result"] == "Failed" {
					l.Printf("target: %s Failed", targetStr)
					continue
				}

				// Write doc to OpenSearch
				if errInsertOS := model.DataInsert(doc); errInsertOS != nil {
					l.Printf("Error Bulk Insert, Job_Name: %s, target :%s, reason :%e", jobName, targetStr, errInsertOS)
					doc = nil
					continue
				}

				l.Println("寫入openSearch成功")

				doc = nil //reset map

			}

		}
	}
	log.Printf("Job '%s' 工作已完成", jobName)

}

//test use
func blackboxConfig(blackboxFile string) (*bec.SafeConfig, error) {

	match, err := regexp.MatchString("^blackbox.*\\.*", blackboxFile)
	if err != nil {
		e := errors.New("blackbox regexp error : " + err.Error())
		return nil, e
	}

	if !match {
		e := errors.New("blackbox 檔案名稱不符合要求")
		return nil, e
	}

	logger := logger.NewNopLogger()

	location := "./blackbox_exporter/" + blackboxFile

	if err := sc.ReloadConfig(location, logger); err != nil {
		level.Error(logger).Log("msg", "Error reloading config", "err", err)
		return nil, err
	}

	return sc, nil
}

//解析map並做分析
func mapResolve(data map[string]interface{}, sc *bec.SafeConfig) {
	l := server.GetServerInstance().GetLogger()

	scrapeConfigs, ok := data["scrape_configs"].([]interface{})
	if !ok {
		log.Fatalf("Invalid YAML structure: 'scrape_configs' not found or has incorrect type")
	}

	for _, scrapeConfig := range scrapeConfigs {

		config, ok := scrapeConfig.(map[interface{}]interface{})
		if !ok {
			l.Println("Invalid scrape config found, skipping...")
			continue
		}

		jobName, ok := config["job_name"].(string)
		if !ok {
			l.Println("Invalid job name found, skipping...")
			continue
		}

		scrapeInterval, ok := config["scrape_interval"].(string)
		if !ok {
			l.Println("Invalid scrape interval found, skipping...")
			continue
		}

		metricsPath, ok := config["metrics_path"].(string)
		if !ok {
			l.Println("Invalid metrics path found, skipping...")
			continue
		}

		// l.Println("Job Name:", jobName)
		// l.Println("Scrape Interval:", scrapeInterval)
		// l.Println("Metrics Path:", metricsPath)

		params, ok := config["params"].(map[interface{}]interface{})
		var paramsValue interface{}
		if ok {

			for param, values := range params {
				l.Printf("----param:%s, values:%v\n", param, values)
				paramsValue = values
			}
		}

		staticConfigs, ok := config["static_configs"].([]interface{})
		if ok {
			for _, staticConfig := range staticConfigs {
				targetConfig, ok := staticConfig.(map[interface{}]interface{})
				if !ok {
					l.Println("Invalid target config found, skipping...")
					continue
				}

				targets, ok := targetConfig["targets"].([]interface{})
				if !ok {
					l.Println("Invalid targets found, skipping...")
					continue
				}

				labelsRaw, labelsOK := targetConfig["labels"].(map[interface{}]interface{})

				tag, tagOK := targetConfig["tag"].(string)

				for _, target := range targets {
					targetStr, ok := target.(string)
					if !ok {
						l.Println("Invalid target found, skipping...")
						continue
					}
					startTime := time.Now()
					// Perform HTTP probe
					doc := make(map[string]interface{})
					module := paramsValue.([]interface{})[0] //module 初步討論只會有一個，所以寫死為0

					exporter.CheckModuleAndDoProbe(module.(string), doc, targetStr, sc)

					l.Println("target: ", targetStr, "的 Process 經過時間: ", time.Since(startTime))

					// Write result to OpenSearch, considering labels and tags
					doc["target"] = targetStr
					// doc["result"] = result

					if labelsOK {
						d := make(map[string]interface{})
						for key, value := range labelsRaw {
							strKey, ok := key.(string)
							if !ok {
								l.Println("Invalid labelsNew type found, skipping...")
								continue
							}
							d[strKey] = value
						}
						doc["labels"] = d
					}

					if tagOK {
						doc["tag"] = tag
					}

					doc["jobName"] = jobName
					doc["params"] = paramsValue
					doc["scrape_interval"] = scrapeInterval
					doc["metrics_path"] = metricsPath

					r, err := json.Marshal(doc)
					if err != nil {
						l.Println(123, err)
					}

					l.Println("Json doc: ", string(r))

					// -----TODO ----- Opensearch Insert

					// // Write doc to OpenSearch
					// // TODO: Implement OpenSearch write operation
					// err = WriteToOpenSearch(doc)
					// if err != nil {
					// 	l.Printf("Failed to write result to OpenSearch: %v", err)
					// 	continue
					// }

					doc = nil

				}

				l.Println("---------------------")
			}
		}
	}
}
