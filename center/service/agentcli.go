package service

import (
	"code.htres.cn/casicloud/alb/center/common"
	"context"
	"encoding/json"
	"time"

	"code.htres.cn/casicloud/alb/apis"
	"code.htres.cn/casicloud/alb/pkg/model"
	"google.golang.org/grpc"
)

// Communicator Center-Agent通信接口
type Communicator interface {
	SendLBRequest(agentRPC string, request *model.LBRequest) (*apis.LBCommandResult, error);
    SendLBCommand(agentRPC string, cmd *apis.LBCommand) (*apis.LBCommandResult, error);
}

// AgentClient 实现与Agent的grpc通信
type AgentClient struct {
	AgentTimeout time.Duration
}

// NewAgentClient 构造函数
func NewAgentClient(timeout time.Duration) *AgentClient {
	return &AgentClient{timeout}
}

// SendLBRequest send request to remote
func (r *AgentClient) SendLBRequest(agentRPC string, request *model.LBRequest) (*apis.LBCommandResult, error) {
	raw, err := json.Marshal(request)
	if err != nil {
		common.SysLogger.Errorf("serialize LBRequest error: %v\n", err)
		return nil, err
	}

	cmd := &apis.LBCommand{
		Cmd:  apis.CmdDeploy,
		Data: raw,
	}
	ret, err := r.SendLBCommand(agentRPC, cmd)
	if err != nil {
		common.SysLogger.Errorf("execute LBCommand error: %v\n", err)
		return nil, err
	}

	return ret, nil
}

// SendLBCommand send lb command
func (r *AgentClient) SendLBCommand(agentRPC string, cmd *apis.LBCommand) (*apis.LBCommandResult, error) {
	conn, err := grpc.Dial(agentRPC, grpc.WithInsecure())
	if err != nil {
		common.SysLogger.Errorf("did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()
	c := apis.NewCommanderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), r.AgentTimeout)
	defer cancel()

	ret, err := c.Execute(ctx, cmd)
	if err != nil {
		common.SysLogger.Errorf("execute LBCommand error: %v\n", err)
		return nil, err
	}

	return ret, err
}
