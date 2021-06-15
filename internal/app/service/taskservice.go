package service

import (
	"context"
	pb "grpc-rusprofile-task/proto"
)

type TaskService interface {
	CompanyByInn(ctx context.Context, request *pb.CompanyByINNRequest) (*pb.CompanyByINNResponse, error)
}
