package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	v1 "zflow/api/registry"
	"zflow/app/registry/core"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	// 创建 gRPC 服务器
	lis, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// 注册 registry 服务
	v1.RegisterRegistryServer(s, core.NewRegistry())

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
