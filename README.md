# Demo

## 项目简介
利用rpc协议跨机房同步数据，从一台向多台同步数据(相当于1主多从)

## 使用Demo
1. 启动server1
```
cd cmd
go build
./cmd -conf ../configs
```
2. 启动server2
```
cd cmd
go build
./cmd -conf ../configs2
```
3. 发送请求数据给server2，同时server2向server1同步数据
```
go test -v test/sync_test.go
```


## TODO
- [ ] 实现raft协议，多台之间数据一致
- [ ] 跨机房的安全性问题



## 依赖
- [Kratos](https://github.com/bilibili/kratos)
