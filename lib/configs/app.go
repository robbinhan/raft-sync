package configs

import (
	"fmt"
	"github.com/bilibili/kratos/pkg/conf/paladin"
)

// RaftPeer ...
type RaftPeer struct {
	NodeID string
	Addr   string
}

// RaftOptions ...
type RaftOptions struct {
	NodeID            string
	Addr              string
	StoreDir          string
	Bootstrap         bool
	HeartbeatTimeout  int64
	ElectionTimeout   int64
	SnapshotInterval  int64
	SnapshotThreshold uint64
	Peers             []RaftPeer
}

// ApplicationConfig ...
type ApplicationConfig struct {
	Peers       []string
	RaftOptions *RaftOptions `toml:"raft"`
}

// Application ...
func Application() ApplicationConfig {

	app := ApplicationConfig{}

	fmt.Printf("app toml:%v\n", paladin.Get("application.toml"))

	if err := paladin.Get("application.toml").UnmarshalTOML(&app); err != nil {
		// 不存在时，将会为nil使用默认配置
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}

	fmt.Printf("app struct:%v\n", *app.RaftOptions)

	return app
}
