package model

// LBAgent 负载均衡代理 管理虚拟主机负载均衡起停等工作
type LBAgent struct {
	ID          int64     `json:"lb_agent_id" gorm:"column:lb_agent_id;auto_increment;primary_key"`
	Description ADCString `json:"description" gorm:"type:varchar(511)"` //描述
	LbPoolID    int64     `json:"lb_pool_id"  gorm:"column:lb_pool_id;not null;type: bigint(20)"`
	AgentType   int8      `json:"type,omitempty" gorm:"column:type;not null,type: tinyint(1)"`

	Deleted    bool    `json:"deleted" gorm:"type:tinyint(1);not null;default:0"` //是否删除
	CreateTime ADCTime `json:"create_time" gorm:"column:create_time"`             //创建时间
	ModifyTime ADCTime `json:"modify_time" gorm:"column:modify_time"`             //修改时间
	DeleteTime ADCTime `json:"delete_time" gorm:"column:delete_time"`             //删除时间
}

// TableName 返回数据库表名称
func (*LBAgent) TableName() string {
	return "lb_agent"
}

const (
	// AgentOK 状态OK
	AgentOK int32 = 0
	// AgentFail 状态出错
	AgentFail int32 = 1
	// AgentUnknown 状态未知
	AgentUnknown int32 = 3
)

// LBAgentStauts Agent实时状态信息, 定时写入到etcd中, key: /adc/lbagent/{agent id}/
type LBAgentStauts struct {
	// 对应LBMC中存储的LBAgent Id, agent启动时作为参数传入
	ID   int64
	Role string

	ControllerRPC string // agent controller rpc address for loadbalance control
	API           string // agent api rest url for status control

	HostIP   string `json:"ip"`
	HostName string

	Policies []LBPolicy

	TimeAt ADCTime `json:"created_time"`
}
