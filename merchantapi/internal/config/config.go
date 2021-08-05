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
	CloudStorage string // 云存储工具
	OssStorage   struct {
		AccessKeyId     string
		SecretAccessKey string
		Endpoint        string
		Bucket          string
	}
	S3Storage struct {
		AccessKeyId     string
		SecretAccessKey string
		Region          string
		Bucket          string
	}
}
