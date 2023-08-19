package controller

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var CTX context.Context

func Redisinit() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
	CTX = context.Background()
}
