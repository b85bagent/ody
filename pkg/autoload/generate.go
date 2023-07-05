package autoload

type Config struct {
	Opensearch map[string]OpensearchConfig `yaml:"opensearch"`
	Const      map[string]interface{}
}

type OpensearchConfig struct {
	Host     []string `yaml:"host"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}
