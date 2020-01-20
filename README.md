# Demo

## 项目简介
利用rpc协议跨机房同步数据，从一台向多台同步数据(相当于1主多从)

## 使用Demo
1. 启动node1(leader节点)
```
./cmd/cmd -conf ./node1/configs
```
2. 启动node2
```
./cmd/cmd -conf ./node2/configs -join=1
```
3. 启动node3
```
./cmd/cmd -conf ./node3/configs -join=1
```
4. 启动node4
```
./cmd/cmd -conf ./node4/configs -join=1
```

5. 发送请求数据给leader，其它节点会同步到数据
```
go test -v test/adddata_test.go
```


## TODO
- [x] 实现raft协议，多台之间数据一致
- [ ] 跨机房的安全性问题



## 依赖
- [Kratos](https://github.com/bilibili/kratos)
