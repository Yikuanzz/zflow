package main

import (
	"zflow/app/zflow/model"
	"zflow/utils/micro"
)

var (
	serviceName = "service_zz"
	serviceAddr = "127.0.0.1:50051"
)

var (
	nodeTypes = map[string]*model.NodeType{
		"node_type_1": {
			UID: "node_type_1",
		},
	}
)

var (
	connTypes = map[string]*model.ConnectionType{
		"conn_type_1": {
			UID: "conn_type_1",
		},
	}
)

func main() {
	micro := micro.NewMicro(serviceName, serviceAddr, nodeTypes, connTypes)
	micro.Run()
}
