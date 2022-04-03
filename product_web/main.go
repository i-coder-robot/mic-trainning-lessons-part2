package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/internal"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/internal/register"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/product_web/handler"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/util"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var (
	consulRegistry register.ConsulRegistry
	randomId       string
)

func init() {
	randomPort := util.GenRandomPort()
	if !internal.AppConf.Debug {
		internal.AppConf.ProductWebConfig.Port = randomPort
	}
	randomId = uuid.NewV4().String()
	consulRegistry = register.NewConsulRegistry(internal.AppConf.ProductWebConfig.Host,
		internal.AppConf.ProductWebConfig.Port)
	consulRegistry.Register(internal.AppConf.ProductWebConfig.SrvName, randomId,
		internal.AppConf.ProductWebConfig.Port, internal.AppConf.ProductWebConfig.Tags)
}

func main() {
	ip := internal.AppConf.ProductWebConfig.Host
	port := util.GenRandomPort()
	if internal.AppConf.Debug {
		port = internal.AppConf.ProductWebConfig.Port
	}
	addr := fmt.Sprintf("%s:%d", ip, port)

	r := gin.Default()
	productGroup := r.Group("/v1/product")
	{
		productGroup.GET("/list", handler.ProductListHandler)
		productGroup.GET("/detail/:id", handler.DetailHandler)
		productGroup.POST("/add", handler.AddHandler)
		productGroup.POST("/update", handler.UpdateHandler)
		productGroup.POST("/del/:id", handler.DelHandler)
	}
	r.GET("/health", handler.HealthHandler)

	go func() {
		err := r.Run(addr)
		if err != nil {
			zap.S().Panic(addr + "启动失败" + err.Error())
		} else {
			zap.S().Info(addr + "启动成功")
		}
	}()
	q := make(chan os.Signal)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q
	err := consulRegistry.DeRegister(randomId)
	if err != nil {
		zap.S().Panic("注销失败" + randomId + ":" + err.Error())
	} else {
		zap.S().Info("注销成功:" + randomId)
	}
}
