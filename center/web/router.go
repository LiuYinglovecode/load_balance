package web

import (
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/center/dao"
	"code.htres.cn/casicloud/alb/center/queue"
	"code.htres.cn/casicloud/alb/center/service"
	"github.com/gin-gonic/gin"
)

// SetupRouter 建立路由关系
func SetupRouter() (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Recovery())
	// todo: 配置gin的日志格式
	router.Use(gin.Logger())

	lc, err := setupLBRequestController()
	if err != nil {
		common.SysLogger.Errorf("LBRequest handler service create failed, reason: %v", err)
		return nil, err
	}

	// lbrequest
	router.POST("/1/lb", lc.AcceptRequest)
	router.GET("/1/lb/:request_id", lc.QueryRequest)

	// lbr
	lbRecordController := NewLBRecordController(common.GlobalConfig)
	router.GET("/1/lbr/:lbr_id", lbRecordController.GetLBRecord)
	router.GET("/1/lbr", lbRecordController.ListLBRecord)
	router.POST("/1/lbr", lbRecordController.CreateLBRecord)
	router.PUT("/1/lbr/:lbr_id", lbRecordController.UpdateLBRecord)
	router.DELETE("/1/lbr/:lbr_id", lbRecordController.DropLBRecord)

	router.POST("/1/lbr/acquire", lbRecordController.GetAutoAllocatedLbr)

	// lbp
	lbPoolController := NewLBPoolController(common.GlobalConfig)
	router.GET("/1/lbp/:lbp_id", lbPoolController.GetLBPool)
	router.GET("/1/lbp", lbPoolController.ListLBPool)
	router.POST("/1/lbp", lbPoolController.CreateLBPool)
	router.PUT("/1/lbp/:lbp_id", lbPoolController.UpdateLBPool)
	router.DELETE("/1/lbp/:lbp_id", lbPoolController.DropLBPool)

	return router, nil
}

// setupLBRequestSvc 创建LBRequest相关的路由处理函数
func setupLBRequestController() (*LBRequestController, error) {
	lbreqRepo := dao.NewLBRequestRepository(common.GetDataSource())
	lbaRepo := dao.NewLBAgentRepository(common.GetDataSource())

	dbConfig := dao.Config{
		DSN: common.GlobalConfig.DBArgs,
	}
	lbrRepo := dao.NewLBRecordDao(&dbConfig)
	lbpRepo := dao.NewLBPoolDaoImpl(&dbConfig)

	qm, err := queue.NewEtcdMessageQueueHandler(common.GlobalConfig.EtcdEndpoints,
		common.GlobalConfig.EtcdCAPath,
		common.GlobalConfig.EtcdCertPath,
		common.GlobalConfig.EtcdKeyPath,
		lbaRepo)
	if err != nil {
		common.SysLogger.Errorf("LBRequest message quque create failed, reason %v", err)
		return nil, err
	}

	agentcli := service.NewAgentClient(common.GlobalConfig.AgentTimeout)

	svc := service.NewRequestHandlerImpl(lbrRepo, lbpRepo, lbreqRepo, qm, agentcli)

	ch := make(chan interface{})
	pr := func() {
		// recoevery from groutine pacnic
		defer func() {
			if r := recover(); r != nil {
				common.SysLogger.Errorf("ProcessRequest goroutine panic, recovery from %v", r)
				ch <- r
			}
		}()
		svc.ProcessRequest()
	}

	// 启动线程, 处理队列
	// todo: 是否应该启动多个
	go pr()

	// 监听处理线程是否终止，如果异常终止重启
	go func() {
		for {
			e := <-ch
			common.SysLogger.Errorf("ProcessRequest goroutine aborted, re-run..., reason: %v,", e)
			go pr()
		}
	}()

	return NewLBRequestController(svc), nil
}
