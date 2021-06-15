package app

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-rusprofile-task/configs"
	"grpc-rusprofile-task/internal/app/handler"
	service2 "grpc-rusprofile-task/internal/app/service"
	pb "grpc-rusprofile-task/proto"
	"log"
	"net"
	"net/http"
	"time"
)

func StartServerGRPS(ctx context.Context, cfg *configs.Config, errCh chan error) {
	l, lErr := net.Listen("tcp", cfg.GRPC.Addr)
	if lErr != nil {
		errCh <- lErr
		return
	}
	srv := grpc.NewServer()
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    15 * time.Second,
		DisableCompression: true,
	}
	cli := &http.Client{Transport: tr}
	service := service2.NewTaskService(cli, cfg)
	pb.RegisterTaskServer(srv, handler.NewHandler(service))
	reflection.Register(srv)
	log.Printf("Starting GRPC server on: %s\n", cfg.GRPC.Addr)
	if srvErr := srv.Serve(l); srvErr != nil {
		errCh <- srvErr
		return
	}
}

func StartServerHTTP(ctx context.Context, cfg *configs.Config, errCh chan error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(header string) (string, bool) {
			return header, true
		}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
	)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50000000)),
	}
	pb.RegisterTaskHandlerFromEndpoint(ctx, mux, cfg.GRPC.Addr, opts)
	log.Printf("Starting HTTP server on: %s\n", cfg.HTTP.Addr)
	if err := http.ListenAndServe(cfg.HTTP.Addr, wsproxy.WebsocketProxy(mux)); err != nil {
		errCh <- err
		return
	}
}

