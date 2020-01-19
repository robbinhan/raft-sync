package api

import (
	"context"
	"paysync/lib/configs"
	"strings"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"

	"google.golang.org/grpc"
)

// target server addrs.
var target = "direct://default/"

// SetClientTarget ...
func SetClientTarget(appConfig configs.ApplicationConfig) {
	var peers strings.Builder
	for _, peer := range appConfig.Peers {
		peers.WriteString(peer)
		peers.WriteString(",")
	}
	target = target + peers.String()
	log.Info("target is : %s\n", target)
}

// NewClient new grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (PaySyncClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), target)
	if err != nil {
		return nil, err
	}
	return NewPaySyncClient(cc), nil
}

// 生成 gRPC 代码
//go:generate kratos tool protoc --grpc --bm api.proto
