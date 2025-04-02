package config

// BaseConfig 基础配置
type BaseConfig struct {
	AppName string `yaml:"app_name"`
	RunMode string `yaml:"run_mode"`
}

// Instance 全局配置实例
var Instance = &Config{}

// Config 总配置结构
type Config struct {
	BaseConfig
	Server            ServerConfig   `yaml:"server"`
	Resource          ResourceConfig `yaml:"resource"`
	Log               LoggerConfig   `yaml:"log"`
	AccessTokenSecret string         `yaml:"access_token_secret"`
	TTS               TTSConfig      `yaml:"tts"`
}
