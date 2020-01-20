package connect

import (
	"context"
	"paysync/api"
	"paysync/lib/configs"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
)

type rpcClient struct {
	c       api.PaySyncClient
	sendMsg chan string
}

type clientsManage struct {
	clients []*rpcClient
}

// ClientsManage ...
var ClientsManage clientsManage = clientsManage{}

func (cm clientsManage) Sync(msg string) {
	for _, wrapRPCClient := range ClientsManage.clients {
		wrapRPCClient.sendMsg <- msg
	}
}

// InitConnects ...
func InitConnects(appConfig configs.ApplicationConfig) {

	for _, peer := range appConfig.Peers {
		log.Info("peer:%v\n", peer)
		cfg := &warden.ClientConfig{}
		c, err := api.NewClient(cfg)
		if err != nil {
			log.Error("new rpc client err: %v\n", err)
			continue
		}

		sc := make(chan string, 100)
		wrapRPCClient := &rpcClient{c, sc}

		ClientsManage.clients = append(ClientsManage.clients, wrapRPCClient)

		go manage(wrapRPCClient)
	}

}

func manage(wrapRPCClient *rpcClient) {

	done := make(chan struct{})
	c := wrapRPCClient.c

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// log.Info("call ping\n")
			// ctx := context.Background()
			// _, err := c.Ping(ctx, &empty.Empty{})
			// if err != nil {
			// 	log.Info("write:%v", err)
			// 	return
			// }
		case name := <-wrapRPCClient.sendMsg:
			log.Info("call SayHello\n")
			ctx := context.Background()
			_, err := c.SayHello(ctx, &api.HelloReq{Name: name})
			if err != nil {
				log.Info("write:%v", err)
				return
			}
		}
	}
}
