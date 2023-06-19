package handler

import (
	"Agent/exporter"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

func BlackboxProcess() {
	data := yamlToMap()
	TimeControl(data)
}

//讀取Yaml檔轉成map
func yamlToMap() (data map[string]interface{}) {
	yamlFile, err := ioutil.ReadFile("snmp.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	return data
}

//解析map並做分析
func mapResolve(data map[string]interface{}) {
	scrapeConfigs, ok := data["scrape_configs"].([]interface{})
	if !ok {
		log.Fatalf("Invalid YAML structure: 'scrape_configs' not found or has incorrect type")
	}

	for _, scrapeConfig := range scrapeConfigs {

		config, ok := scrapeConfig.(map[interface{}]interface{})
		if !ok {
			log.Println("Invalid scrape config found, skipping...")
			continue
		}

		jobName, ok := config["job_name"].(string)
		if !ok {
			log.Println("Invalid job name found, skipping...")
			continue
		}

		scrapeInterval, ok := config["scrape_interval"].(string)
		if !ok {
			log.Println("Invalid scrape interval found, skipping...")
			continue
		}

		metricsPath, ok := config["metrics_path"].(string)
		if !ok {
			log.Println("Invalid metrics path found, skipping...")
			continue
		}

		// fmt.Println("Job Name:", jobName)
		// fmt.Println("Scrape Interval:", scrapeInterval)
		// fmt.Println("Metrics Path:", metricsPath)

		params, ok := config["params"].(map[interface{}]interface{})
		var paramsValue interface{}
		if ok {

			for param, values := range params {
				fmt.Printf("----param:%s, values:%v\n", param, values)
				paramsValue = values
			}
		}

		staticConfigs, ok := config["static_configs"].([]interface{})
		if ok {
			for _, staticConfig := range staticConfigs {
				targetConfig, ok := staticConfig.(map[interface{}]interface{})
				if !ok {
					log.Println("Invalid target config found, skipping...")
					continue
				}

				targets, ok := targetConfig["targets"].([]interface{})
				if !ok {
					log.Println("Invalid targets found, skipping...")
					continue
				}

				labelsRaw, labelsOK := targetConfig["labels"].(map[interface{}]interface{})

				tag, tagOK := targetConfig["tag"].(string)

				for _, target := range targets {
					targetStr, ok := target.(string)
					if !ok {
						log.Println("Invalid target found, skipping...")
						continue
					}
					startTime := time.Now()
					// Perform HTTP probe
					doc := make(map[string]interface{})
					module := paramsValue.([]interface{})[0] //module 初步討論只會有一個，所以寫死為0
					fmt.Println("module.(string): ", module.(string))
					exporter.CheckModule(module.(string), doc, targetStr)

					fmt.Println(targetStr, " 經過時間: ", time.Since(startTime))

					// Write result to OpenSearch, considering labels and tags
					doc["target"] = targetStr
					// doc["result"] = result

					if labelsOK {
						d := make(map[string]interface{})
						for key, value := range labelsRaw {
							strKey, ok := key.(string)
							if !ok {
								log.Println("Invalid labelsNew type found, skipping...")
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
						log.Println(123, err)
					}

					fmt.Println("Json doc: ", string(r))

					// -----TODO ----- Opensearch Insert

					// // Write doc to OpenSearch
					// // TODO: Implement OpenSearch write operation
					// err = WriteToOpenSearch(doc)
					// if err != nil {
					// 	log.Printf("Failed to write result to OpenSearch: %v", err)
					// 	continue
					// }

					doc = nil

				}

				fmt.Println("---------------------")
			}
		}
	}
}

//建立定時器
func TimeControl(data map[string]interface{}) {
	scrapeConfigs, ok := data["scrape_configs"].([]interface{})
	if !ok {
		log.Fatalf("Invalid YAML structure: 'scrape_configs' not found or has incorrect type")
	}
	count := 0
	for _, scrapeConfig := range scrapeConfigs {
		count++
		log.Println("count: ", count)
		config, ok := scrapeConfig.(map[interface{}]interface{})
		if !ok {
			log.Println("Invalid scrape config found, skipping...")
			continue
		}

		jobName, ok := config["job_name"].(string)
		if !ok {
			log.Println("Invalid job name found, skipping...")
			continue
		}

		scrapeInterval, ok := config["scrape_interval"].(string)
		if !ok {
			log.Printf("Failed to parse scrape_interval for job '%s'", jobName)
			continue
		}

		timeControl, err := time.ParseDuration(scrapeInterval)
		if err != nil {
			log.Printf("Failed to parse scrape_interval for job '%s': %v", jobName, err)
			continue
		}

		go func(config map[interface{}]interface{}) {

			// 建立定時器，定期執行工作
			ticker := time.NewTicker(timeControl)
			defer ticker.Stop()

			for range ticker.C {
				//定時任務觸發
				dataResolve(config)
			}
		}(config)
	}

	// 防止主程式退出
	select {}

}

func dataResolve(config map[interface{}]interface{}) {

	jobName, ok := config["job_name"].(string)
	if !ok {
		log.Println("Invalid job name found, skipping...")
		return
	}

	scrapeInterval, ok := config["scrape_interval"].(string)
	if !ok {
		log.Println("Invalid scrape interval found, skipping...")
		return
	}

	metricsPath, ok := config["metrics_path"].(string)
	if !ok {
		log.Println("Invalid metrics path found, skipping...")
		return
	}

	// fmt.Println("Job Name:", jobName)
	// fmt.Println("Scrape Interval:", scrapeInterval)
	// fmt.Println("Metrics Path:", metricsPath)

	params, ok := config["params"].(map[interface{}]interface{})
	var paramsValue interface{}
	if ok {

		for param, values := range params {
			fmt.Printf("----param:%s, values:%v\n", param, values)
			paramsValue = values
		}
	}

	staticConfigs, ok := config["static_configs"].([]interface{})
	if ok {
		for _, staticConfig := range staticConfigs {
			targetConfig, ok := staticConfig.(map[interface{}]interface{})
			if !ok {
				log.Println("Invalid target config found, skipping...")
				continue
			}

			targets, ok := targetConfig["targets"].([]interface{})
			if !ok {
				log.Println("Invalid targets found, skipping...")
				continue
			}

			labelsRaw, labelsOK := targetConfig["labels"].(map[interface{}]interface{})

			tag, tagOK := targetConfig["tag"].(string)

			for _, target := range targets {
				targetStr, ok := target.(string)
				if !ok {
					log.Println("Invalid target found, skipping...")
					continue
				}
				startTime := time.Now()
				// Perform HTTP probe
				doc := make(map[string]interface{})
				module := paramsValue.([]interface{})[0] //module 初步討論只會有一個，所以寫死為0
				fmt.Println("module.(string): ", module.(string))
				exporter.CheckModule(module.(string), doc, targetStr)

				fmt.Println(targetStr, " 經過時間: ", time.Since(startTime))

				// Write result to OpenSearch, considering labels and tags
				doc["target"] = targetStr
				// doc["result"] = result

				if labelsOK {
					d := make(map[string]interface{})
					for key, value := range labelsRaw {
						strKey, ok := key.(string)
						if !ok {
							log.Println("Invalid labelsNew type found, skipping...")
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
					log.Println(123, err)
				}

				fmt.Println("Json doc: ", string(r))

				// -----TODO ----- Opensearch Insert

				// // Write doc to OpenSearch
				// // TODO: Implement OpenSearch write operation
				// err = WriteToOpenSearch(doc)
				// if err != nil {
				// 	log.Printf("Failed to write result to OpenSearch: %v", err)
				// 	continue
				// }

				doc = nil

			}

			fmt.Println("---------------------")
		}
	}

}
