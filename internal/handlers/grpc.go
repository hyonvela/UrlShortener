package handlers

import (
	"context"

	"example.com/m/internal/usecase"
	grpcV1 "example.com/m/pkg/grpc.v1"
	"example.com/m/pkg/logging"
)

type Server struct {
	uc  usecase.Usecase
	log *logging.Logger
	grpcV1.UnimplementedUrlShortenerServiceServer
}

func NewGRPCHandler(uc usecase.Usecase, log *logging.Logger) *Server {
	return &Server{uc: uc, log: log}
}

func (s *Server) ShortenUrl(ctx context.Context, req *grpcV1.UrlRequest) (*grpcV1.UrlResponse, error) {
	short, err := s.uc.ShortenUrl(req.Url, ctx)
	if err != nil {
		return nil, err
	}
	return &grpcV1.UrlResponse{Url: short}, nil
}

func (s *Server) GetLongUrl(ctx context.Context, req *grpcV1.UrlRequest) (*grpcV1.UrlResponse, error) {
	long, err := s.uc.GetLongUrl(req.Url, ctx)
	if err != nil {
		return nil, err
	}
	return &grpcV1.UrlResponse{Url: long}, nil
}
