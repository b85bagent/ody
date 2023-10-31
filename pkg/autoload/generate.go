package autoload

type Config struct {
	Opensearch OpensearchSettings     `yaml:"opensearch"`
	RabbitMQ   RabbitMQSettings       `yaml:"rabbitMQ"`
	Const      map[string]interface{} `yaml:"const"`
}

type OpensearchSettings struct {
	Index      string                      `yaml:"index"`
	Opensearch map[string]OpensearchConfig `yaml:",inline"`
}

type OpensearchConfig struct {
	Host     []string `yaml:"host"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

type RabbitMQSettings struct {
	RabbitMQ map[string]RabbitMQConfig `yaml:",inline"`
}

type RabbitMQConfig struct {
	Host               []string `yaml:"host"`
	Username           string   `yaml:"username"`
	Password           string   `yaml:"password"`
	RabbitMQExchange   string   `yaml:"RabbitMQExchange"`
	RabbitMQRoutingKey string   `yaml:"RabbitMQRoutingKey"`
	RabbitMQQueue      []string `yaml:"RabbitMQQueue"`
	Enable             bool     `yaml:"enable"`
}
