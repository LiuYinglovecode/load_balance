package center

import (
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/center/web"
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// StartServer 启动http服务
func StartServer() error {
	port := ":" + strconv.Itoa(common.GlobalConfig.Port)

	router, err:= web.SetupRouter()
	if err != nil {
		common.SysLogger.Errorf("Failed setup router, reason: %v", err)
		return err
	}

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.SysLogger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	common.SysLogger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		common.SysLogger.Fatalf("Server Shutdown: %v", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		common.SysLogger.Info("timeout of 5 seconds.")
	}
	common.SysLogger.Info("Server exiting")
	return nil
}



