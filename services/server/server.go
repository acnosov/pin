package server

import (
	"context"
	"database/sql"
	pb "github.com/aibotsoft/gen/surebetpb"
	"github.com/aibotsoft/micro/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type Server struct {
	cfg *config.Config
	log *zap.SugaredLogger
	//store *Store
	gs *grpc.Server
	pb.UnimplementedSurebetServer
}

func (s *Server) CheckLine(ctx context.Context, req *pb.CheckLineRequest) (*pb.CheckLineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckLine not implemented")
}

func NewServer(cfg *config.Config, log *zap.SugaredLogger, db *sql.DB) *Server {
	return &Server{
		cfg: cfg,
		log: log,
		//store: NewStore(cfg, log, db),
		gs: grpc.NewServer(),
	}
}
func (s *Server) Serve() error {
	addr := net.JoinHostPort("", "50051")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "net.Listen error")
	}
	pb.RegisterSurebetServer(s.gs, s)
	s.log.Info("gRPC Proxy Server listens on addr ", addr)
	return s.gs.Serve(lis)
}
func (s *Server) GracefulStop() {
	s.log.Debug("begin gRPC server gracefulStop")
	s.gs.GracefulStop()
	s.log.Debug("end gRPC server gracefulStop")
}
