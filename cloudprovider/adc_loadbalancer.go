package cloudprovider

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const configMapAnnotationKey = "k8s.co/cloud-provider-config"
const serviceForwardMethodAnnotationKey = "k8s.co/keepalived-forward-method"

//AdcLoadBalancer loadbalancer
type AdcLoadBalancer struct {
	kubeClient *kubernetes.Clientset
	cli        *LBController
}

var _ cloudprovider.LoadBalancer = &AdcLoadBalancer{}

//NewAdcBMLoadBalancer new AdcLoadBalancer struct
func NewAdcBMLoadBalancer(kubeClient *kubernetes.Clientset, cli *LBController) cloudprovider.LoadBalancer {
	return &AdcLoadBalancer{kubeClient, cli}
}

// TODO: Break this up into different interfaces (LB, etc) when we have more than one type of service

// GetLoadBalancer returns whether the specified load balancer exists, and
// if so, what its status is.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (a *AdcLoadBalancer) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	policies, err := a.cli.GetLoadBalancerFromEtcd(service.Namespace, service.Name)
	if err == nil && len(policies) != 0 {
		log.WithField("namepsace",service.Namespace).WithField("service", service.GetName()).
			Info("GetLoadBalancer success")
		return &v1.LoadBalancerStatus{
			Ingress: []v1.LoadBalancerIngress{{IP: policies[0].Record.IP.String()}},
		}, true, nil
	}
	log.WithField("namepsace", service.Namespace).WithField("service name", service.GetName()).
		WithField("reason", err).WithField("policies", policies).Info("GetLoadBalancer called failed")
	return nil, false, nil
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (a *AdcLoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	if len(nodes) == 0 {
		msg := "there are no available nodes for LoadBalancer service %s/%s"
		log.WithField("namespace", service.Namespace).WithField("service", service.Name).Errorf(msg)
		return nil, fmt.Errorf(msg)
	}
	status := &v1.LoadBalancerStatus{}
	loadBalancerIP, err := a.addLBReq(service, nodes, false)
	if err != nil || len(loadBalancerIP) == 0 {
		log.WithField("namespace", service.Namespace).WithField("service", service.Name).
			WithField("reason", err).Error("EnsureLoadBalancer called faild")
		return status, err
	}
	log.WithField("namespace", service.Namespace).WithField("service", service.Name).
		WithField("ip",loadBalancerIP).WithField("ports", service.Spec.Ports).
		Info("EnsureLoadBalancer called success")
	status.Ingress = []v1.LoadBalancerIngress{{IP: loadBalancerIP}}
	return status, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (a *AdcLoadBalancer) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	log.Infof("UpdateLoadBalancer: %s/%s", clusterName, service.GetName())
	if len(nodes) == 0 {
		msg := "there are no available nodes for LoadBalancer service %s/%s"
		log.WithField("namespace", service.Namespace).WithField("service", service.Name).Errorf(msg)
		return fmt.Errorf(msg)
	}
	ip, err := a.addLBReq(service, nodes, true)
	if err != nil || len(ip) == 0 {
		log.WithField("namespace", service.Namespace).WithField("service", service.Name).
			WithField("reason", err).Error("UpdateLoadBalancer called faild")
		return err
	}
	log.WithField("namespace", service.Namespace).WithField("service", service.Name).
		Info("UpdateLoadBalancer called success")
	return nil
}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it
// exists, returning nil if the load balancer specified either didn't exist or
// was successfully deleted.
// This construction is useful because many cloud providers' load balancers
// have multiple underlying components, meaning a Get could say that the LB
// doesn't exist even if some part of it is still laying around.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (a *AdcLoadBalancer) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	log.WithField("namespace", service.Namespace).WithField("service", service.Name).
		Info("EnsureLoadBalancerDeleted called")
	ps, err := a.cli.GetLoadBalancerFromEtcd(service.Namespace, service.Name)
	if err != nil {
		log.WithField("reason", err).Error("GetLoadBalancerFromEtcd failed")
		// 如果etcd中读取失败, 认为lb已经被删除, 不返回错误
		return nil
	}
	ns := service.Namespace
	svc := ns + "/" + service.Name

	reqs := a.buildDeleteLBrequestsFromPolicies(ns, svc, ps)

	log.Info("Sent Delete LBRequest to LBMC client------\n")
	err = a.sentAndWait(reqs, a.cli.DeleteLoadBalancer)
	if err != nil {
		return err
	}

	// 如果处理成功, 从etcd中删除key, key的格式为:
	// prefix/namespace/service
	err = a.cli.DeletePolicies(service.Namespace, service.Name)
	if err != nil {
		log.WithField("reason", err).Errorf("Delete LBPolicy from etcd failed")
		return err
	}
	log.Infof("EnsureLoadBalancerDeleted success, %s/%s\n", clusterName, service.Name)
	return nil
}

