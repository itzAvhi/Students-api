package config

type HTTPS struct {
	Address string
}

type Config struct {
	Env          string `yaml:"env" env:"Env" env-required:"True"`
	Storage_Path string `yaml:"Storage_Path"`
	HTTPS        `yaml:"http_server:"`
}
