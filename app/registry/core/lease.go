package core

import (
	"time"

	v1 "zflow/api/registry"
)

// lease 生成租约
func lease(in *v1.ServiceInstance) *v1.Lease {
	return &v1.Lease{
		Name:       in.Name,
		Id:         in.Id,
		ExpireUnix: time.Now().Add(time.Duration(in.TtlSec) * time.Second).Unix(),
	}
}
