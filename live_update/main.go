package main

import (
	"context"
	"fmt"
	"time"

	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/rs/cors"
)

func main() {
	pt := polling.Default
	wt := websocket.Default
	wt.CheckOrigin = func(req *http.Request) bool {
		return true
	}
	server := socketio.NewServer(&engineio.Options{
		//PingInterval: pingInterval,
		//PingTimeout:  pingTimeout,
		// Socket.IO默认支持两种传输方式：polling和websocket。确保客户端和服务器都支持至少一种传输方式。
		Transports: []transport.Transport{
			pt,
			wt,
		},
	})
	// server.Adapter(&socketio.RedisAdapterOptions{
	// 	Addr:    "/tmp/docker/redis.sock",
	// 	Network: "unix",
	// })
	// // 使用 redis 作为适配器
	// _, err := server.Adapter(&socketio.RedisAdapterOptions{
	// 	Addr:     "127.0.0.1",
	// 	Password: "",
	// })
	// if err != nil {
	// 	fmt.Println("socket adapter redis set error.", err)
	// }
	roomID := "liveUpdate"

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		s.Join(roomID) // 加入房间
		return nil
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		// server.Remove(s.ID())
		fmt.Println("报错：meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// Add the Remove session id. Fixed the connection & mem leak
		fmt.Println("关闭closed", reason)
	})

	// 监听 notice 事件
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
		server.BroadcastToRoom("", roomID, "bcard", msg+" 上线了") // 广播给所有用户
	})
	go server.Serve()
	defer server.Close()
	go subLiveUpdate(server, roomID)
	// 设置CORS中间件选项
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost"}, // 明确指定允许的源
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true, // 允许请求携带凭证
		// 可以根据需要添加其他CORS配置项
	})
	log.Println("Serving at :3032...")
	mux := http.NewServeMux()
	mux.Handle("/socket.io/", server)
	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":3032", handler))

}

func subLiveUpdate(serv *socketio.Server, roomID string) {
	var ctx = context.Background()
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "",               // Redis密码，如果没有设置则为空
		DB:       0,                // 使用默认DB
	})
	// 订阅一个频道
	pubsub := rdb.Subscribe(ctx, "mychannel")
	log.Println("redis Subscribe at mychannel...")
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
			// 广播给房间roomID里面的 change 事件
			serv.BroadcastToRoom("", roomID, "change", time.Now().String()+"：数据更新") // 广播给所有用户
		}
	}()
}
