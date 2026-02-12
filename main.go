package main

import (
	"context"
	"fmt"
	logger "go-base/pkg"
	gozero "go-base/pkg/go-zero"
	demov1 "go-base/pkg/proto/demo/v1"
)

func main() {
	// 初始化 logger
	if err := logger.LoggerInit("./logs"); err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		return
	}

	gozero.Init("demo-rpc", gozero.WithEndpoints("192.168.8.245:9091"))

	demoSvc := demov1.NewDemoServiceClient(gozero.MustConn("demo-rpc"))
	ctx := context.Background()
	resp, err := demoSvc.CreateDemo(ctx, &demov1.CreateDemoRequest{Name: "test", Description: "test", Status: 1})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
