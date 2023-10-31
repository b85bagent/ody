package handler

import (
	"encoding/json"
	"log"
	"newProject/server"
	"os"
	"sync"

	"github.com/b85bagent/rabbitmq"
	"github.com/streadway/amqp"
)

var wg sync.WaitGroup

// 開啟rabbitMQ監聽
func ListenRabbitMQ(reload chan bool) error {
	rabbitMQ := server.GetServerInstance().GetRabbitMQArg()

	for _, v := range rabbitMQ.RabbitMQQueue {
		rabbitMQArg := getRabbitMQArg(rabbitMQ, v)
		localResponse := initRPCResponse(rabbitMQArg)
		wg.Add(1)
		go handleRabbitMQMessage(rabbitMQArg, localResponse, reload)
	}

	wg.Wait()
	return nil
}

func getRabbitMQArg(rabbitMQ server.RabbitMQArg, queueName string) rabbitmq.RabbitMQArg {
	return rabbitmq.RabbitMQArg{
		Host:               rabbitMQ.Host[0],
		Username:           rabbitMQ.Username,
		Password:           rabbitMQ.Password,
		RabbitMQExchange:   rabbitMQ.RabbitMQExchange,
		RabbitMQRoutingKey: rabbitMQ.RabbitMQRoutingKey,
		RabbitMQQueue:      queueName,
	}
}

func initRPCResponse(arg rabbitmq.RabbitMQArg) rabbitmq.RPCResponse {
	t := make(map[string]interface{})
	t["message"] = "Agent get MQ message Successfully"
	return rabbitmq.RPCResponse{
		Status:     rabbitmq.Response_Success,
		StatusCode: rabbitmq.Response_Success_Code,
		Response:   t,
		Queue:      arg.RabbitMQQueue,
	}
}

// rabbitmq handle
func handleRabbitMQMessage(arg rabbitmq.RabbitMQArg, response rabbitmq.RPCResponse, reload chan bool) {
	defer wg.Done()

	localResponse := response // 使用本地副本，避免數據競爭

	err := rabbitmq.ListenRabbitMQUsingRPC(arg, localResponse, func(msg amqp.Delivery, ch *amqp.Channel, localResponse rabbitmq.RPCResponse) error {
		//handle and do something here

		// response
		err := replyToPublisher(localResponse, ch, msg)
		if err != nil {
			log.Printf("Error replyToPublisher : %s", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Println("ListenRabbitMQUsingRPC ERROR: ", err)
	}
}

// 將回應發送回去
func replyToPublisher(localResponse rabbitmq.RPCResponse, ch *amqp.Channel, msg amqp.Delivery) error {

	response, err := json.Marshal(localResponse)
	if err != nil {
		log.Printf("Error marshaling to JSON: %v\n", err)
		return err
	}

	// 發送回應到 reply_to 隊列
	errReplay := ch.Publish(
		"",          // exchange
		msg.ReplyTo, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: msg.CorrelationId,
			Body:          response,
		})
	if errReplay != nil {
		log.Printf("Failed to publish a message Replay: %s", errReplay)
		return errReplay
	}

	// 發送 ack 確認消息已經被處理
	err = msg.Ack(false)
	if err != nil {
		log.Printf("Error acknowledging message : %s", err)
		return err
	}

	return nil
}

func SaveYAMLToFile(content []byte, filepath string) error {
	l := server.GetServerInstance().GetLogger()

	err := os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		log.Println("無法寫入檔案：", err)
		return err
	}

	l.Println("已成功儲存 YAML 檔案：", filepath)
	return nil
}
