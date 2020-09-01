package dao

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type LBRecordDao interface {
	GetLBRecord(id int64, record *model.LBRecord) error
	ListLBRecord(conditions map[string]interface{}, records *[]model.LBRecord) error
	CreateLBRecord(record *model.LBRecord) (*model.LBRecord, error)
	UpdateLBRecord(record *model.LBRecord) (*model.LBRecord, error)
	UpdateAttribute(id int64, attributes map[string]interface{}) (*model.LBRecord, error)
	DropLBRecord(id int64) error
}

type LBRecordDaoImpl struct {
	config *Config
	db     *gorm.DB
}

// NewLBRecordDao BaseDao constructor
func NewLBRecordDao(c *Config) (d *LBRecordDaoImpl) {
	d = &LBRecordDaoImpl{
		config:  c,
		db: NewMysql(c),
	}
	d.db.SingularTable(true)
	return
}

// Close release all mysql resource .
func (dao *LBRecordDaoImpl) Close() {
	dao.db.Close()
}

func (dao *LBRecordDaoImpl) GetLBRecord(id int64, record *model.LBRecord) (err error)  {
	err = dao.db.Where(&model.LBRecord{ID: id, Deleted:false}).First(record).Error
	if err != nil {
		dao.config.Logger.Error(err)
	}
	return
}

func (dao *LBRecordDaoImpl) ListLBRecord(conditions map[string]interface{}, records *[]model.LBRecord) (err error)  {
	conditions["deleted"] = false
	err = dao.db.Find(records, conditions).Error
	if err != nil {
		dao.config.Logger.Error(err)
	}
	return
}

func (dao *LBRecordDaoImpl) CreateLBRecord(record *model.LBRecord) (*model.LBRecord, error)  {
	err := dao.db.Create(record).Error
	if err != nil {
		dao.config.Logger.Error(err)
	}
	return record,err
}

func (dao *LBRecordDaoImpl) UpdateLBRecord(record *model.LBRecord) (*model.LBRecord, error)  {
	err := dao.db.Save(record).Error
	if err != nil {
		dao.config.Logger.Error(err)
	}
	return record,err
}

func (dao *LBRecordDaoImpl) UpdateAttribute(id int64, attributes map[string]interface{}) (*model.LBRecord, error) {
	record := &model.LBRecord{}
	err := dao.db.Where(&model.LBRecord{ID: id}).First(record).Error
	if err != nil {
		dao.config.Logger.Error(err)
		return record, err
	}
	err = dao.db.Model(record).Where(&model.LBRecord{ID: id}).Update(attributes).Error
	if err != nil {
		dao.config.Logger.Error(err)
		return record, err
	}
	return record, nil
}

func (dao LBRecordDaoImpl)DropLBRecord(id int64) error  {
	return dao.db.Model(&model.LBRecord{}).Where("lb_record_id = ?", id).Update("deleted", true).Error
}