// The LB needs to be configured with instance addresses on the same
// subnet as the LB (aka opts.SubnetID).  Currently we're just
// guessing that the node's InternalIP is the right address - and that
// should be sufficient for all "normal" cases.
func (a *AdcLoadBalancer) nodeAddressForLB(node *v1.Node) (string, error) {
	addrs := node.Status.Addresses
	if len(addrs) == 0 {
		return "", fmt.Errorf("ErrNoAddressFound")
	}

	for _, addr := range addrs {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address, nil
		}
	}
	return addrs[0].Address, nil
}

// addLBReq call lbmc client
func (a *AdcLoadBalancer) addLBReq(service *v1.Service, nodes []*v1.Node, update bool) (string, error) {
	if !update {
		if ip, err := a.cli.CheckLbReady(service.Namespace, service.Name); err == nil && ip != "" {
			return ip, err
		}
	}

	reqs, policies, e := a.buildLBrequests(service, nodes, update)
	if e != nil {
		log.WithField("reason", e).Error("Build LBRequest failed")
		return "", e
	}

	log.Info("Sent Add LBRequest to LBMC client------\n")
	err := a.sentAndWait(reqs,a.cli.CreateLoadBalancer)
	if err != nil {
		return "", nil
	}

	// 如果处理成功，存储到etcd中, 格式为 prefix/namespace/service
	err = a.cli.PersistPolicies(service.Namespace, service.Name, &policies)
	if err != nil {
		log.WithField("reason", err).Errorf("Persist LBPolicy into etcd failed")
		return "", err
	}
	log.Infof("addLoadBalancer success, %s\n", service.Name)
	return service.Spec.LoadBalancerIP, err
}

// buildLBrequests 根据service信息 构建lbrequest请求
func (a *AdcLoadBalancer) buildLBrequests(service *v1.Service, nodes []*v1.Node, update bool) ([]*model.LBRequest, []model.LBPolicy, error) {
	ns := service.Namespace
	svc := ns + "/" + service.Name
	loadBalancerIP := service.Spec.LoadBalancerIP

	// todo: 目前只支持单独IP的负载均衡
	// 对每一个port生成一个请求
	var reqs []*model.LBRequest
	var policies []model.LBPolicy

	ports := service.Spec.Ports
	for _, port := range ports {
		realServers := []model.RealServer{}
		for _, node := range nodes {
			addr, err := a.nodeAddressForLB(node)
			if err != nil {
				continue
			}

			rs := model.RealServer{
				Name: node.Name,
				IP:   addr,
				Port: int32(port.NodePort),
			}
			realServers = append(realServers, rs)
		}

		lbr := model.NewLBRecordIP("", loadBalancerIP, port.Port)
		po := &model.LBPolicy{
			lbr,
			realServers,
		}

		var i model.ActionType
		if update {
			i = model.ActionUpdate
		} else {
			i = model.ActionAdd
		}
		req := model.NewLBRequest(ns, svc, i, po)
		reqs = append(reqs, req)
		policies = append(policies, *po)
	}
	if len(reqs) == 0 {
		return nil, nil, fmt.Errorf("ErrNoServicePort")
	}
	return reqs, policies, nil
}

func (a *AdcLoadBalancer) buildDeleteLBrequestsFromPolicies(namespace, svc string, ps []model.LBPolicy) ([]*model.LBRequest) {
	var reqs []*model.LBRequest
	for _, v := range ps {
		req := model.NewLBRequest(namespace, svc, model.ActionDelete, &v)
		reqs = append(reqs, req)
	}
	return reqs
}

// sentAndWait 向lbmc发送请求, 等待请求处理返回结果
// 请求处理成功返回nil
func (a *AdcLoadBalancer) sentAndWait(reqs []*model.LBRequest, f func (request *model.LBRequest) (string, error)) error {
	done := make(chan error)
	// 发送请求
	for _, v := range reqs {
		go func(r *model.LBRequest) {
			var err error
			reqID, err := f(r)
			if err != nil {
				log.WithField("reason", err).WithField("ip", r.Policy.Record.IP.String()).
					WithField("port", r.Policy.Record.Port).Error("Send lbrequest to lbmc failed")
				done <- err
				return
			}

			err = a.cli.waitforLbReady(reqID)
			if err != nil {
				log.WithField("reason", err).WithField("ip", r.Policy.Record.IP.String()).
					WithField("port", r.Policy.Record.Port).Error("Wait request for ready failed")
			}
			done <- err
		}(v)
	}

	// 阻塞
	var err error
	for i := 0; i < len(reqs); i++ {
		err = <-done
		if err != nil {
			log.WithField("err", err).Errorf("waitforLbReady got error")
			return  err
		}
	}
	return nil
}