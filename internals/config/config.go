package config

type HTTPS struct {
	Address string
}

type Config struct {
	Env          string
	Storage_Path string
	HTTPS
}
