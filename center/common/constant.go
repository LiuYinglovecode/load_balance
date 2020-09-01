package common

// InfomerKeyPrefx agent 定时汇报到etcd中的key前缀
// 存在etcd中定时汇报的路径是/alb/agent/{agentid}
const InfomerKeyPrefx = "/alb/agent/"

// LBReqKeyPrefix center接受请求后存放在etcd中的key前缀
// 完整路径是/adc/lbrequest/{agentid}/timestamp
const LBReqKeyPrefix = "/adc/lbrequest/"
