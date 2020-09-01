apis
===
contanis internal apis

# 创建grpc go代码的方法
```bash
protoc -I apis/ apis/agent_apis.proto --go_out=plugins=grpc:apis
```

# 使用apis说明
api处理过程通过grpc对外提供服务
创建服务器端程序如下:
```golang
rpcServer, err := agent.NewRPCServer(config)
	if err != nil {
		t.Fatal(err)
	}

	quit := make(chan int)
	go func() {
		rpcServer.Run()
		quit <- 0
	}()
	// processing...
	rpcServer.Stop()
	<-quit
```

客户端调用如下:
```golang
    centerConf := &center.Config{
		AgentTimeout: time.Second * 60,
		Logger:       logger,
	}

	cmd := &apis.LBCommand{
		Cmd: apis.CmdDeploy,
	}
	r, err := centerConf.SendLBCommand(rpcServer.GetAddr(), cmd)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("get command result: code: %d, msg: %s", r.Code, r.Msg)
```