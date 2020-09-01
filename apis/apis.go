package apis

// CmdType commandType
type CmdType = string

const (
	//CmdDeploy command deploy haproxy
	CmdDeploy CmdType = "deploy"
)

const (
	// RetOk command request run ok
	RetOk int32 = iota
	// RetFail command request run fail
	RetFail
	// RetTimeout command request run timeout
	RetTimeout
	// RetUnknownCommand unknown command
	RetUnknownCommand
	// RetNotImplement not implement yet
	RetNotImplement
)

// UnknownCommandRet command result of unknown command
var UnknownCommandRet = &LBCommandResult{
	Code: RetUnknownCommand,
	Msg:  "Unknown command",
}

// NotImplementRet command result of not implement yet
var NotImplementRet = &LBCommandResult{
	Code: RetNotImplement,
	Msg:  "not implement yet",
}

// OKRet command result of process success
var OKRet = &LBCommandResult{
	Code: RetOk,
	Msg:  "success",
}

// CommandExecutor handle command
type CommandExecutor interface {
	Execute(*LBCommand) *LBCommandResult
	GetCMDType() CmdType
}
