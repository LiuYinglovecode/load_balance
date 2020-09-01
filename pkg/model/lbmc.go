package model

// ActionType command type for LBCommand
type ActionType = int32

const (
	//ActionAdd request type
	ActionAdd ActionType = 1
	//ActionUpdate request type
	ActionUpdate ActionType = 2
	//ActionDelete request type
	ActionDelete ActionType = 3
	//ActionGet request type
	ActionGet ActionType = 4
	//ActionStop request type
	ActionStop ActionType = 5

	//StatusUnHandle Request 未处理
	StatusUnHandle int32 = 0
	//StatusInProcessing Request 处理中
	StatusInProcessing int32 = 1
	//StatusHandleSuccess Request 处理成功
	StatusHandleSuccess int32 = 2
	//StatusHandleFailed Request 处理失败
	StatusHandleFailed int32 = 3
	//StatusUnknown  Request 处理情况未知
	StatusUnknown int32 = 999
)

// LBRequest client向LBMC发起的请求
// LBMC通过rest接口接收, 将请求信息存入etcd中,
// 如果通过k8s申请负载均衡需要传入namespace信息
// etcd作为消息队列如何设计key值
// key: /adc/lbreq/{action type}/{domain或者ip:port}
// proteus:generate
type LBRequest struct {
	RequestID int64  `json:"request_id,omitempty" gorm:"column:lb_request_id;auto_increment;primary_key"`
	User      ADCString  `json:"user,omitempty" gorm:"column:owner;type:varchar(255)"`
	Service   ADCString `json:"service,omitempty" gorm:"column:service;type:varchar(511)"`

	Status int32 `json:"status"`
	Action int32 `json:"action"`

	//Policy LBPolicy `json:"policy,omitempty" gorm:"column:policy;type:text"`
	Policy LBPolicy `json:"policy,omitempty" gorm:"column:policy;type:text"`

	CreatedAt ADCTime  `json:"create_time,omitempty" gorm:"column:create_time"`
	UpdatedAt ADCTime  `json:"modify_time,omitempty"  gorm:"column:update_time"`
	FinishAt  ADCTime  `json:"finish_time,omitempty"  gorm:"column:finish_time"`
	DeletedAt ADCTime `json:"delete_time,omitempty"  gorm:"column:delete_time"`
}

func (*LBRequest) String() string {
	panic("implement me")
}

// TableName 返回数据库表名称
func (*LBRequest) TableName() string {
	return "lb_request"
}

// NewLBRequest 构造函数
// user是申请者信息, 非k8s申请，填写userid
// k8s申请 填写namespace名称
// service 是申请负载均衡的服务名称
// 非k8s申请 service可以传空字符串
func NewLBRequest(user string, service string, action int32, policy *LBPolicy) *LBRequest {
	return &LBRequest{
		User:    NewADCString(user) ,
		Action:  action,
		Policy:  *policy,
		Service: NewADCString(service),
	}
}
