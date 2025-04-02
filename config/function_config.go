package config

// FunctionConfig 函数服务配置
type FunctionConfig struct {
	Weather WeatherConfig `yaml:"weather"`
}

type WeatherConfig struct {
	Qweather QweatherConfig `yaml:"qweather"`
}

type QweatherConfig struct {
	ApiKey  string `yaml:"api_key"`
	ApiHost string `yaml:"api_host"`
}
