package agent

// Config for agent controller
type Config struct {
	RPC           string   // listened address for rpc controller
	WorkDir       string   // 工作目录，存储输出日志，Agent持久化的相关内容
	AgentID       int64    // lbagent id
	LBMCUrl       string   // lbmc url
	LogLevel      string   // loglevel
	AuditLogPath  string   // audit log path
	SysLogPath    string   // system log path
	Endpoints     []string // etcd 地址
	EtcdCAPath    string
	EtcdCertPath  string
	EtcdKeyPath   string
	KeepalivedCfg *KeepalivedConfig // keepalive的配置类
}

// KeepalivedConfig keepalived相关配置
type KeepalivedConfig struct {
	INet            string
	VirutalRouterID string
	State           string
	Priority        int
	UnicastSrcIP    string
	UnicastPeer     []string
	VirtualIP       string
}
