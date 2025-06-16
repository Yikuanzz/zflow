# 服务注册中心

```proto
// 服务注册，返回 LeaseID 供后续心跳续期
rpc Register(ServiceInstance) returns (Lease) {}
// 心跳续期；若返回 NOT_FOUND 需重新 Register
rpc KeepAlive(Lease) returns (Lease) {}
// 主动注销
rpc Deregister(Lease) returns (google.protobuf.Empty) {}
// 一次性拉取
rpc Discover(Query) returns (Services) {}
// 长连接订阅；服务器发现变化即推送
rpc Watch(Query) returns (stream Services) {}
```
