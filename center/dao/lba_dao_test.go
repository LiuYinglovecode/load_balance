package dao

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

const lbagent_sql = `
CREATE TABLE lb_agent (
  lb_agent_id Integer PRIMARY KEY,
  lb_pool_id  Integer NOT NULL,
 
  type tinyint(1) NOT NULL DEFAULT '0',
  description varchar(511),

  deleted tinyint(1) NOT NULL DEFAULT '0',
  create_time datetime,
  modify_time datetime DEFAULT NULL,
  delete_time datetime DEFAULT NULL
)`


func TestLBAgentRepositoryImpl_CRUD(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	r := NewLBAgentRepository(db)

	agent := &model.LBAgent{
		Description: model.NewADCString("wow"),
		LbPoolID: 1,
	}

	//create table
	r.DB().Exec(lbagent_sql)

	// create
	err = r.Create(agent);
	assert.NoError(t, err)
	assert.NotEqual(t, agent.ID, 0)
	id := agent.ID
	// get
	nr, err := r.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, "wow", nr.Description.String())
	// list
	agents, err := r.List(make(map[string]interface{}))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(agents))

	// update
	s := "wow and wow"
	nr.Description = model.NewADCString(s)
	err = r.Update(nr)
	assert.NoError(t, err)
	assert.Equal(t, nr.Description.String(), s)
	// delete
	err = r.DeleteByID(id)
	assert.NoError(t, err)
}

