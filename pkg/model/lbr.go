package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// LBType type of lbreocrd
type LBType = int

const (
	// TypeIP 负载均衡资源是IP端口号类型的。如 106.75.69.124:5050
	TypeIP LBType = 0
	// TypeDomain 负载均衡类型是域名形式的。如 adc.htres.cn
	TypeDomain LBType = 1
)

// LBRStatus type define
type LBRStatus = int32

const (
	//Applying LBR状态: 用户已经申请，但是未获审批，
	Applying LBRStatus = 0
	//Applied LBR状态: 审批通过，但client端未使用
	Applied LBRStatus = 1
	//InUse LBR状态: 代理通道已经建立，正常使用中
	InUse LBRStatus = 2
)

// LBRecord 负载均衡资源
type LBRecord struct {
	ID   int64  `json:"lb_record_id,omitempty" gorm:"column:lb_record_id;auto_increment;primary_key"`
	Type LBType `json:"type,omitempty" gorm:"type:tinyint(1);not null;default:0"` // LBType

	Owner ADCString `json:"owner,omitempty" gorm:"not null;type:varchar(255)"` // 资源所有者，用户中心的用户ID

	IP     ADCString `json:"ip,omitempty" gorm:"type:varchar(511)"`     // public ip address include ipv4 and ipv6 address
	Port   int32     `json:"port,omitempty" gorm:"type:int"`            // port for loadbalance
	Domain ADCString `json:"domain,omitempty" gorm:"type:varchar(511)"` // domain name for using

	Deleted    bool    `json:"deleted,omitempty" gorm:"type:tinyint(1);not null;default:0"` //是否删除
	CreateTime ADCTime `json:"create_time,omitempty"`                                       //创建时间
	ModifyTime ADCTime `json:"modify_time,omitempty"`

	Status LBRStatus `json:"status,omitempty"`
	Name   string `json:"name,omitempty"` // 用户对改负载均衡的描述
}

// NewLBRecord 生成IP类型的LBR
func NewLBRecordIP(owner string, ip string, port int32) LBRecord {
	return LBRecord{
		Owner: NewADCString(owner),
		IP: NewADCString(ip),
		Port: port,
		Type: TypeIP,
	}
}
// NewLBRecordDomain 生成Domain类型的LBR
func NewLBRecordDomain(owner int64, domain string) LBRecord {
	return LBRecord{
		Owner: NewADCString(owner),
		Domain: NewADCString(domain),
		Type: TypeDomain,
	}
}

// LBPool 负载均衡资源池。形式如下：
// 106.75.69.9,8000-20000  亦或 [0-9a-Z]{0,1}.htres.cn [0-9a-Z]{1,2}.htres.cn
type LBPool struct {
	ID   int64  `json:"lb_pool_id,omitempty" gorm:"column:lb_pool_id;auto_increment;primary_key"`
	Type LBType `json:"type,omitempty" gorm:"type:tinyint(1);not null;default:0"` // LBType

	IP          ADCString `json:"ip,omitempty" gorm:"type:varchar(511)"`           // public ip address include ipv4 and ipv6 address
	StartPort   int32  `json:"start_port,omitempty" gorm:"type:int"`            // port start
	EndPort     int32  `json:"end_port,omitempty" gorm:"type:int"`              // port end
	DomainRegex ADCString `json:"domain_regex,omitempty" gorm:"type:varchar(511)"` // domain regex

	Agents []LBAgent   `json:"agents,omitempty" gorm:"foreignkey:LbPoolID;association_foreignkey:ID"`           // one to many agents

	Deleted    bool    `json:"deleted,omitempty" gorm:"type:tinyint(1);not null;default:0"` //是否删除
	CreateTime ADCTime `json:"create_time,omitempty"`                                       //创建时间
	ModifyTime ADCTime `json:"modify_time,omitempty"`                                       //修改时间
}

// RealServer 反向代理后端的真实服务器
type RealServer struct {
	IP   string `json:"ip,omitempty"`
	Port int32  `json:"port,omitempty"`
	Name string `json:"name,omitempty"`
}

// LBPolicy 反向代理策略，包括公网IP,以及后端IP,今后会把复杂的策略添加进来，主要为了生成HAProxy配置文件
type LBPolicy struct {
	Record    LBRecord     `json:"record,omitempty"`
	Endpoints []RealServer `json:"endpoints,omitempty"`
}

//GetID get idenity string of LBPolicy
func (p LBPolicy) GetID() string {
	var sb strings.Builder
	switch p.Record.Type {
	case TypeIP:
		sb.WriteString(p.Record.IP.String())
		sb.WriteString("_")
		sb.WriteString(strconv.Itoa(int(p.Record.Port)))
	case TypeDomain:
		sb.WriteString(p.Record.Domain.String())
	default:
		sb.WriteString("unknown_")
	}

	return sb.String()
}

// Equals two policy if every thing is the same
func (p LBPolicy) Equals(b LBPolicy) bool {
	result := true
	result = result && p.Record.ID == b.Record.ID
	result = result && p.GetID() == b.GetID()
	result = result && len(p.Endpoints) == len(b.Endpoints)
	for i, v := range p.Endpoints {
		if v != b.Endpoints[i] {
			return false
		}
	}
	return result
}

// Value implements the driver Valuer interface.
func (l LBPolicy) Value() (driver.Value, error) {
	jsonData, err := json.Marshal(l)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Scan implements the Scanner interface.
func (l *LBPolicy) Scan(value interface{}) error {
	if value == nil {
		l.Record, l.Endpoints = LBRecord{}, []RealServer{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		err := json.Unmarshal(v, l)
		if err != nil {
			return err
		}
		return nil
	case string:
		err := json.Unmarshal([]byte(v), l)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("can not scan %v to LBPlicy", value)
}
