
[toc]
# 创建项目

## 安装依赖

``` shell
# < go 1.16
go get -u github.com/TarsCloud/TarsGo/tars/tools/tarsgo
go get -u github.com/TarsCloud/TarsGo/tars/tools/tars2go
# >= go 1.16
go install github.com/TarsCloud/TarsGo/tars/tools/tarsgo@latest
go install github.com/TarsCloud/TarsGo/tars/tools/tars2go@latest
```


## 创建服务 
https://doc.tarsyun.com/#/dev/tarsgo/tarsgo.md

```shell

tarsgo make App Server Servant GoModuleName
# 例如：
tarsgo make TestApp HelloGo SayHello github.com/Tars/test

tarsgo make Demo UserServer UserServerObj demo-user

```
Demo: 项目名称
UserServer: 服务名
UserServerObj: tars后台的Servant名
demo-user: 服务的go mod名称



## 生成golang协议
```
tars2go  -outdir=tars-protocol -module=demo-user UserServerObj.tars
```


## 上传部署

请注意修改 makefile.tars.gomod.mk 

部署失败常见问题：

1、部署后启动时找不到执行问题
2、... 