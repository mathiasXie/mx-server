package config

// TTSConfig TTS服务配置
type TTSConfig struct {
	Microsoft  MicrosoftConfig  `yaml:"microsoft"`
	VolcEngine VolcEngineConfig `yaml:"volcengine"`
	Aliyun     AliyunTTSConfig  `yaml:"aliyun"`
}

// MicrosoftConfig 微软TTS配置
type MicrosoftConfig struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
}

// VolcengineConfig 豆包TTS配置
type VolcEngineConfig struct {
	APIID    string `yaml:"api_id"`
	Token    string `yaml:"token"`
	Endpoint string `yaml:"endpoint"`
}

// AliyunTTSConfig 阿里云TTS配置
type AliyunTTSConfig struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
}
