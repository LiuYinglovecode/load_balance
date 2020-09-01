package agent

import (
	"encoding/json"

	"code.htres.cn/casicloud/alb/apis"
	"code.htres.cn/casicloud/alb/pkg/model"
)

// DeployExecutor implement deply command
type DeployExecutor struct {
	controller LBController
}

// Execute impelement interface CommandExecutor
func (d *DeployExecutor) Execute(cmd *apis.LBCommand) *apis.LBCommandResult {
	var request model.LBRequest
	err := json.Unmarshal(cmd.Data, &request)
	if err != nil {
		return &apis.LBCommandResult{
			Code: apis.RetFail,
			Msg:  err.Error(),
		}
	}

	switch request.Action {
	case model.ActionAdd:
		err := d.controller.StartLB(request.Policy)
		if err != nil {
			return &apis.LBCommandResult{
				Code: apis.RetFail,
				Msg:  err.Error(),
			}
		}
	case model.ActionUpdate:
		err := d.controller.UpdateLB(request.Policy)
		if err != nil {
			return &apis.LBCommandResult{
				Code: apis.RetFail,
				Msg:  err.Error(),
			}
		}
	case model.ActionStop:
		err := d.controller.StopLB(request.Policy)
		if err != nil {
			return &apis.LBCommandResult{
				Code: apis.RetFail,
				Msg:  err.Error(),
			}
		}
	}

	return apis.OKRet
}

// GetCMDType impelement interface CommandExecutor
func (d *DeployExecutor) GetCMDType() apis.CmdType {
	return apis.CmdDeploy
}
