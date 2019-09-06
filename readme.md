# GIM

IM 系统

## protobuf

生成

```bash
protoc --go_out=plugins=grpc:internal/rpc_service -I protocol message.proto
```