package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumberjaye/shtrafov_net/internal/pb"
	"github.com/cucumberjaye/shtrafov_net/pkg/parser"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type searcherServer struct {
	pb.UnimplementedSearcherServer
}

func NewSearcherServer() pb.SearcherServer {
	return &searcherServer{}
}

func (s *searcherServer) GetProfile(ctx context.Context, pf *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			profile, err := parser.GetCompanyInfo(pf.GetInn())
			if err != nil {
				if errors.Is(err, parser.ErrNotFound) {
					return nil, status.Error(codes.NotFound, err.Error())
				}
				return nil, status.Error(codes.Internal, err.Error())
			}

			fmt.Println(profile)

			return profile, nil
		}
	}
}
