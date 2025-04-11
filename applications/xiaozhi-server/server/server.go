package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/loader"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/loader/resource"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/middleware"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/router"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(configFile, env string) error {
	err := config.Instance.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	myFigure := figure.NewFigure(config.Instance.AppName, "", true)
	myFigure.Print()

	ctx := context.Background()
	ginMode := map[string]string{
		"dev":  gin.DebugMode,
		"prod": gin.ReleaseMode,
	}
	if mode, ok := ginMode[env]; ok {
		gin.SetMode(mode)
	}
	r := gin.Default()
	go logger.InitLog()
	r.Use(middleware.CustomRecovery())

	r.Use(logger.HttpLog(config.Instance.Log))
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.PrometheusMiddleware())

	// 初始化资源 如数据库、redis等
	err = loader.InitResource()
	if err != nil {
		logger.CtxError(&gin.Context{}, "初始化资源失败", err.Error())
		return errors.New("初始化资源失败: " + err.Error())
	}
	defer func() {
		readDB, _ := resource.GetResource().ReadDB.DB()
		readDB.Close()
		writeDB, _ := resource.GetResource().WriteDB.DB()
		writeDB.Close()
		resource.GetResource().RedisClient.Close()
		resource.GetResource().TTSRpcClient.Conn.Close()
	}()

	_ = router.InitRouter(ctx, r)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	addr := fmt.Sprintf(":%d", config.Instance.Server.Port)
	fmt.Printf("Starting server at %s:...\n", addr)
	logger.CtxInfo(ctx, "Starting server at %s:...", addr)
	err = r.Run(addr)
	if err != nil {
		return errors.New("服务启动失败: " + err.Error())
	}

	return nil
}
