package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"code.htres.cn/casicloud/alb/apis"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var logger = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: &logrus.TextFormatter{DisableColors: true},
	Level:     logrus.DebugLevel,
}

var rpc = "0.0.0.0:6789"
var config = &Config{
	RPC: rpc,
}

func Test_RPCServerInit(t *testing.T) {
	store, err := NewLBPolicyStore(config)
	assert.NoError(t, err)

	rpcServer, err := NewRPCServer(config, store)
	if err != nil {
		t.Fatal(err)
	}

	quit := make(chan int)
	go func() {
		rpcServer.Run()
		quit <- 0
	}()

	rpcServer.Stop()
	<-quit
}

func Test_RPCServerResponse(t *testing.T) {
	ifaces, err := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			fmt.Printf("get ip address %+v\n", ip)
		}
	}

	store, err := NewLBPolicyStore(config)
	assert.NoError(t, err)

	rpcServer, err := NewRPCServer(config, store)
	if err != nil {
		t.Fatal(err)
	}

	quit := make(chan int)
	go func() {
		rpcServer.Run()
		quit <- 0
	}()
	fmt.Printf("rpc server listen on: %s\n", rpcServer.GetAddr())
	conn, err := grpc.Dial(rpcServer.GetAddr(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %+v", err)
	}

	defer conn.Close()
	c := apis.NewCommanderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	cmd := apis.LBCommand{
		Cmd: apis.CmdDeploy,
	}
	r, err := c.Execute(ctx, &cmd)

	if err != nil {
		t.Fatalf("Execute: %+v", err)
	}

	fmt.Printf("command result is %+v", r)

	rpcServer.Stop()
	<-quit
}

func TestDeployHAProxy(t *testing.T) {
	ifaces, err := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			fmt.Printf("get ip address %+v\n", ip)
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	config.WorkDir = dir

	store, err := NewLBPolicyStore(config)
	assert.NoError(t, err)

	rpcServer, err := NewRPCServer(config, store)
	if err != nil {
		t.Fatal(err)
	}

	quit := make(chan int)
	go func() {
		rpcServer.Run()
		quit <- 0
	}()
	fmt.Printf("rpc server listen on: %s\n", rpcServer.GetAddr())
	conn, err := grpc.Dial(rpcServer.GetAddr(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %+v", err)
	}

	defer conn.Close()
	c := apis.NewCommanderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	policy := model.LBPolicy{
		Record: model.LBRecord{
			Name: "shanyou_lb",
			IP:   model.NewADCString("192.168.100.200"),
			Port: 80,
			Type: model.TypeIP},
		Endpoints: []model.RealServer{
			{Name: "sever1", IP: "106.75.69.8", Port: 80},
			{Name: "sever2", IP: "106.75.69.8", Port: 80},
			{Name: "sever3", IP: "106.75.69.8", Port: 80},
		},
	}

	request := model.LBRequest{
		RequestID: 1,
		User:      model.NewADCString("2"),
		Action:    model.ActionAdd,
		Policy:    policy,
	}

	raw, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(string(raw))
	cmd := apis.LBCommand{
		Cmd:  apis.CmdDeploy,
		Data: raw,
	}
	r, err := c.Execute(ctx, &cmd)

	if err != nil {
		t.Fatalf("Execute: %+v", err)
	}

	fmt.Printf("command result is %+v", r)

	rpcServer.Stop()
	<-quit
}
