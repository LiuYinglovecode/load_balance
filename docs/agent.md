Agent 设计
===
agent负载管理代理服务的，如HAProxy， Nginx等，部署在虚拟主机上。主要功能如下：

1. 负责负载均衡服务的启动停止更改配置

2. 提供API接口让Center来控制

3. 向etcd注册自己，汇报状态

4. 提供metrics向metrics收集器汇报各种监控指标

5. agent要能以审计日志的形式记录每一次center操作

6. agent运行全过程要提供日志

7. agent需要实现重启后可以重新恢复运行前的状态

设计上要求agent对除宿主机外其他硬件和服务达到最小的依赖。

# 模型
## `RPCServer` 
实现grpc服务接口，处理控制命令请求

## `Controller`
处理部署请求，管理`LBPolicy`。维护LBPolicy与容器一对一的关系。

## `Informer`
负责向`etcd`汇报自己的状态

## `Logger`
记录运行日志，审计日志。是个全局单例，在agent初始化时根据配置文件创建。通过全局变量在各个模块中调用。

## `lbagent`
agent命令工具入口，主要实现启动RPCServer, 在本地管理查看状态等。
# 详细设计
## grpc 实现远程Center控制
在agent上开起grpc服务。提供Command类型的调用。center通过`LBCommand`发送到agent，调用`Commander.Execute`。agent负责实现这个接口。处理不同`LBCommand请求`。处理完成后通过`LBCommandResult`返回处理结果。

为了支持部署控制docker container。实现`deply_executor`来处理。


## 故障恢复机制
如果agent出现崩溃或者宿主机宕机，在agent重启后可以恢复。
## 日志
