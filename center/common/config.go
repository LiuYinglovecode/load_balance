package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"sync"
	"time"
)

// Config for center
type Config struct {
	WorkDir  string // 工作目录，存储输出日志，Agent持久化的相关内容
	LogLevel int    // 日志输出级别

	AuditLogPath string // audit log path
	SysLogPath   string // system log path

	AgentTimeout time.Duration // 连接agent超时时间

	Port int // HttpServer启动监听的端口号

	DBArgs string   // 数据库连接参数, eg. id:password@tcp(your-host-uri.com:3306)/dbname?foo=bar
	Dialect string  // 数据库类型, eg. mysql

	EtcdEndpoints []string //etcd 集群地址

	EtcdCAPath   string
	EtcdCertPath string
	EtcdKeyPath  string
}

// 全局配置
// todo: 设置为imuutable
var GlobalConfig Config

// db 是连接数据库使用的datasource
var db *gorm.DB

var once sync.Once

// GetDataSource 获取当前数据库连接
// singleton pattern
func GetDataSource() *gorm.DB {
	once.Do(func() {
		db = initDatabase(&GlobalConfig)
	})
	return db
}

func initDatabase(config *Config) *gorm.DB {
	db, err := gorm.Open(config.Dialect, config.DBArgs)
	if err != nil {
		// 连接数据库失败, 推出程序
		SysLogger.Errorf("connect to database failed, reason: %v", err)
		os.Exit(1)
	}

	// todo: 改为可配置
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.DB().SetConnMaxLifetime(time.Hour)

	return db
}