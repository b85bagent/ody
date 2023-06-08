
{
content, err := ReadYAMLFromFile("snmp.yaml")
 if err != nil {
  log.Println("Read Yaml error: ", err)
 }

 rabbitmq := RabbitMQ.NewRabbitMQSimple("lex")

 // Push Message to Queue
 rabbitmq.PublishSimple(content)

 fmt.Println("send success")

 // catch the message from Queue
 //rabbitmq.ConsumeSimple()
}

func ReadYAMLFromFile(filepath string) (string, error) {
 content, err := ioutil.ReadFile(filepath)
 if err != nil {
  log.Println("無法讀取檔案：", err)
  return "", err
 }
 return string(content), nil
}
