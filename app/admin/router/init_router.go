package router

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"go.uber.org/zap"

	common "opt-switch/common/middleware"
	devicerouter "opt-switch/app/device/router"
)

// InitRouter 路由初始化，不要怀疑，这里用到了
func InitRouter() {
	var r *gin.Engine
	h := sdk.Runtime.GetEngine()
	if h == nil {
		log.Fatal("not found engine...")
		os.Exit(-1)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}

	// the jwt middleware
	authMiddleware, err := common.AuthInit()
	if err != nil {
		log.Fatalf("JWT Init Error, %s", err.Error())
	}

	// 注册系统路由
	g := InitSysRouter(r, authMiddleware)

	// 初始化设备服务 - create a zap logger for the device service
	deviceLogger, err := zap.NewProduction()
	if err != nil {
		log.Warnf("Failed to create device logger: %v", err)
	} else {
		if err := devicerouter.InitDeviceService(deviceLogger); err != nil {
			log.Errorf("Failed to initialize device service: %v", err)
			// Continue even if device service fails to start
		}
	}

	// 注册设备路由到 /api/v1
	devicerouter.InitDeviceRouter(g)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由只有演示代码
	InitExamplesRouter(r, authMiddleware)
}
