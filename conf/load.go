package conf

import (
	"os"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
)

var config *Config

func C() *Config {
	//默认配置
	if config == nil {
		config = Default()
	}
	return config
}

//加载配置，把外部配置读到config全局变量里面来
//yaml文件 --> conf

func LoadConfigFromYaml(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	//默认值
	config = C()
	return yaml.Unmarshal(content, config)
}

// 从环境变量读取配置
func LoadConfigFromEnv() error {
	config = C()
	return env.Parse(config)
}
