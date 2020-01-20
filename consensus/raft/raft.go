package raft

import (
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"paysync/api"
	"paysync/lib/configs"
	"strings"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	libraft "github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

// JoinRaft ...
func JoinRaft(raftOptions *configs.RaftOptions) {
	cfg := &warden.ClientConfig{}

	c, err := api.NewClient(cfg)
	if err != nil {
		log.Error("new rpc client err: %v", err)
		return
	}
	ctx := context.Background()
	c.AddVoter(ctx, &api.VoterReq{Id: raftOptions.NodeID, Addr: raftOptions.Addr})
}

func newRaftTransport(addr string) (*libraft.NetworkTransport, error) {
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	transport, err := libraft.NewTCPTransport(address.String(), address, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}
	return transport, nil
}

// MysqlFSM ...
type MysqlFSM struct {
	logs [][]byte
}

// MockSnapshot ...
type MockSnapshot struct {
	logs [][]byte
}

// Persist saves the FSM snapshot out to the given sink.
func (s *MockSnapshot) Persist(sink libraft.SnapshotSink) error {
	log.Info("MockSnapshot Persist")

	// log.Info("MockSnapshot logs  len : %d, cap : %d\n", len(s.logs), cap(s.logs))
	var snapshotBytes strings.Builder
	for _, log := range s.logs {
		snapshotBytes.Write(log)
	}

	if _, err := sink.Write([]byte(snapshotBytes.String())); err != nil {
		log.Error("snapshot Write fail,err:%s\n", err)
		sink.Cancel()
		return err
	}
	if err := sink.Close(); err != nil {
		log.Error("snapshot Close fail,err:%s\n", err)
		sink.Cancel()
		return err
	}

	return nil
}

// Release ...
func (s *MockSnapshot) Release() {
	log.Info("MockSnapshot Release")
	s.logs = nil
}

// RaftNode ...
var RaftNode *libraft.Raft

// StartRaft ...
func StartRaft(raftOptions *configs.RaftOptions) *libraft.Raft {
	config := libraft.DefaultConfig()
	config.LocalID = libraft.ServerID(raftOptions.NodeID)

	if raftOptions.HeartbeatTimeout > 0 {
		config.HeartbeatTimeout = time.Duration(raftOptions.HeartbeatTimeout) * time.Second
	}

	if raftOptions.ElectionTimeout > 0 {
		config.ElectionTimeout = time.Duration(raftOptions.ElectionTimeout) * time.Second
	}

	if raftOptions.SnapshotInterval > 0 {
		config.SnapshotInterval = time.Duration(raftOptions.SnapshotInterval) * time.Second
	}

	if raftOptions.SnapshotThreshold > 0 {
		config.SnapshotThreshold = raftOptions.SnapshotThreshold
	}

	leadershipCh := make(chan bool)
	config.NotifyCh = leadershipCh

	go func() {
		for {
			select {
			case <-leadershipCh:
				log.Info("leadership change \n")
			}
		}
	}()

	snapshotStore, err := libraft.NewFileSnapshotStore(raftOptions.StoreDir, 1, os.Stderr)
	if err != nil {
		log.Error("new snapshotStore err: %s\n", err)
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftOptions.StoreDir,
		"raft-log.bolt"))
	if err != nil {
		log.Error("new logStore err: %s\n", err)
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftOptions.StoreDir,
		"raft-stable.bolt"))
	if err != nil {
		log.Error("new stableStore err: %s\n", err)
	}

	transport, err := newRaftTransport(raftOptions.Addr)
	if err != nil {
		log.Error("newRaftTransport err:%s\n", err.Error())
	}

	var configuration libraft.Configuration
	configuration.Servers = append(configuration.Servers, libraft.Server{
		Suffrage: libraft.Voter,
		ID:       libraft.ServerID(raftOptions.NodeID),
		Address:  transport.LocalAddr(),
	})

	fsm := &MysqlFSM{}

	// 初始化集群
	if raftOptions.Bootstrap {
		log.Info("will BootstrapCluster\n")
		libraft.BootstrapCluster(config, logStore, stableStore, snapshotStore, transport, configuration)
	}

	r, err := libraft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Error("new raft fail: %s\n", err.Error())
	}

	RaftNode = r
	return r
}

// Apply 日志条目commit后
func (fsm *MysqlFSM) Apply(thelog *libraft.Log) interface{} {
	// follower节点连接上会重做该方法
	log.Info("fsm apply,log is %s", thelog.Data)
	// execute sql

	fsm.logs = append(fsm.logs, thelog.Data)

	// log.Info("fsm.logs; len : %d, cap : %d", len(fsm.logs), cap(fsm.logs))
	return nil
}

// Snapshot ...
func (fsm *MysqlFSM) Snapshot() (libraft.FSMSnapshot, error) {
	log.Info("fsm Snapshot")
	// log.Info("fsm.logs; len : %d, cap : %d", len(fsm.logs), cap(fsm.logs))
	snapshot := &MockSnapshot{}

	snapshot.logs = fsm.logs

	fsm.logs = nil

	// log.Info("snapshot.logs; len : %d, cap : %d", len(snapshot.logs), cap(snapshot.logs))
	// log.Info("fsm.logs; len : %d, cap : %d", len(fsm.logs), cap(fsm.logs))

	return snapshot, nil
}

// Restore ...
func (fsm *MysqlFSM) Restore(serialized io.ReadCloser) error {
	log.Info("fsm Restore")

	// // var b []byte // A Buffer needs no initialization.
	// b := make([]byte, 20)
	// n, err := serialized.Read(b)
	// if err != nil {
	// 	log.Error("read buffer err:%s", err)
	// }
	// log.Info("read bytes len:%d,content:%s", n, string(b[:]))

	return nil
}
