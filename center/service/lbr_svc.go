package service

import (
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/center/dao"
	"code.htres.cn/casicloud/alb/pkg/model"
	"fmt"
	"math/rand"
)

const maxRetry = 10

// LbRecordService 处理lbrecord 相关的业务逻辑
type LbRecordService interface {
	GetAutoAllocatedLbr(userID int64) (*model.LBRecord, error)
}

// lbRecordServiceImpl 实现 LbRecordService接口
type lbRecordServiceImpl struct {
	lbRecordDao dao.LBRecordDao
	lbPoolDao   dao.LBPoolDao
}

// NewLbRecordService 构造函数
func NewLbRecordService(config common.Config) LbRecordService {
	dbConfig := dao.Config{
		DSN: config.DBArgs,
	}
	lbRecordDao := dao.NewLBRecordDao(&dbConfig)
	lbPoolDao := dao.NewLBPoolDaoImpl(&dbConfig)

	return &lbRecordServiceImpl{
		lbRecordDao: lbRecordDao,
		lbPoolDao:   lbPoolDao,
	}
}

// GetAutoAllocatedLbr 获取自动分配的lbr
// 暂时使用/后续逻辑需要修改
// 逻辑: 根据ip池和userID随机分配，有冲突则重试
func (l *lbRecordServiceImpl) GetAutoAllocatedLbr(userID int64) (*model.LBRecord, error) {
	// 获取可用ip池
	var pools []model.LBPool
	conditions := map[string]interface{}{"deleted": 0}
	err := l.lbPoolDao.ListLBPool(conditions, &pools)
	if err != nil {
		return nil, err
	}

	// 重试10次
	retry := 0
	for {
		// 随机选择一个ip和port
		lbp := pools[rand.Intn(len(pools))]
		ip := lbp.IP

		left := int(lbp.StartPort)
		right := int(lbp.EndPort)

		port := rand.Intn(right - left) + int(lbp.StartPort)

		lbrConditions := map[string]interface{}{
			"deleted": 0,
			"ip":      ip,
			"port":    port,
		}

		var lbrs []model.LBRecord
		err = l.lbRecordDao.ListLBRecord(lbrConditions, &lbrs)
		if err != nil {
			return nil, err
		}

		// 加入数据库
		if len(lbrs) == 0 {
			lbr := &model.LBRecord{
				IP:     ip,
				Port:   int32(port),
				Name:   "CICD PIPELINE",
				Status: 0,
				Owner:  model.NewADCString(userID),
			}
			lbr, err := l.lbRecordDao.CreateLBRecord(lbr)
			if err != nil {
				return nil, err
			}
			// 返回结果
			return lbr, nil
		}
		retry++
		if retry == maxRetry {
			return nil, fmt.Errorf("无可用ip")
		}
	}
}
