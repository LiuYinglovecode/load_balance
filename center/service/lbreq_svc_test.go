package service

import (
	"code.htres.cn/casicloud/alb/apis"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestHandlerImpl_AcceptRequest(t *testing.T) {
	h := NewRequestHandlerImpl(&LBRecordRepositoryStub{}, &LBPoolRepositoryStub{}, &LBRequestRepositoryStub{}, &MessageQueueHandlerStub{}, nil)

	request := &model.LBRequest{
		User:   model.NewADCString("12345"),
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

	got, err := h.AcceptRequest(request)
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestRequestHandlerImpl_QueryRequest(t *testing.T) {
	h := NewRequestHandlerImpl(&LBRecordRepositoryStub{}, &LBPoolRepositoryStub{}, &LBRequestRepositoryStub{}, &MessageQueueHandlerStub{}, nil)
	got, err := h.QueryRequest(3827)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), got)
}

var i = 0

func TestRequestHandlerImpl_ProcessRequest(t *testing.T) {
	h := NewRequestHandlerImpl(&LBRecordRepositoryStub{}, &LBPoolRepositoryStub{}, &LBRequestRepositoryStub{}, &MessageQueueHandlerStub{}, nil)
	h.ProcessRequest()
	assert.Equal(t, i, 2)
}

// MessageQueueHandlerStub  stub
type MessageQueueHandlerStub struct {
	// blank
}

// WatchAndDequeue stub
func (m *MessageQueueHandlerStub) WatchAndDequeue() (*model.LBRequest, string) {
	return &model.LBRequest{
		RequestID: 9898,
		User:      model.NewADCString("12345"),
		Status:    0,
		Action:    1,
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
	}, "1"
}

// Enqueue stub
func (m *MessageQueueHandlerStub) Enqueue(queueName string, request *model.LBRequest) error {
	return nil
}

// LBRequestRepositoryStub stub
type LBRequestRepositoryStub struct {
	// blank
}

func (l *LBRequestRepositoryStub) DB() *gorm.DB {
	panic("implement me")
}

// Create stub
func (l *LBRequestRepositoryStub) Create(request *model.LBRequest) error {
	request.RequestID = 9898
	return nil;
}

// GetById stub
func (l *LBRequestRepositoryStub) GetByID(id int64) (*model.LBRequest, error) {
	return &model.LBRequest{
		RequestID: 9898,
		User:      model.NewADCString("12345"),
		Status:    0,
		Action:    1,
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
	}, nil
}

// Update stub
func (l *LBRequestRepositoryStub) Update(request *model.LBRequest) error {
	i = i + 1
	return nil
}

// DeleteById stub
func (l *LBRequestRepositoryStub) DeleteByID(id int64) error {
	panic("implement me")
}

// ConmmunicatorStub  实现Conmunicator接口
type ConmmunicatorStub struct {
}

// SendLBRequest stub
func (c *ConmmunicatorStub) SendLBRequest(agentRPC string, request *model.LBRequest) (*apis.LBCommandResult, error) {
	return &apis.LBCommandResult{Code: 200, Msg: "ok"}, nil
}

// SendLBCommand stub
func (c *ConmmunicatorStub) SendLBCommand(agentRPC string, cmd *apis.LBCommand) (*apis.LBCommandResult, error) {
	return &apis.LBCommandResult{Code: 200, Msg: "ok"}, nil
}

type LBRecordRepositoryStub struct {
}

func (*LBRecordRepositoryStub) GetLBRecord(id int64, record *model.LBRecord) error {
	panic("implement me")
}

func (*LBRecordRepositoryStub) ListLBRecord(conditions map[string]interface{}, records *[]model.LBRecord) error {
	lbr := model.LBRecord{
		Owner:  model.NewADCString("12345"),
		IP:     model.NewADCString("127.0.0.1"),
		Port:   1023,
		Type:   0,
		Status: model.Applied,
	}

	*records = append(*records, lbr)
	return nil
}

func (*LBRecordRepositoryStub) CreateLBRecord(record *model.LBRecord) (*model.LBRecord, error) {
	panic("implement me")
}

func (*LBRecordRepositoryStub) UpdateLBRecord(record *model.LBRecord) (*model.LBRecord, error) {
	panic("implement me")
}

func (*LBRecordRepositoryStub) UpdateAttribute(id int64, attributes map[string]interface{}) (*model.LBRecord, error) {
	return nil, nil
}

func (*LBRecordRepositoryStub) DropLBRecord(id int64) error {
	panic("implement me")
}

type LBPoolRepositoryStub struct {
}

func (*LBPoolRepositoryStub) GetLBPool(id int64, Pool *model.LBPool) error {
	panic("implement me")
}

func (*LBPoolRepositoryStub) ListLBPool(conditions map[string]interface{}, pools *[]model.LBPool) error {
	pool := model.LBPool{
		ID:        1,
		IP:        model.NewADCString("127.0.0.1"),
		StartPort: 10,
		EndPort:   30000,
		Agents: []model.LBAgent{
			{ID: 1},
		},
	}
	*pools = append(*pools, pool)
	return nil
}

func (*LBPoolRepositoryStub) CreateLBPool(Pool *model.LBPool) (*model.LBPool, error) {
	panic("implement me")
}

func (*LBPoolRepositoryStub) UpdateLBPool(Pool *model.LBPool) (*model.LBPool, error) {
	panic("implement me")
}

func (*LBPoolRepositoryStub) UpdateAttribute(id int64, attributes map[string]interface{}) (*model.LBPool, error) {
	panic("implement me")
}

func (*LBPoolRepositoryStub) DropLBPool(id int64) error {
	panic("implement me")
}
