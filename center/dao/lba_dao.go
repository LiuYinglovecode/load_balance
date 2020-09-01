package dao

import (
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/jinzhu/gorm"
)

// LBAgentRepository  是LBAgent的存储相关操作接口
type LBAgentRepository interface {
	// Create 存储
	Create(agent *model.LBAgent) error
	// GetById 通过Id获取
	GetByID(id int64) (*model.LBAgent, error)
	// Update 更新, 空的字段不会更新到数据库
	// 执行后 参数agent会和数据库中字段同步
	Update(agent *model.LBAgent) error
	// DeleteById 通过Id删除
	DeleteByID(id int64) error
	// DB 返回gorm.DB
	DB() *gorm.DB
	// List
	List(conditions map[string]interface{}) ([]model.LBAgent, error)
}

// LBAgentRepositoryImpl 实现了LBAgentRepository接口
type LBAgentRepositoryImpl struct {
	db *gorm.DB
}

// NewLBAgentRepository 构造函数
func NewLBAgentRepository(db *gorm.DB) LBAgentRepository {
	return &LBAgentRepositoryImpl{db}
}

// Create 存储到数据库
func (l *LBAgentRepositoryImpl) Create(agent *model.LBAgent) error {
	if l.db.NewRecord(agent) {
		agent.ID = 0
	}
	if err := l.db.Create(agent).Error; err != nil {
		common.SysLogger.Errorf("Failed create lbagent %+v, reason: %s", *agent, err)
		return err
	}
	common.SysLogger.Infof("Created lbagent %+v", *agent)
	return nil
}

// GetByID 根据ID从数据库读取对象
func (l *LBAgentRepositoryImpl) GetByID(id int64) (*model.LBAgent, error) {
	agent := &model.LBAgent{}
	if err := l.db.Where("lb_agent_id = ?", id).First(agent).Error; err != nil {
		common.SysLogger.Errorf("Failed get agent %+v, reason: %s", *agent, err)
		return nil, err
	}
	common.SysLogger.Infof("Get agent %+v", *agent)
	return agent, nil
}

// Update 更新
func (l *LBAgentRepositoryImpl) Update(agent *model.LBAgent) error {
	if err := l.db.Model(agent).Update(agent).Error; err != nil {
		common.SysLogger.Errorf("Failed updated agent %+v, reason: %s", *agent, err)
		return err
	}
	common.SysLogger.Infof("Updated agent %+v", *agent)
	return nil
}

// DeleteByID 根据ID删除 lbagent
func (l *LBAgentRepositoryImpl) DeleteByID(id int64) error {
	var agent model.LBAgent
	if err := l.db.Where("lb_agent_id = ?", id).Delete(&agent).Error; err != nil {
		common.SysLogger.Errorf("Failed delete agent, id: %+v, reason: %s", agent, err)
		return err
	}
	common.SysLogger.Infof("Deleted agent %d", id)
	return nil
}

// DB 返回数据库连接对象
func (l *LBAgentRepositoryImpl) DB() *gorm.DB {
	return l.db
}

// List 返回Agent列表
func (l *LBAgentRepositoryImpl) List(conditions map[string]interface{}) ([]model.LBAgent, error) {
	conditions["deleted"] = false
	agents := []model.LBAgent{}
	if err := l.db.Find(&agents, conditions).Error; err != nil {
		common.SysLogger.Errorf("Failed list agent, conditions: %+v, reason: %s", conditions, err)
		return nil, err
	}
	common.SysLogger.Infof("List lbagent on condition: %+v", conditions)
	return agents, nil
}
