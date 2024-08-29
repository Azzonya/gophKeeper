// Package app manages the lifecycle of the GophKeeper client application,
// handling initialization, starting the TUI, listening for shutdown signals,
// and performing cleanup during shutdown.
package app

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/emptypb"
	"gophKeeper/client/internal/client"
	"gophKeeper/client/internal/conf"
	"gophKeeper/client/internal/tui"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// App represents the main application structure, containing the gRPC client,
// TUI (text user interface), and methods for controlling the application's lifecycle.
type App struct {
	// grpc client
	grpcClient *client.GophKeeperClient

	// TUI
	TUI *tui.TUI

	// cache
	redisClient *redis.Client

	exitCode int
}

// Init initializes the application by setting up the gRPC client and TUI components.
func (a *App) Init() {
	var err error
	//var err error

	// grpc client
	{
		a.grpcClient, err = client.NewGophKeeperClient(
			conf.Conf.EnableTLS,
			conf.Conf.ServerAddress,
			conf.Conf.CAFile,
			conf.Conf.ClientCertFile,
			conf.Conf.ClientKeyFile,
		)
		errCheck(err, "NewGophKeeperClient")
	}

	// redis
	{
		a.redisClient = redis.NewClient(&redis.Options{
			Addr:     conf.Conf.RedisAddress,
			Password: conf.Conf.RedisPassword,
			DB:       conf.Conf.RedisDb,
		})

		err = a.redisClient.Set(context.Background(), "key", "value", 0).Err()
		if err != nil {
			log.Fatalf("Could not set cache: %v", err)
		}
	}

	// TUI
	{
		a.TUI = tui.NewTUI(a.grpcClient, a.redisClient)
	}
}

func (a *App) PreStartHook() {
	slog.Info("PreStartHook")
	a.grpcClient.IsServerAvailable(context.Background(), &emptypb.Empty{}, true)
}

// Start runs the main application logic, including launching the TUI.
func (a *App) Start() {
	slog.Info("Starting")

	// TUI
	{
		go func() {
			if err := a.TUI.Run(); err != nil {
				errCheck(err, "TUI Run")
			}
		}()
	}

	// client
	{
		a.grpcClient.Start(context.Background(), time.Minute)
	}
}

// Listen waits for shutdown signals like SIGTERM and handles graceful shutdown when received.
func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCtxCancel()

	// Wait signal
	slog.Info("Listening for shutdown signal...")
	<-signalCtx.Done()
	a.grpcClient.Wait()
	slog.Info("Shutdown signal received")
}

// Exit terminates the application with the specified exit code.
func (a *App) Exit() {

	if err := a.redisClient.Close(); err != nil {
		slog.Error("Error closing Redis client")
	} else {
		slog.Info("Redis client closed")
	}

	slog.Info("Exit")

	os.Exit(a.exitCode)
}

// errCheck handles errors by logging them and exiting the application if an error occurs.
func errCheck(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}
