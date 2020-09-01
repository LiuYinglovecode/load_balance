package agent

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	//审计日志 负责记录center各种请求以及处理过程
	auditLogger = logrus.New()
	//系统日志 负责记录程序运行情况
	sysLogger = logrus.New()
)

// InitLogger 初始化日志
func InitLogger(config *Config) error {
	// TODO: 请参考 http://xiaorui.cc/2018/01/11/golang-logrus%E7%9A%84%E9%AB%98%E7%BA%A7%E9%85%8D%E7%BD%AEhook-logrotate/ 实现日志rotate
	auditLogger.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	sysLogger.SetFormatter(&logrus.TextFormatter{DisableColors: true})

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}

	sysLogger.SetLevel(level)
	//审计日志只需要用info就可以了
	auditLogger.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile(config.AuditLogPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		auditLogger.Out = file
	} else {
		return err
	}

	sysLogFile, err := os.OpenFile(config.SysLogPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		sysLogger.Out = sysLogFile
	} else {
		return err
	}
	return nil
}
