package agent

import (
	"context"
	"net"

	"code.htres.cn/casicloud/alb/apis"
	"code.htres.cn/casicloud/alb/apis/executor"
	"google.golang.org/grpc"
)

// RPCServer for control agent haproxy deploy
type RPCServer struct {
	lis          net.Listener
	serv         *grpc.Server
	execRegistry executor.Registry
	controller   *Controller
}

//NewRPCServer create new rpc server by config
func NewRPCServer(config *Config, store LBPolicyStore) (*RPCServer, error) {
	lis, err := net.Listen("tcp", config.RPC)
	if err != nil {
		sysLogger.Errorf("RpcServer failed to listen %s, reason: %v\n", config.RPC, err)
		return nil, err
	}

	cont, err := NewController(config, store)
	if err != nil {
		sysLogger.Errorf("Create LBController failed reason: %s\n", err)
		return nil, err
	}

	//注册命令执行器
	registry := executor.Registry{}
	deployExec := &DeployExecutor{
		controller: cont,
	}
	registry.RegisterExecutor(deployExec)

	serv := grpc.NewServer()
	server := &RPCServer{
		serv:         serv,
		lis:          lis,
		execRegistry: registry,
		controller:   cont,
	}

	apis.RegisterCommanderServer(serv, server)
	return server, nil
}

// Run controller
func (s *RPCServer) Run() error {
	return s.serv.Serve(s.lis)
}

// GetAddr get server address for dail
func (s *RPCServer) GetAddr() string {
	return s.lis.Addr().String()
}

// Stop stop server
func (s *RPCServer) Stop() {
	s.serv.Stop()
}

// Execute implements apis CommanderServer interfaces
func (s *RPCServer) Execute(ctx context.Context, cmd *apis.LBCommand) (*apis.LBCommandResult, error) {
	//TODO: 请实现具体的命令处理逻辑

	exec := s.execRegistry.GetExecutor(cmd.Cmd)
	// try to deploy haproxy
	if exec != nil {
		return exec.Execute(cmd), nil
	}

	return apis.UnknownCommandRet, nil
}
