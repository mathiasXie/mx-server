package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig 加载配置文件
func (c *Config) LoadConfig(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}
