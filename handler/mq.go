package handler

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/streadway/amqp"
)

//連接資訊格式：amqp://帳號:密碼@RabbitMQ伺服器地址:Port (預設端口為5672)

const MQUrl = "amqp://lex:s850429s@127.0.0.1:5672/Lex"

//rabbitMQ結構體
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string //隊列名稱
	Exchange  string //交換機名稱
	Key       string //bind Key 名稱
	MQurl     string //連接資訊
}

//創建結構體實例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MQurl: MQUrl}
}

//斷開channel和connection
func (r *RabbitMQ) ReleaseRes() {
	r.channel.Close()
	r.conn.Close()
}

//錯誤處理函數
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//創建簡單模式下RabbitMQ實例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//創建RabbitMQ實例
	rabbitmq := NewRabbitMQ(queueName, "", "")

	var err error

	//創建 connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MQurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+
		"itmq!")

	//創建 channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")

	return rabbitmq
}

//Simple Mode Producer
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申請隊列，如果隊列不存在會自動創建，存在則跳過創建
	_, err := r.channel.QueueDeclare(
		r.QueueName, //隊列名
		false,       //是否持久化
		false,       //是否自動刪除
		false,       //是否具有排他性
		false,       //是否阻塞處理
		nil,         //額外的屬性
	)
	if err != nil {
		log.Println(err)
	}
	//調用channel發送消息到隊列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, //如果為true，根據自身exchange類型和routekey規則無法找到符合條件的隊列會把消息返還給發送者
		false, //如果為true，當exchange發送消息到隊列後發現隊列上沒有消費者，則會把消息返還給發送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//Simple Mode Consumer
func (r *RabbitMQ) ConsumeSimple() {
	//1.申請隊列，如果隊列不存在會自動創建，存在則跳過創建
	q, err := r.channel.QueueDeclare(
		r.QueueName, //隊列名
		false,       //是否持久化
		false,       //是否自動刪除
		false,       //是否具有排他性
		false,       //是否阻塞處理
		nil,         //額外的屬性
	)
	if err != nil {
		log.Println(err)
	}

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer 用來區分多個消費者
		true,   // auto-ack 是否自動應答
		false,  // exclusive 是否獨有
		false,  // no-local 設置為true，表示不能將同一個Connection中生產者發送的消息傳遞給這個Connection中 的消費者
		false,  // no-wait 列是否阻塞
		nil,    // 額外的屬性
	)
	if err != nil {
		log.Println(err)
	}

	forever := make(chan bool)
	//消息邏輯處理
	go func() {
		for d := range msgs {
			//處理yaml to 本地端file
			err := SaveYAMLToFile(d.Body, "./test.yaml")
			if err != nil {
				log.Println("save yaml file error:", err)
			}

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func SaveYAMLToFile(content []byte, filepath string) error {
	err := ioutil.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		log.Println("無法寫入檔案：", err)
		return err
	}
	log.Println("已成功儲存 YAML 檔案：", filepath)
	return nil
}
