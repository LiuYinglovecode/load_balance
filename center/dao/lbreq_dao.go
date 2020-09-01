package dao

import (
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/jinzhu/gorm"

	// gorm使用的数据库驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// LBRequestRepository  是LBRequest的存储相关操作接口
type LBRequestRepository interface {
	// Create 存储
	Create(request *model.LBRequest) error
	// GetById 通过Id获取
	GetByID(id int64) (*model.LBRequest, error)
	// Update 更新, request为空的字段不会更新到数据库
	// 执行后 request会和数据库中字段同步
	Update(request *model.LBRequest) error
	// DeleteById 通过Id删除
	DeleteByID(id int64) error
	// DB 返回gorm.DB
	DB() *gorm.DB
}

// LBRequestRepositoryImpl 实现了LBRequestRepository接口
type LBRequestRepositoryImpl struct {
	db *gorm.DB
}

// NewLBRequestRepository 构造函数
func NewLBRequestRepository(db *gorm.DB) LBRequestRepository {
	return &LBRequestRepositoryImpl{db}
}

// Create 存储到数据库
func (l *LBRequestRepositoryImpl) Create(request *model.LBRequest) error {
	if l.db.NewRecord(request) {
		request.RequestID = 0
	}
	if err := l.db.Create(request).Error; err != nil {
		common.SysLogger.Errorf("Failed create lbrequest %+v, reason: %s", *request, err)
		return err
	}
	common.SysLogger.Infof("Created lbrequest %+v", *request)
	return nil
}

// GetByID 根据ID从数据库读取对象
func (l *LBRequestRepositoryImpl) GetByID(id int64) (*model.LBRequest, error) {
	request := &model.LBRequest{}
	if err := l.db.Where("lb_request_id = ?", id).First(request).Error; err != nil {
		common.SysLogger.Errorf("Failed get request %+v, reason: %s", *request, err)
		return nil, err
	}
	common.SysLogger.Infof("Get request %+v", *request)
	return request, nil
}

// Update 更新
func (l *LBRequestRepositoryImpl) Update(request *model.LBRequest) error {
	if err := l.db.Model(request).Update(request).Error; err != nil {
		common.SysLogger.Errorf("Failed updated request %+v, reason: %s", *request, err)
		return err
	}
	common.SysLogger.Infof("Updated request %+v", *request)
	return nil
}

// DeleteByID 根据ID删除 lbrequest
func (l *LBRequestRepositoryImpl) DeleteByID(id int64) error {
	var request model.LBRequest
	if err := l.db.Where("lb_request_id = ?", id).Delete(&request).Error; err != nil {
		common.SysLogger.Errorf("Failed delete request, id: %+v, reason: %s", request, err)
		return err
	}
	common.SysLogger.Infof("Deleted request %d", id)
	return nil
}

// DB 返回数据库连接对象
func (l *LBRequestRepositoryImpl) DB() *gorm.DB {
	return l.db
}
