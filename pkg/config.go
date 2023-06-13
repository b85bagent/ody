package pkg

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func configInit() (*Config, error) {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
		return nil, err
	}

	// 定义一个Config类型的变量来存储解析后的配置信息
	var config Config

	// 解析配置文件
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
		return nil, err
	}

	return &config, nil

}
