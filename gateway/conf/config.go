package conf

type config struct {
	Https struct {
		Port int `yaml:"port"`
	} `yaml:"https"`
	Http struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
	Database struct {
		DB     string `yaml:"DB"`
		DBRead string `yaml:"DBRead"`
	} `yaml:"Database"`
}

var Config config
