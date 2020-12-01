package config

type Config struct {
	Server    Server    `yaml:"server"`
	Discovery Discovery `yaml:"discovery"`
	DB        DB        `yaml:"db"`
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
	Consul Consul `yaml:"consul"`
}

type Consul struct {
	Address string `yaml:"address"`
	Health  struct {
		Timeout                        int `yaml:"timeout"`
		Interval                       int `yaml:"interval"`
		DeregisterCriticalServiceAfter int `yaml:"deregisterCriticalServiceAfter"`
	} `yaml:"health"`
}

type DB struct {
	Mysql struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mysql"`
}
