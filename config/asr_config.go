package config

// ASRConfig ASR服务配置
type ASRConfig struct {
	Vosk   VoskConfig `yaml:"vosk" json:"vosk"`
	Aliyun AliyunAsr  `yaml:"aliyun" json:"aliyun"`
}

// VoskConfig Vosk配置
type VoskConfig struct {
	Model string `yaml:"model" json:"model"`
}

type AliyunAsr struct {
	ApiKey string `yaml:"api_key" json:"api_key"`
}
