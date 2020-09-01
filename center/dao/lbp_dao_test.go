package dao

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLBPoolDaoImpl_CreateLBPool(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBPoolDaoImpl(config)
	defer dao.Close()
	pool := &model.LBPool{
		Type:        model.TypeDomain,
		DomainRegex: model.NewADCString("[0-9a-Z]{1,2}.test.cn"),
	}
	record, err := dao.CreateLBPool(pool)

	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, record.ID)

	pool = &model.LBPool{
		Type:      model.TypeIP,
		IP:        model.NewADCString("106.75.69.9"),
		StartPort: 18000,
		EndPort:   19000,
	}
	record, err = dao.CreateLBPool(pool)

	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, record.ID)
}

func TestLBPoolDaoImpl_GetLBPool(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBPoolDaoImpl(config)
	defer dao.Close()
	var id = int64(1)
	pool := model.LBPool{}
	err := dao.GetLBPool(id, &pool)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, pool)
}
func TestLBPoolDaoImpl_ListLBPool(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBPoolDaoImpl(config)
	defer dao.Close()
	var pools []model.LBPool
	err := dao.ListLBPool(map[string]interface{}{"ip": "106.74.152.34"}, &pools)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, pools)
}

func TestLBPoolDaoImpl_UpdateLBPool(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBPoolDaoImpl(config)
	defer dao.Close()
	pool := &model.LBPool{
		ID:2,
		IP:        model.NewADCString("106.75.69.9"),
		StartPort: 8000,
		EndPort:   8999,
	}
	record, err := dao.UpdateLBPool(pool)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, record.ID)
}

func TestLBPoolDaoImpl_UpdateAttribute(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBPoolDaoImpl(config)
	defer dao.Close()
	id := int64(4)
	attributes := map[string]interface{}{
		"start_port": 9000,
		"end_port": 9999,
	}
	record, err := dao.UpdateAttribute(id, attributes)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, record.ID)
}

func TestLBPoolDaoImpl_DropLBPool(t *testing.T) {
	var config = getDBConfig()
	dao := NewLBPoolDaoImpl(config)
	defer dao.Close()
	var id = int64(1)
	err := dao.DropLBPool(id)
	if err != nil {
		t.Error(err)
	}
}