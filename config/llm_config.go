package config

// LLMConfig LLM服务配置
type LLMConfig struct {
	VolcEngine OpenAiConfig `yaml:"volcengine"`
	Aliyun     OpenAiConfig `yaml:"aliyun"`
	GLM        OpenAiConfig `yaml:"glm"`
}

type OpenAiConfig struct {
	APIKey         string   `yaml:"api_key"`
	BaseURL        string   `yaml:"base_url"`
	DefaultModelID string   `yaml:"default_model_id"`
	Models         []string `yaml:"models"`
}
