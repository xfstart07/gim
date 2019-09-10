# GIM

IM 系统

## TODO LIST

第一阶段：

- [x]  1. 通过 API 接口实现用户的注册，并将信息存入 redis
- [x]  2. 服务端启动 http，grpc 服务，client 连接上 grpc 服务，采用 stream 进行消息推送
- [x]  3. 用户可以通过命令行进行实现私聊，群聊，通过接口将消息发送给服务端，服务端通过 grpc stream 流将消息推送给客户端
- [x]  4. 与服务端连接断开后，用户的下线功能

## protobuf

生成

```bash
protoc --go_out=plugins=grpc:internal/rpc_service -I protocol message.proto
```

## 消息

### 注册

向服务器注册账户

```bash
curl -X POST --header 'Content-Type: application/json' -d '{"user_name": "leon"}' http://localhost:8081/registerAccount
```

返回信息

```json
{"code":"0","message":"OK","data":{"user_id":1567996897857327000,"user_name":"baby"}}
```

## 使用

服务端启动

```bash
make run_server
```

客户端启动

```bash 
make run_client
```

要启动多个客户端需要指定用户信息

```bash
go run cmd/client/client.go --user_id=1568012126668462000 --username=kevin
```