package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"paysync/api"
	"paysync/connect"
	"paysync/consensus/raft"
	"paysync/internal/di"
	"paysync/lib/configs"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
)

var (
	join = flag.String("join", "", "join raft peer address")
)

func main() {
	flag.Parse()
	log.Init(nil) // debug flag: log.dir={path}
	defer log.Close()
	log.Info("paysync start")
	paladin.Init()

	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}

	appConfig := configs.Application()
	api.SetClientTarget(appConfig)
	connect.InitConnects(appConfig)

	raft.StartRaft(appConfig.RaftOptions)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info("paysync exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
