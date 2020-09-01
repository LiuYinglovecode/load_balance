package dao

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"testing"
)

const sql = `CREATE TABLE lb_request (
  lb_request_id INTEGER PRIMARY KEY,
  owner varchar(255),
  service varchar(255),

  status varchar(255),
  action tinyint(1),
  policy text,


  finish_time datetime DEFAULT NULL,
  create_time datetime,
  update_time datetime DEFAULT NULL,
  delete_time datetime DEFAULT NULL,
  note varchar(511)
)`

func TestLBRequest_CRUD(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	r := NewLBRequestRepository(db)

	request := &model.LBRequest{
		User: model.NewADCString("12345"),
		Status: 0,
		Action: 1,
		Policy: model.LBPolicy{
			Record: model.LBRecord{
				IP:   model.NewADCString("192.168.100.200"),
				Port: 80,
				Type: model.TypeIP},
			Endpoints: []model.RealServer{
				{Name: "sever1", IP: "106.74.100.99", Port: 80},
				{Name: "sever2", IP: "106.74.100.98", Port: 80},
				{Name: "sever3", IP: "106.74.100.97", Port: 80},
			},
		},
	}

	//create table
	r.DB().Exec(sql)

	// create
	err = r.Create(request);
	assert.NoError(t, err)
	assert.NotEqual(t, request.RequestID, 0)
	id := request.RequestID
	// get
	nr, err := r.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, "12345", nr.User.String())
	// update
	nr.User = model.NewADCString("23456")
	err = r.Update(nr)
	assert.NoError(t, err)
	assert.Equal(t, nr.User.String(), "23456")
	// delete
	err = r.DeleteByID(id)
	assert.NoError(t, err)
}

