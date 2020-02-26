# GIM

Golang 实现 的 IM 系统。

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

### Ref 

You've got Mail.

<https://v.qq.com/x/cover/xnbwwtilv114645.html>
