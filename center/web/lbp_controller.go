package web

import (
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/center/dao"
	"code.htres.cn/casicloud/alb/pkg/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LBPoolController struct {
	lbPoolDao dao.LBPoolDao
}

func NewLBPoolController(config common.Config) LBPoolController {
	dbConfig := dao.Config{
		DSN:config.DBArgs,
	}
	lbPoolDao := dao.NewLBPoolDaoImpl(&dbConfig)
	return LBPoolController{
		lbPoolDao,
	}
}

func (controller LBPoolController) GetLBPool (context *gin.Context)  {
	lbrId, _ := strconv.ParseInt(context.Params.ByName("lbr_id"), 10, 64)
	pool := model.LBPool{}
	err := controller.lbPoolDao.GetLBPool(int64(lbrId), &pool)
	if nil != err {
		context.String(http.StatusInternalServerError, err.Error())
	}
	var result, _ = json.Marshal(&pool)
	context.String(http.StatusOK, string(result))
}

func (controller LBPoolController) ListLBPool (context *gin.Context) {
	var pools []model.LBPool
	params := context.Request.URL.Query()
	err := controller.lbPoolDao.ListLBPool(getConditions(params), &pools)
	if nil != err {
		context.String(http.StatusInternalServerError, err.Error())
	}
	var result, _ = json.Marshal(pools)
	context.String(http.StatusOK, string(result))
}

func (controller LBPoolController) CreateLBPool (context *gin.Context) {
	var pool model.LBPool
	if err := context.ShouldBindJSON(&pool); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pool1, err := controller.lbPoolDao.CreateLBPool(&pool)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	var result, _ = json.Marshal(pool1)
	context.String(http.StatusOK, string(result))
}

func (controller LBPoolController) UpdateLBPool (context *gin.Context)  {
	id, _ := strconv.ParseInt(context.Params.ByName("lbp_id"), 10, 64)
	var attributes map[string]interface{}
	if err := context.ShouldBindJSON(&attributes); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pool, err := controller.lbPoolDao.UpdateAttribute(id, attributes)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	var result, _ = json.Marshal(pool)
	context.String(http.StatusOK, string(result))
}

func (controller LBPoolController) DropLBPool (context *gin.Context)  {
	id, _ := strconv.ParseInt(context.Params.ByName("lbp_id"), 10, 64)
	err := controller.lbPoolDao.DropLBPool(int64(id))
	if nil != err {
		context.String(http.StatusInternalServerError, err.Error())
	}
	context.String(http.StatusOK, "delete successfully!")
}
