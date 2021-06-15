package handler

import (
	"context"
	"grpc-rusprofile-task/internal/app/service"
	pb "grpc-rusprofile-task/proto"
)

type Handler struct {
	taskService service.TaskService
}

func NewHandler(taskService service.TaskService) *Handler {
	return &Handler{taskService: taskService}
}


func (h Handler) CompanyByINN(ctx context.Context, request *pb.CompanyByINNRequest) (*pb.CompanyByINNResponse, error) {
	return h.taskService.CompanyByInn(ctx, request)

}

//func responseHandler(ctx context.Context, resp interface{}, statusCode int) (*pb.CompanyByINNResponse, error) {
//	md, ok := metadata.FromIncomingContext(ctx)
//	if !ok {
//		fmt.Println("asdasdasd")
//	}
//
//	md.
//	_, isNotFound := resp.(pb.ErrorNotFound)
//	if isNotFound{
//		grpc.SendHeader()
//	}
//}



