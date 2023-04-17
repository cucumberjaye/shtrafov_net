package main

import (
	"github.com/cucumberjaye/shtrafov_net/pkg/grpc"
	"github.com/cucumberjaye/shtrafov_net/pkg/gw"
	"github.com/rs/zerolog/log"
)

func main() {
	go grpc.StartGRPCServer()
	err := gw.Run()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
