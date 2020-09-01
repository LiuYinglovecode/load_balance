package web

import (
	"code.htres.cn/casicloud/alb/center/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

func TestLBRecordController(t *testing.T) {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	config := &common.Config{
		WorkDir:  ".",
		LogLevel: 0,

		AuditLogPath: ".adc",
		SysLogPath:   ".alb",

		AgentTimeout: 5 * time.Second,

		Port: 8080,

		DBArgs:        "root:123456@tcp(172.17.60.108:3306)/adc?charset=utf8&parseTime=True",
		Dialect:       "mysql",
		EtcdEndpoints: []string{"localhost:1234"},
	}
	lbRecordController := NewLBRecordController(*config)
	lbPoolController := NewLBPoolController(*config)
	r := gin.Default()

	r.GET("/1/lbr/:lbr_id", lbRecordController.GetLBRecord)
	r.GET("/1/lbr",lbRecordController.ListLBRecord)
	r.POST("/1/lbr", lbRecordController.CreateLBRecord)
	r.PUT("/1/lbr/:lbr_id", lbRecordController.UpdateLBRecord)
	r.DELETE("/1/lbr/:lbr_id", lbRecordController.DropLBRecord)

	r.GET("/1/lbp/:lbp_id", lbPoolController.GetLBPool)
	r.GET("/1/lbp",lbPoolController.ListLBPool)
	r.POST("/1/lbp", lbPoolController.CreateLBPool)
	r.PUT("/1/lbp/:lbp_id", lbPoolController.UpdateLBPool)
	r.DELETE("/1/lbp/:lbp_id", lbPoolController.DropLBPool)
	return r
}



