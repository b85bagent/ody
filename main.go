package main

import (
	RabbitMQ "Agent/handler"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// exporter.RunExporter()
	content, err := ReadYAMLFromFile("snmp.yaml")
	if err != nil {
		log.Println("Read Yaml error: ", err)
	}

	rabbitmq := RabbitMQ.NewRabbitMQSimple("lex")
	rabbitmq.PublishSimple(content)
	
	fmt.Println("send success")

}

func ReadYAMLFromFile(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println("無法讀取檔案：", err)
		return "", err
	}
	return string(content), nil
}
