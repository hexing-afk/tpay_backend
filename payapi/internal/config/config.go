package config

import "github.com/tal-tech/go-zero/rest"

type Config struct {
	rest.RestConf
	Timezone string
	Mysql    struct {
		DataSource string
	}
	Redis struct {
		Host string
		Pass string
		DB   int
	}
	Static struct {
		Path string
	}
}
