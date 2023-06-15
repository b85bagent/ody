package main

import (
	"Agent/exporter"
	"Agent/pkg"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

func processYAML(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fmt.Println("Key:", key)
			processYAML(value)
		}
	case []interface{}:
		for _, item := range v {
			processYAML(item)
		}
	default:
		fmt.Println("Value:", v)
	}
}

func main() {
	// 載入Server
	pkg.AutoLoader()

	testHttp()

}

func testHttp() {

	yamlFile, err := ioutil.ReadFile("snmp.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

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

		fmt.Println("Job Name:", jobName)
		fmt.Println("Scrape Interval:", scrapeInterval)
		fmt.Println("Metrics Path:", metricsPath)

		fmt.Println(config["params"].(map[interface{}]interface{}))
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
					exporter.ProbeHttp(doc, targetStr)

					fmt.Println(targetStr, " 經過時間: ", time.Since(startTime))

					// Write result to OpenSearch, considering labels and tags
					doc["target"] = targetStr
					// doc["result"] = result

					if labelsOK {
						fmt.Println("Labels into check")
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

					log.Printf("label: %T tag:%t ", doc["labels"], doc["tag"])

					fmt.Println("doc: ", doc)
					r, err := json.Marshal(doc)
					if err != nil {
						log.Println(123, err)
					}



					fmt.Println("Json doc: ", string(r))

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
