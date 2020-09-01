package queue

import "code.htres.cn/casicloud/alb/pkg/model"

// MessageQueueHandler 实现消息队列读写接口
type MessageQueueHandler interface {
	MessageQueueReader
	MessageQueueWriter
}

// MessageQueueReader 读消息队列接口
type MessageQueueReader interface {
	// 监听全部消息队列,
	// 如果消息存在则出列处理, 否则阻塞等待
	// 同时返回队列名称
	WatchAndDequeue() (*model.LBRequest, string)
}

// MessageQueueWriter 写消息队列接口
type MessageQueueWriter interface {
	// 将请求加入指定队列,
	// 如果队列已满则返回错误信息
	Enqueue(queueName string, request *model.LBRequest) error
}

