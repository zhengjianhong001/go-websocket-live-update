


登录校验



利用redis的pub/sub机制，用户的更新事件会被所有节点消费。从而无需cookie之类的机制确保用户的消息落到对应的节点上，便于扩缩容。

## 效果
页面会实时显示 redis mq 推送的内容

![alt text](image.png)

