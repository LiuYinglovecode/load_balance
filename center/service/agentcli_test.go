package service

import (
	"fmt"
	"testing"
	"time"

	"code.htres.cn/casicloud/alb/agent"
	"code.htres.cn/casicloud/alb/apis"
	"code.htres.cn/casicloud/alb/pkg/model"
)

var rpc = "0.0.0.0:6789"
var config = &agent.Config{
	RPC:    rpc,
}

func TestAgentClientSendCommand(t *testing.T) {
	rpcServer, err := agent.NewRPCServer(config, nil)
	if err != nil {
		t.Fatal(err)
	}

	quit := make(chan int)
	go func() {
		rpcServer.Run()
		quit <- 0
	}()

	centerConf := &AgentClient{
		AgentTimeout: time.Second * 60,
	}

	cmd := &apis.LBCommand{
		Cmd: apis.CmdDeploy,
	}
	r, err := centerConf.SendLBCommand(rpcServer.GetAddr(), cmd)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("get command result: code: %d, msg: %s", r.Code, r.Msg)
	rpcServer.Stop()
	<-quit
}

func TestAgentClientSendLBRequest(t *testing.T) {
	rpcServer, err := agent.NewRPCServer(config, nil)
	if err != nil {
		t.Fatal(err)
	}

	quit := make(chan int)
	go func() {
		rpcServer.Run()
		quit <- 0
	}()

	centerConf := &AgentClient{
		AgentTimeout: time.Second * 60,
	}

	var userID string = "1000"

	lbr := model.LBRecord{
		ID:    1,
		Type:  model.TypeIP,
		Owner: model.NewADCString(userID),
		IP:    model.NewADCString("106.75.69.100"),
		Port:  5901,
	}

	endPoints := []model.RealServer{
		{
			IP:   "10.10.10.1",
			Port: 3721,
			Name: "rs1",
		},
	}

	request := &model.LBRequest{
		Action:    model.ActionAdd,
		User:      model.NewADCString(userID),
		RequestID: 123456,
		Policy: model.LBPolicy{
			Record:       lbr,
			Endpoints: endPoints,
		},
	}
	r, err := centerConf.SendLBRequest(rpcServer.GetAddr(), request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("get command result: code: %d, msg: %s", r.Code, r.Msg)
	rpcServer.Stop()
	<-quit
}
