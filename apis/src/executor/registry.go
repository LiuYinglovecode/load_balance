package executor

import (
	"fmt"

	"code.htres.cn/casicloud/alb/apis"
)

// Registry registry for executor
type Registry struct {
	cmdExecutors []apis.CommandExecutor
}

//RegisterExecutor 注册executor
func (r *Registry) RegisterExecutor(exec apis.CommandExecutor) {
	if r.GetExecutor(exec.GetCMDType()) != nil {
		panic(fmt.Sprintf("duplicate executor: %s", exec.GetCMDType()))
	}

	r.cmdExecutors = append(r.cmdExecutors, exec)
}

//GetExecutor 根据cmdtype获取handler
func (r *Registry) GetExecutor(cmdType apis.CmdType) apis.CommandExecutor {
	for _, exec := range r.cmdExecutors {
		if exec.GetCMDType() == cmdType {
			return exec
		}
	}

	return nil
}

// MustGetExecutor 指定命令执行器不存在时，报错
func (r *Registry) MustGetExecutor(cmdType apis.CmdType) apis.CommandExecutor {
	exec := r.GetExecutor(cmdType)

	if exec == nil {
		panic(fmt.Sprintf("executor not register! %s", cmdType))
	}

	return exec
}
