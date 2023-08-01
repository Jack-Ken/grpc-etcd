package main

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"test/etcd"
	"test/idl"
	"time"
)

func main() {
	// 注册自定义的ETCD解析器
	etcdResolverBuilder := etcd.NewResolver([]string{"127.0.0.1:2379"}, zap.L())
	resolver.Register(etcdResolverBuilder)
	// 连接服务器
	conn, err := grpc.Dial(etcdResolverBuilder.Scheme()+":///user-service-1/v1", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(
		`{"loadBalancingPolicy":"weight_lb_picker"}`,
	))
	if err != nil {
		log.Fatal("net.Connect err: %v", err.Error())
	}
	defer conn.Close()

	// 建立gRPC连接
	grpcClient := idl.NewDemoServiceClient(conn)
	// 创建发送结构体
	req := idl.HiRequest{
		Name: "Angry Potato",
	}
	for {
		time.Sleep(2 * time.Second)
		grpcClient.SayHi(context.Background(), &req)
		//resp, _ := grpcClient.SayHi(context.Background(), &req)
		//fmt.Println(resp.Message)
	}
}
