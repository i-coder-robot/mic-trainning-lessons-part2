package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/internal"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/product_srv/biz"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/proto/pb"
	"github.com/i-coder-robot/mic-trainning-lessons-part2/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func init() {
	internal.InitDB()
}

func main() {

	/*
		1 生成proto对应的文件
		2 建立biz目录，实现接口
		3 拷贝之前main文件中的函数
	*/

	port := util.GenRandomPort()
	if internal.AppConf.Debug {
		port = internal.AppConf.ProductSrvConfig.Port
	}
	addr := fmt.Sprintf("%s:%d", internal.AppConf.ProductSrvConfig.Host, port)
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, &biz.ProductServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Error("account_srv启动异常:" + err.Error())
		panic(err)
	}
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d",
		internal.AppConf.ConsulConfig.Host,
		internal.AppConf.ConsulConfig.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		panic(err)
	}
	checkAddr := fmt.Sprintf("%s:%d", internal.AppConf.ProductSrvConfig.Host, port)
	check := &api.AgentServiceCheck{
		GRPC:                           checkAddr,
		Timeout:                        "3s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5s",
	}
	randUUID := uuid.New().String()
	reg := api.AgentServiceRegistration{
		Name:    internal.AppConf.ProductSrvConfig.SrvName,
		ID:      randUUID,
		Port:    port,
		Tags:    internal.AppConf.ProductSrvConfig.Tags,
		Address: internal.AppConf.ProductSrvConfig.Host,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%s启动在%d", randUUID, port))
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}

}
