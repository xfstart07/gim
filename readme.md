# GIM

Golang 实现 的 IM 系统。

## TODO LIST

第一阶段：

- [x]  1. 通过api接口实现用户的注册，并将信息存入 redis
- [x]  2. 服务端启动 http，grpc 服务，client 能连接上
- [x]  3. 用户可以通过命令行进行实现私聊，群聊
- [x]  4. 与服务端连接断开后，用户的下线功能
- [x]  5. 格式化客户端信息打印
- [x]  6. 聊天信息存储文件中。
- [x]  7. 接入 etcd 做客户端和服务端的服务发现。
- [x]  8. 完善客户端重连和服务端连接下线功能
- [x]  9. 通过 redis 的 pubsub 分发消息
- [x]  10. 重构代码，将对象接口化，合理化
- [x]  11. 用户退出，注销
- [x]  12. 客户端获取所有用户列表
- [ ]  13. 添加测试用例

## protobuf

生成 GRPC 接口描述文件

```bash
make gen_proto
```

## 消息

### 注册

向服务器注册账户

```bash
curl -X POST --header 'Content-Type: application/json' -d '{"user_name": "leon"}' http://localhost:8081/registerAccount
```

返回信息

```json
{"code":"0","message":"OK","data":{"user_id":1567996897857327000,"user_name":"leon"}}
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

要启动多个客户端需要修改 `config/client.ini` 的用户信息

## 构建

```bash
make build_server
```

## 部署

### 配置 server

请查看 [deployment/nginx/gim.conf](deployment/nginx/gim.conf) 文件.