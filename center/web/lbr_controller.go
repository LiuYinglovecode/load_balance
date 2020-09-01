package web

import (
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/center/dao"
	"code.htres.cn/casicloud/alb/center/service"
	"code.htres.cn/casicloud/alb/pkg/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

// LBRecordController 处理lbr相关的web请求
type LBRecordController struct {
	lbRecordDao dao.LBRecordDao
	lbRecordSvc service.LbRecordService
}
// NewLBRecordController 构造函数
func NewLBRecordController(config common.Config) LBRecordController {
	dbConfig := dao.Config{
		DSN:config.DBArgs,
	}
	lbRecordDao := dao.NewLBRecordDao(&dbConfig)
	lbrsvc := service.NewLbRecordService(config)
	return LBRecordController{
lbRecordDao,
lbrsvc,
	}
}

// GetAutoAllocatedLbr 获取为用户自动分配的Ip
func (controller LBRecordController) GetAutoAllocatedLbr(context *gin.Context) {
	userID, err := strconv.ParseInt(context.Query("user_id"), 10 ,64)
	if err != nil {
		context.JSON(http.StatusOK, model.NewApiResult(model.ErrParam, "解析请求参数失败, 原因:" + err.Error(), ""))
		return
	}

	lbr, err := controller.lbRecordSvc.GetAutoAllocatedLbr(userID)
	if err != nil {
		context.JSON(http.StatusOK, model.NewApiResult(model.ErrInternal, err.Error(), ""))
		return
	}

	data := make(map[string]string)

	data["ip"] = lbr.IP.String()
	data["port"] = strconv.Itoa(int(lbr.Port))

	context.JSON(http.StatusOK, model.NewApiResult(model.Ok, "ok", data))
}

func (controller LBRecordController) GetLBRecord (context *gin.Context)  {
	lbrId, _ := strconv.ParseInt(context.Params.ByName("lbr_id"), 10, 64)
	record := model.LBRecord{}
	err := controller.lbRecordDao.GetLBRecord(int64(lbrId), &record)
	if nil != err {
		context.String(http.StatusInternalServerError, err.Error())
	}
	var result, _ = json.Marshal(&record)
	context.String(http.StatusOK, string(result))
}

func (controller LBRecordController) ListLBRecord (context *gin.Context)  {
	var records []model.LBRecord
	params := context.Request.URL.Query()
	err := controller.lbRecordDao.ListLBRecord(getConditions(params), &records)
	if nil != err {
		context.String(http.StatusInternalServerError, err.Error())
	}
	var result, _ = json.Marshal(records)
	context.String(http.StatusOK, string(result))
}

func getConditions (params url.Values) (conditions map[string]interface{}) {
	conditions = map[string]interface{}{}
	for key, value := range params{
		if 0 < len(value){
			conditions[key] = value[0]
		}
	}
	return conditions
}

func (controller LBRecordController) CreateLBRecord (context *gin.Context)  {
	var record model.LBRecord
	if err := context.ShouldBindJSON(&record); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record1, err := controller.lbRecordDao.CreateLBRecord(&record)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	var result, _ = json.Marshal(record1)
	context.String(http.StatusOK, string(result))
}

func (controller LBRecordController) UpdateLBRecord (context *gin.Context)  {
	id, _ := strconv.ParseInt(context.Params.ByName("lbr_id"), 10, 64)
	var attributes map[string]interface{}
	if err := context.ShouldBindJSON(&attributes); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record, err := controller.lbRecordDao.UpdateAttribute(id, attributes)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	var result, _ = json.Marshal(record)
	context.String(http.StatusOK, string(result))
}

func (controller LBRecordController) DropLBRecord (context *gin.Context)  {
	id, _ := strconv.ParseInt(context.Params.ByName("lbr_id"), 10, 64)
	err := controller.lbRecordDao.DropLBRecord(int64(id))
	if nil != err {
		context.String(http.StatusInternalServerError, err.Error())
	}
	context.String(http.StatusOK, "delete successfully!")
}
