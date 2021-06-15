package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"grpc-rusprofile-task/configs"
	"grpc-rusprofile-task/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//go:embed configs.json
var fs embed.FS

const configName = "configs.json"

func main()  {
	data, readErr := fs.ReadFile(configName)
	if readErr != nil {
		log.Fatal(readErr)
	}
	cfg := configs.NewConfig()
	if unmErr := json.Unmarshal(data, cfg); unmErr != nil {
		log.Fatal(unmErr)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	go app.StartServerGRPS(ctx, cfg, errCh)
	go app.StartServerHTTP(ctx, cfg, errCh)
	go func(ctx context.Context, errCh chan error) {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%v", <-sigCh)
	}(ctx, errCh)
	log.Printf("\nService terminated: %v", <-errCh)
}
