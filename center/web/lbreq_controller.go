package web

import (
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/center/service"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// LBRequestController web layer cnotroller
type LBRequestController struct {
	service service.LBRequestHandler
}

// NewLBRequestController 构造函数
func NewLBRequestController(handler service.LBRequestHandler) *LBRequestController {
	return &LBRequestController{handler}
}

// AcceptRequest 接受client的lbrequest, 进行处理
// todo: 不应该直接将error信息返回给client
func (l *LBRequestController) AcceptRequest(c *gin.Context) {
	var req model.LBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.SysLogger.Errorf("Parse lbrequest body failed, reason: %v", err)
		c.JSON(http.StatusOK, model.NewApiResult(model.ErrParam, "解析请求失败", ""))
		return
	}

	reqid, err := l.service.AcceptRequest(&req)
	if err != nil {
		common.SysLogger.Errorf("Handle LBRequest failed, reason: %v", err)
		c.JSON(http.StatusOK, model.NewApiResult(model.ErrUnknown, err.Error(), ""))
	} else {
		data := make(map[string]string)

		data["request_id"] = strconv.FormatInt(reqid, 10)

		c.JSON(http.StatusOK, model.NewApiResult(model.Ok, "ok", data))
	}
}

// QueryRequest 客户端查询lbrequest处理进度
func (l *LBRequestController) QueryRequest(c *gin.Context) {
	strReqid := c.Param("request_id")
	reqId, err := strconv.ParseInt(strReqid, 10, 64)
	if err != nil {
		common.SysLogger.Errorf("Parse request_id failed, reason: %v", err)
		c.JSON(http.StatusOK, model.NewApiResult(model.ErrParam, "解析请求失败, 错误的id", ""))
		return
	}
	req, err := l.service.QueryRequest(reqId)
	if err != nil {
		common.SysLogger.Errorf("Handle LBRequest failed, reason: %v", err)
		c.JSON(http.StatusOK, model.NewApiResult(model.ErrUnknown, err.Error(), ""))
		return
	}

	data := make(map[string]int32)

	data["status"] = req

	c.JSON(http.StatusOK, model.NewApiResult(model.Ok, "ok", data))
}

// DeleteRequest 客户端终止对指定lbrequest的处理
// 如果已经处理完毕，删除请求所作的操作
func (l *LBRequestController) DeleteRequest(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"msg": "not implemented"})
}

// UpdateRequest 更新已经发送的请求
// 暂未实现
func (l *LBRequestController) UpdateRequest(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"msg": "not implemented"})
}
