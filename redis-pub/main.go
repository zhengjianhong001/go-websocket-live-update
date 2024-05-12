package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "",               // Redis密码，如果没有设置则为空
		DB:       0,                // 使用默认DB
	})

	// 每隔一段时间发布消息到频道
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		err := rdb.Publish(ctx, "mychannel", "Hello, World!").Err()
		if err != nil {
			panic(err)
		}
		fmt.Println("Message published to channel 'mychannel'")
	}
}
