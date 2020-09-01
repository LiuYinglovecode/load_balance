package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	//AuditLogger 审计日志 负责记录center各种请求以及处理过程
	AuditLogger = logrus.New()
	//SysLogger 系统日志 负责记录程序运行情况
	SysLogger = logrus.New()
)

// InitLogger 初始化日志
func InitLogger(config *Config) error {
	// TODO: 请参考 http://xiaorui.cc/2018/01/11/golang-logrus%E7%9A%84%E9%AB%98%E7%BA%A7%E9%85%8D%E7%BD%AEhook-logrotate/ 实现日志rotate
	AuditLogger.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	SysLogger.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	file, err := os.OpenFile(config.AuditLogPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		AuditLogger.Out = file
	} else {
		return err
	}

	SysLogger.Out = os.Stdout
	// sysLogFile, err := os.OpenFile(config.SysLogPath, os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// SysLogger.Out = sysLogFile
	// } else {
	// return err
	// }
	return nil
}
