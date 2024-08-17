// Package app manages the lifecycle of the GophKeeper client application,
// handling initialization, starting the TUI, listening for shutdown signals,
// and performing cleanup during shutdown.
package app

import (
	"context"
	"fmt"
	"gophKeeper/client/internal/client"
	"gophKeeper/client/internal/conf"
	"gophKeeper/client/internal/tui"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// App represents the main application structure, containing the gRPC client,
// TUI (text user interface), and methods for controlling the application's lifecycle.
type App struct {
	// grpc client
	grpcClient *client.GophKeeperClient
	TUI        *tui.TUI

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

	// TUI
	{
		a.TUI = tui.NewTUI(a.grpcClient)
	}

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
}

// Listen waits for shutdown signals like SIGTERM and handles graceful shutdown when received.
func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCtxCancel()

	// Wait signal
	slog.Info("Listening for shutdown signal...")
	<-signalCtx.Done()
	slog.Info("Shutdown signal received")
}

// Exit terminates the application with the specified exit code.
func (a *App) Exit() {
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
