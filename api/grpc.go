package api

import (
	"context"
	"url-shortener-golang/storage"

	pb "url-shortener-golang/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedURLShortenerServer
	storage storage.Storage
}

// NewGRPCServer создаёт новый экземпляр gRPC-сервера
func NewGRPCServer(s storage.Storage) *GRPCServer {
	return &GRPCServer{storage: s}
}

// Shorten создаёт сокращённую ссылку
func (s *GRPCServer) Shorten(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	shortURL, err := s.storage.Post(req.OriginalUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ShortenResponse{ShortUrl: shortURL}, nil
}

// Resolve возвращает оригинальную ссылку по сокращённой
func (s *GRPCServer) Resolve(ctx context.Context, req *pb.ResolveRequest) (*pb.ResolveResponse, error) {
	originalURL, err := s.storage.Get(req.ShortUrl)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ResolveResponse{OriginalUrl: originalURL}, nil
}
