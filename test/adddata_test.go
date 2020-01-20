package test

import (
	"context"
	"paysync/api"
	"paysync/lib/configs"
	"testing"

	"github.com/bilibili/kratos/pkg/net/rpc/warden"
)

func TestAddData(t *testing.T) {
	cfg := &warden.ClientConfig{}

	appConf := configs.ApplicationConfig{
		Peers: []string{"0.0.0.0:9001"},
	}
	api.SetClientTarget(appConf)
	c, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("new rpc client err: %v\n", err)
	}
	ctx := context.Background()
	c.AddData(ctx, &api.DataReq{Data: "insert into t_qudao values(1,2,3,4);", Cmd: "sql"})
}
