package test

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func init() {
	host := "127.0.0.1:6379"
	pass := "123456a"
	//db := 8 // adminapi-8
	//db := 9 // merchantapi-9
	db := 10 // payapi-10

	// redis连接
	redisObj := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass, // no password set
		DB:       db,   // use default DB
	})

	// 检测redis连接是否正常
	if err := redisObj.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis连接失败err:%v", err)
	}

	Redis = redisObj
}
