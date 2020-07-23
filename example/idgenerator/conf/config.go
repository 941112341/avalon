package conf

type config struct {
	Database struct {
		DB     string `yaml:"DB"`
		DBRead string `yaml:"DBRead"`
	} `yaml:"Database"`
}

var Config config
