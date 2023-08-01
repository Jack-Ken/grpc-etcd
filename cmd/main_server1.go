package main

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	demosrv "test/client"
	"test/etcd"
	"test/idl"
)

func main() {
	serverInfo := &etcd.Server{
		Name:    "user-service-1",
		Addr:    "127.0.0.1:8081",
		Version: "v1",
		Weight:  2,
	}

	listener, err := net.Listen("tcp", serverInfo.Addr)

	if err != nil {
		log.Fatal("net.Listen err: %v", err.Error())
	}

	// 创建注册器
	etcdRegister := etcd.NewRegister([]string{"127.0.0.1:2379"}, zap.L())
	defer etcdRegister.Close()
	_, err = etcdRegister.Register(serverInfo, 3)
	if err != nil {
		log.Fatal("register Server to etcd failed: %v", err.Error())
	}

	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	idl.RegisterDemoServiceServer(grpcServer, &demosrv.DemoService{})

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("grpcServer.Serve err: %v", err.Error())
	}
}
