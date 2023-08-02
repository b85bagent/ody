package autoload

import (
	"errors"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

func configInit(configFile string) (*Config, error) {

	match, err := regexp.MatchString("^config.*\\.*", configFile)
	if err != nil {
		e := errors.New("config regexp error : " + err.Error())
		return nil, e
	}

	if !match {
		e := errors.New("config 檔案名稱不符合要求，請用-h 確認指令以及符合的yaml格式")
		return nil, e
	}

	filePath := "./yaml/"

	configFile = filePath + configFile

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// 定义一个Config类型的变量来存储解析后的配置信息
	var config Config

	// 解析配置文件
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil

}
