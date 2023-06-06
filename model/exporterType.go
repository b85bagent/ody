package model

type Config struct {
	Module  ModuleConfig `yaml:"module"`
	Targets []string     `yaml:"targets"`
	Metrics []Metric     `yaml:"metrics"`
}

type ModuleConfig struct {
	Name    string `yaml:"name"`
	Timeout string `yaml:"timeout"`
	Retries int    `yaml:"retries"`
	Walk    bool   `yaml:"walk"`
	Auth    Auth   `yaml:"auth"`
}

type Auth struct {
	Community string `yaml:"community"`
}

type Metric struct {
	Name     string         `yaml:"name"`
	Oid      string         `yaml:"oid"`
	Index    []string       `yaml:"index,omitempty"`
	CPUUsage CPUUsageConfig `yaml:"cpu_usage,omitempty"`
}

type CPUUsageConfig struct {
	CPU1 float64 `yaml:"cpu1"`
	CPU2 float64 `yaml:"cpu2"`
}

type Data struct {
	Up       float64
	CPUUsage CPUUsageData
}

type CPUUsageData struct {
	Enabled bool
	CPU1    float64
	CPU2    float64
}
