package agent

import "code.htres.cn/casicloud/alb/pkg/model"

// LoadBalancerHandler 负责创建负载均衡的业务逻辑
// 使用ProxyControler实现代理的创建、删除等功能
type LoadBalancerHandler interface {
	// 添加新的负载均衡代理
	// lbname 是负载均衡的名称，命名规则是 ???
	// portMap 是负载均衡的端口映射规则
	// 添加成功后返回成功信息， 目前架构应该返回container id
	AddLoadBalancer(lbname string, policy []model.LBPolicy) (string, error)

	UpdateLoadBalancer()

	DeleteLoadBalancer()
}
