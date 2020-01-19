package raft

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"paysync/lib/configs"
	"sync"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	libraft "github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

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

// MockFSM ...
type MockFSM struct {
	sync.Mutex
	logs           [][]byte
	configurations []libraft.Configuration
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

	fsm := &MockFSM{}

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

// Apply ...
func (fsm *MockFSM) Apply(thelog *libraft.Log) interface{} {
	log.Info("fsm apply,log is %s\n", thelog.Data)
	return nil
}

// Snapshot ...
func (fsm *MockFSM) Snapshot() (libraft.FSMSnapshot, error) {
	log.Info("fsm Snapshot\n")
	return nil, nil
}

// Restore ...
func (fsm *MockFSM) Restore(io.ReadCloser) error {
	log.Info("fsm Restore\n")
	return nil
}
