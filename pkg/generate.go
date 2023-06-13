package pkg

type Config struct {
	Opensearch map[string]OpensearchConfig `yaml:"opensearch"`
}

type OpensearchConfig struct {
	Host     []string `yaml:"host"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}
