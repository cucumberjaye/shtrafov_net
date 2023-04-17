package gw

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cucumberjaye/shtrafov_net/internal/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/rs/zerolog/log"
)

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger/swagger.json")
}

func Run() error {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("connection to server failed with error: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	rmux := runtime.NewServeMux()
	err = pb.RegisterSearcherHandler(ctx, rmux, conn)
	if err != nil {
		return fmt.Errorf("register http hanlers failed with error: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.HandleFunc("/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir("swagger/swagger-ui"))
	mux.Handle("/swaggerui/", http.StripPrefix("/swaggerui", fs))

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Info().Msg("client startig")
	return http.ListenAndServe(":8001", mux)
}
