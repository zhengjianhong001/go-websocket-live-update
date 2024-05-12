package main

import (
	"context"
	"fmt"

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

	// 订阅一个频道
	pubsub := rdb.Subscribe(ctx, "mychannel")

	// 在Go协程中等待消息
	go func() {
		for {
			// 获取消息
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}

			// 处理消息
			fmt.Printf("Received message from %s: %s\n", msg.Channel, msg.Payload)
			// 这里可以添加你的消息处理逻辑
		}
	}()

	// 阻塞主线程，防止程序退出
	select {}
}
