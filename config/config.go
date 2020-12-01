package config

import "flag"

var ConfPath string

func init() {
	flag.StringVar(&ConfPath, "c", "./config.yml", "配置文件的地址")
}

type Server struct {
	Ip  string `yaml:"ip"`
	Tcp struct {
		Port int `yaml:"port"`
	} `yaml:"tcp"`
	Http struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
}

type Discovery struct {
	Consul struct {
		Address string `yaml:"address"`
		Health  struct {
			Timeout                        int `yaml:"timeout"`
			Interval                       int `yaml:"interval"`
			DeregisterCriticalServiceAfter int `yaml:"deregisterCriticalServiceAfter"`
		} `yaml:"health"`
	} `yaml:"consul"`
}

type DB struct {
	Mysql struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Schema   string `yaml:"schema"`
	} `yaml:"mysql"`
}
