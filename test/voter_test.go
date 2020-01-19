package test

import (
	"context"
	"paysync/api"
	"paysync/lib/configs"
	"testing"

	"github.com/bilibili/kratos/pkg/net/rpc/warden"
)

func TestAddVoter(t *testing.T) {
	cfg := &warden.ClientConfig{}

	appConf := configs.ApplicationConfig{
		Peers: []string{"0.0.0.0:9000"},
	}
	api.SetClientTarget(appConf)

	c, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("new rpc client err: %v\n", err)
	}
	ctx := context.Background()
	c.AddVoter(ctx, &api.VoterReq{Id: "N2", Addr: "127.0.0.1:7002"})
}
