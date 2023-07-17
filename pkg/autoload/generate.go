package autoload

type Config struct {
	Opensearch OpensearchSettings     `yaml:"opensearch"`
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
