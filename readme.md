# GIM

IM 系统

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