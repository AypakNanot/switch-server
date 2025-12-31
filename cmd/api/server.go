package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/config/source/file"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"opt-switch/app/admin/models"
	"opt-switch/app/admin/router"
	"opt-switch/app/jobs"
	"opt-switch/common/database"
	"opt-switch/common/global"
	common "opt-switch/common/middleware"
	"opt-switch/common/middleware/handler"
	"opt-switch/common/storage"
	"opt-switch/web"
	ext "opt-switch/config"
)

var (
	configYml string
	apiCheck  bool
	StartCmd  = &cobra.Command{
		Use:          "server",
		Short:        "Start API server",
		Example:      "go-admin server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

var AppRouters = make([]func(), 0)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().BoolVarP(&apiCheck, "api", "a", false, "Start server with check api data")

	//注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
}

func setup() {
	// 注入配置扩展项
	config.ExtendConfig = &ext.ExtConfig
	//1. 读取配置
	config.Setup(
		file.NewSource(file.WithPath(configYml)),
		database.Setup,
		storage.Setup,
	)

	//2. 运行时内存优化（在配置读取后立即执行）
	initRuntime()

	//注册监听函数
	queue := sdk.Runtime.GetMemoryQueue("")
	queue.Register(global.LoginLog, models.SaveLoginLog)
	queue.Register(global.OperateLog, models.SaveOperaLog)
	queue.Register(global.ApiCheck, models.SaveSysApi)
	go queue.Run()

	usageStr := `starting api server...`
	log.Info(usageStr)
}

func run() error {
	if config.ApplicationConfig.Mode == pkg.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	initRouter()

	for _, f := range AppRouters {
		f()
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.ApplicationConfig.Host, config.ApplicationConfig.Port),
		Handler: sdk.Runtime.GetEngine(),
		ReadTimeout:  time.Duration(config.ApplicationConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ApplicationConfig.WriterTimeout) * time.Second,
	}

	go func() {
		jobs.InitJob()
		jobs.Setup(sdk.Runtime.GetDb())

	}()

	if apiCheck {
		var routers = sdk.Runtime.GetRouter()
		q := sdk.Runtime.GetMemoryQueue("")
		mp := make(map[string]interface{})
		mp["List"] = routers
		message, err := sdk.Runtime.GetStreamMessage("", global.ApiCheck, mp)
		if err != nil {
			log.Infof("GetStreamMessage error, %s \n", err.Error())
			//日志报错错误，不中断请求
		} else {
			err = q.Append(message)
			if err != nil {
				log.Infof("Append message error, %s \n", err.Error())
			}
		}
	}

	go func() {
		// 服务连接
		if config.SslConfig.Enable {
			if err := srv.ListenAndServeTLS(config.SslConfig.Pem, config.SslConfig.KeyStr); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("listen: ", err)
			}
		}
	}()
	fmt.Println(pkg.Red(string(global.LogoContent)))
	tip()
	fmt.Println(pkg.Green("Server run at:"))
	fmt.Printf("-  Local:   %s://localhost:%d/ \r\n", "http", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: %s://%s:%d/ \r\n", "http", pkg.GetLocalHost(), config.ApplicationConfig.Port)
	fmt.Println(pkg.Green("Swagger run at:"))
	fmt.Printf("-  Local:   http://localhost:%d/swagger/admin/index.html \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: %s://%s:%d/swagger/admin/index.html \r\n", "http", pkg.GetLocalHost(), config.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", pkg.GetCurrentTimeStr())
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Info("Shutdown Server ... ")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")

	return nil
}

//var Router runtime.Router

func tip() {
	usageStr := `欢迎使用 ` + pkg.Green(`go-admin `+global.Version) + ` 可以使用 ` + pkg.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s \n\n", usageStr)
}

func initRouter() {
	var r *gin.Engine
	h := sdk.Runtime.GetEngine()
	if h == nil {
		h = gin.New()
		sdk.Runtime.SetEngine(h)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		//os.Exit(-1)
	}
	if config.SslConfig.Enable {
		r.Use(handler.TlsHandler())
	}
	//r.Use(middleware.Metrics())

	// 条件启用中间件（用于内存优化）
	// Sentinel 限流中间件
	if getBoolConfig("application.enableMiddleware.sentinel", true) {
		r.Use(common.Sentinel())
	}
	// RequestID 中间件
	if getBoolConfig("application.enableMiddleware.requestID", true) {
		r.Use(common.RequestId(pkg.TrafficKey))
	}
	// RequestLogger 中间件（始终启用）
	r.Use(api.SetRequestLogger)

	common.InitMiddleware(r)

	// 条件启用前端静态文件（用于内存优化）
	if getBoolConfig("application.enableFrontend", true) {
		// 设置前端静态文件服务 - 直接使用 embed.FS 读取文件
		// 静态资源路由 (css, js, fonts, img 等)
	serveStaticFile := func(c *gin.Context, filePath string) {
		data, err := web.WebFS.ReadFile(filePath)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		// 设置正确的 Content-Type
		switch {
		case strings.HasSuffix(filePath, ".html"):
			c.Header("Content-Type", "text/html; charset=utf-8")
		case strings.HasSuffix(filePath, ".css"):
			c.Header("Content-Type", "text/css; charset=utf-8")
		case strings.HasSuffix(filePath, ".js"):
			c.Header("Content-Type", "application/javascript; charset=utf-8")
		case strings.HasSuffix(filePath, ".json"):
			c.Header("Content-Type", "application/json; charset=utf-8")
		case strings.HasSuffix(filePath, ".png"):
			c.Header("Content-Type", "image/png")
		case strings.HasSuffix(filePath, ".jpg"), strings.HasSuffix(filePath, ".jpeg"):
			c.Header("Content-Type", "image/jpeg")
		case strings.HasSuffix(filePath, ".svg"):
			c.Header("Content-Type", "image/svg+xml")
		case strings.HasSuffix(filePath, ".ico"):
			c.Header("Content-Type", "image/x-icon")
		case strings.HasSuffix(filePath, ".woff"), strings.HasSuffix(filePath, ".woff2"):
			c.Header("Content-Type", "font/woff2")
		}
		c.Data(http.StatusOK, "", data)
	}

	r.GET("/css/*filepath", func(c *gin.Context) {
		filePath := "dist/css" + c.Param("filepath")
		serveStaticFile(c, filePath)
	})
	r.GET("/js/*filepath", func(c *gin.Context) {
		filePath := "dist/js" + c.Param("filepath")
		serveStaticFile(c, filePath)
	})
	r.GET("/fonts/*filepath", func(c *gin.Context) {
		filePath := "dist/fonts" + c.Param("filepath")
		serveStaticFile(c, filePath)
	})
	r.GET("/img/*filepath", func(c *gin.Context) {
		filePath := "dist/img" + c.Param("filepath")
		serveStaticFile(c, filePath)
	})

	// favicon.ico
	r.GET("/favicon.ico", func(c *gin.Context) {
		serveStaticFile(c, "dist/favicon.ico")
	})

	// 根路径 - 返回 index.html
	r.GET("/", func(c *gin.Context) {
		serveStaticFile(c, "dist/index.html")
	})

	// SPA 路由支持：所有其他非 API 路由返回 index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 跳过 API 路径和 swagger
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "API Not Found"})
			return
		}
		if len(path) >= 8 && path[:8] == "/swagger" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Swagger Not Found"})
			return
		}
		if path == "/info" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
			return
		}
		// 返回 index.html
		serveStaticFile(c, "dist/index.html")
	})
	} // end of if enableFrontend

	// 如果未启用前端，根路径返回 API 信息
	if !getBoolConfig("application.enableFrontend", true) {
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "opt-switch API Server (Minimal Mode)",
				"data": gin.H{
					"version": global.Version,
					"mode":    "api-only",
					"docs":    "/swagger/admin/index.html",
				},
			})
		})
	}
}
