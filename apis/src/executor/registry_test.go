package executor

import (
	"reflect"
	"testing"

	"code.htres.cn/casicloud/alb/apis"
	"github.com/stretchr/testify/assert"
)

type MockCommandExcutor struct {
}

func (e *MockCommandExcutor) Execute(cmd *apis.LBCommand) *apis.LBCommandResult {
	return apis.NotImplementRet
}
func (e *MockCommandExcutor) GetCMDType() apis.CmdType {
	return "mock"
}
func TestRegistry(t *testing.T) {
	reg := Registry{}
	origin := new(MockCommandExcutor)
	reg.RegisterExecutor(origin)
	assert.NotNil(t, reg.GetExecutor("mock"))

	executor := reg.MustGetExecutor("mock")

	assert.Equal(t, true, reflect.TypeOf(executor) == reflect.TypeOf(origin))
}

func TestRegistryPanic(t *testing.T) {
	reg := Registry{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	reg.MustGetExecutor("mock")
}
