package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type LoggerConf struct {
	FileDirectory   string   `yaml:"file_directory"`
	MaxSize         int      `yaml:"maxsize"`
	MaxBackups      int      `yaml:"max_backups"`
	MaxAge          int      `yaml:"max_age"`
	Compress        bool     `yaml:"compress"`
	Level           int32    `yaml:"level"`
	LogIDShowHeader bool     `yaml:"log_id_show_header"`
	SkipPaths       []string `yaml:"skip_paths"`
}

type Server struct {
	Port int `yaml:"port" `
}

type Resource struct {
	Mysql  Mysql  `yaml:"mysql"`
	Redis  Redis  `yaml:"redis"`
	TTSRPC TTSRPC `yaml:"tts_rpc"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Charset  string `yaml:"charset"`
	LogLevel int    `yaml:"log_level"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type TTSRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type RDS struct {
	DbName      string `yaml:"Database"`
	Psm         string `yaml:"PSM"`
	LogLevel    int    `yaml:"LogLevel"`
	ShardingNum int    `yaml:"ShardingNum"`
}

type Conf struct {
	AppName           string     `yaml:"app_name"`
	RunMode           string     `yaml:"run_mode"`
	Server            Server     `yaml:"server" env:"GO_ENV"`
	Resource          Resource   `yaml:"resource"`
	Log               LoggerConf `yaml:"log"`
	AccessTokenSecret string     `yaml:"access_token_secret"`
}

var Instance = &Conf{}

func (m *Conf) LoadConfig(configFile string) error {

	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil

}
