package service

import (
	"context"
	"fmt"
	"time"

	pb "paysync/api"
	"paysync/connect"
	"paysync/consensus/raft"
	"paysync/internal/dao"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
	libraft "github.com/hashicorp/raft"
)

// Provider ...
var Provider = wire.NewSet(New, wire.Bind(new(pb.PaySyncServer), new(*Service)))

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

// SayHello grpc demo func.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	fmt.Printf("hello %s\n", req.Name)
	connect.ClientsManage.Sync(req.Name)
	return
}

// AddVoter ...
func (s *Service) AddVoter(ctx context.Context, req *pb.VoterReq) (reply *pb.CommonResp, err error) {
	reply = new(pb.CommonResp)
	reply.Code = 0
	reply.Msg = "ok"
	future := raft.RaftNode.AddVoter(libraft.ServerID(req.Id), libraft.ServerAddress(req.Addr), 0, 10*time.Second)
	if future.Error() != nil {
		log.Error("AddVoter fail: %s\n", future.Error().Error())
		reply.Msg = future.Error().Error()
	}
	return
}

// AddData ...
func (s *Service) AddData(ctx context.Context, req *pb.DataReq) (reply *pb.CommonResp, err error) {
	reply = new(pb.CommonResp)
	reply.Code = 0
	reply.Msg = "ok"

	verifyFuture := raft.RaftNode.VerifyLeader()
	if verifyFuture.Error() != nil {
		log.Error("verifyFuture err:%s\n", verifyFuture)
		reply.Msg = "this node is not leader"
		return
	}

	l := libraft.Log{}
	l.Data = []byte(req.Data)
	future := raft.RaftNode.ApplyLog(l, 10*time.Second)
	if future.Error() != nil {
		log.Error("AddData fail: %s\n", future.Error().Error())
		reply.Msg = future.Error().Error()
	}
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	log.Info("recv Ping")
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}
