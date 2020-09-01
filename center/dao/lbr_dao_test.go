package dao

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLBRecordDaoImpl_CreateLBRecord(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBRecordDao(config)
	defer dao.Close()
	record := &model.LBRecord{
		Owner: model.NewADCString("10000031112923"),
		Type:model.TypeDomain,
		Domain:model.NewADCString("test3"),
	}

	record, err := dao.CreateLBRecord(record)

	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, record.ID)
}

func TestLBRecordDaoImpl_GetLBRecord(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBRecordDao(config)
	defer dao.Close()
	var id = int64(1)
	record := model.LBRecord{}
	err := dao.GetLBRecord(id, &record)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, record)
}

func TestLBRecordDaoImpl_LISTLBRecord(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBRecordDao(config)
	defer dao.Close()
	var records []model.LBRecord
	err := dao.ListLBRecord(map[string]interface{}{"owner": "10000031112923", "domain": "test"}, &records)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, records)
}

func TestLBRecordDaoImpl_UpdateLBRecord(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBRecordDao(config)
	defer dao.Close()
	record := &model.LBRecord{
		ID:1,
		Owner: model.NewADCString(10000031112923),
		Type:model.TypeDomain,
		Name:"test",
		Domain:model.NewADCString("test.domain"),
	}

	record, err := dao.UpdateLBRecord(record)

	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, record.ID)
}

func TestLBRecordDaoImpl_UpdateAttribute(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBRecordDao(config)
	defer dao.Close()
	id := int64(2)
	attributes := map[string]interface{}{
		"name": "test1",
		"domain": "test1.domain",
	}
	record, err := dao.UpdateAttribute(id, attributes)

	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, record.ID)
}

func TestLBRecordDaoImpl_DropLBRecord(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBRecordDao(config)
	defer dao.Close()
	var id = int64(1)
	err := dao.DropLBRecord(id)
	if err != nil {
		t.Error(err)
	}
}

func getDBConfig() *Config {
	return &Config{
		DSN: "root:123456@tcp(172.17.60.108:3306)/adc?charset=utf8&parseTime=True",
	}
}