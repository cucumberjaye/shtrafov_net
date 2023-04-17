package grpc

import (
	"fmt"
	"net"

	"github.com/cucumberjaye/shtrafov_net/internal/pb"
	"github.com/cucumberjaye/shtrafov_net/internal/service"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCServer() {
	var srv *grpc.Server
	var lis net.Listener

	srv = grpc.NewServer()
	searchServer := service.NewSearcherServer()
	pb.RegisterSearcherServer(srv, searchServer)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal().Err(fmt.Errorf("create listener failed with error: %w", err)).Send()
		return
	}

	log.Info().Msg("server startig")
	err = srv.Serve(lis)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("serve failed with error: %w", err)).Send()
	}
}
