package dao

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/jinzhu/gorm"
)

type LBPoolDao interface {
	GetLBPool(id int64, Pool *model.LBPool) error
	ListLBPool(conditions map[string]interface{}, Pools *[]model.LBPool) error
	CreateLBPool(Pool *model.LBPool) (*model.LBPool, error)
	UpdateLBPool(Pool *model.LBPool) (*model.LBPool, error)
	UpdateAttribute(id int64, attributes map[string]interface{}) (*model.LBPool, error)
	DropLBPool(id int64) error
}

type LBPoolDaoImpl struct {
	config *Config
	db     *gorm.DB
}

// NewLBRecordDao BaseDao constructor
func NewLBPoolDaoImpl(c *Config) (d *LBPoolDaoImpl) {
	d = &LBPoolDaoImpl{
		config: c,
		db:     NewMysql(c),
	}
	d.db.SingularTable(true)
	return
}

// Close release all mysql resource .
func (dao *LBPoolDaoImpl) Close() {
	dao.db.Close()
}

func (dao *LBPoolDaoImpl) GetLBPool(id int64, pool *model.LBPool) (err error) {
	po := model.LBPool{ID: id, Deleted: false}

	err = dao.db.Where(&po).First(pool).Error
	if err != nil {
		dao.config.Logger.Error(err)
	}

	var agents []model.LBAgent
	dao.db.Model(&pool).Association("Agents").Find(&agents)

	pool.Agents = agents
	return
}
func (dao *LBPoolDaoImpl) ListLBPool(conditions map[string]interface{}, pools *[]model.LBPool) (err error) {
	conditions["deleted"] = false;
	err = dao.db.Find(pools, conditions).Error
	if err != nil {
		dao.config.Logger.Error(err)
	}

	for i, p := range *pools {
		var agents []model.LBAgent
		dao.db.Model(&p).Association("Agents").Find(&agents)
		(*pools)[i].Agents = agents
	}
	return
}
func (dao *LBPoolDaoImpl) CreateLBPool(pool *model.LBPool) (*model.LBPool, error) {
	err := dao.db.Create(pool).Error
	if err != nil {
		dao.config.Logger.Error(err)
	}
	return pool, err
}
func (dao *LBPoolDaoImpl) UpdateLBPool(pool *model.LBPool) (*model.LBPool, error) {
	err := dao.db.Save(pool).Error
	if err != nil {
		dao.config.Logger.Error(err)
	}
	return pool, err
}
func (dao *LBPoolDaoImpl) UpdateAttribute(id int64, attributes map[string]interface{}) (*model.LBPool, error) {
	pool := &model.LBPool{}
	err := dao.db.Where(&model.LBPool{ID: id}).First(pool).Error
	if err != nil {
		dao.config.Logger.Error(err)
		return pool, err
	}
	err = dao.db.Model(pool).Where(&model.LBPool{ID: id}).Update(attributes).Error
	if err != nil {
		dao.config.Logger.Error(err)
		return pool, err
	}
	return pool, nil
}
func (dao *LBPoolDaoImpl) DropLBPool(id int64) (err error) {
	return dao.db.Model(&model.LBPool{}).Where("lb_pool_id = ?", id).Update("deleted", true).Error
}
