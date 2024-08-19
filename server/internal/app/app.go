// Package app provides functionality to initialize, start, and stop the gophKeeper application.
// It sets up the GRPC server, repository, and database connections based on the configuration.
package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"gophKeeper/pkg/proto/gophkeeper"
	"gophKeeper/server/internal/conf"
	authorizerServiceP "gophKeeper/server/internal/domain/auth/service"
	dataItemsServiceP "gophKeeper/server/internal/domain/data_items/service"
	usersServiceP "gophKeeper/server/internal/domain/users/service"
	grpcHandler "gophKeeper/server/internal/handler/grpc"
	dataItemsUsecaseP "gophKeeper/server/internal/usecase/data_items"
	usersUsecaseP "gophKeeper/server/internal/usecase/users"
	"net"
	"os/signal"

	dataItemsRepoPgP "gophKeeper/server/internal/domain/data_items/repo/pg"
	dataItemsRepoS3P "gophKeeper/server/internal/domain/data_items/repo/s3"
	usersRepoPgP "gophKeeper/server/internal/domain/users/repo/pg"
	"log/slog"
	"os"
)

// App represent the application state containing configuration, GRPC server, database connection, and repository.
type App struct {
	pgpool *pgxpool.Pool

	// auth
	authorizer *authorizerServiceP.Auth

	// users
	usersUsecase *usersUsecaseP.Usecase

	// data itesms
	dataItemsUsecase *dataItemsUsecaseP.Usecase

	// grpc server
	grpcServer *grpc.Server

	exitCode int
}

// Init initializes the application with the provided configuration
func (a *App) Init() {
	var err error

	// pgpool
	{
		a.pgpool, err = pgxpool.New(context.Background(), conf.Conf.PgDsn)
		errCheck(err, "pgxpool.New")
	}

	// auth
	{
		a.authorizer = authorizerServiceP.New(conf.Conf.JwtSecret)
	}

	// users
	{
		usersRepo := usersRepoPgP.New(a.pgpool)
		usersService := usersServiceP.New(usersRepo)
		a.usersUsecase = usersUsecaseP.New(usersService, a.authorizer)
	}

	// data items
	{
		dataItemsPgRepo := dataItemsRepoPgP.New(a.pgpool)
		dataItemsS3Repo, err := dataItemsRepoS3P.NewS3Repo(context.Background(), conf.Conf.S3Endpoint, conf.Conf.S3AccessKey, conf.Conf.S3SecretKey, conf.Conf.S3Bucket)
		errCheck(err, "dataItemsS3Repo")
		dataItemsSerivce := dataItemsServiceP.New(dataItemsPgRepo, dataItemsS3Repo)
		a.dataItemsUsecase = dataItemsUsecaseP.New(dataItemsSerivce)
	}

	// grpc server
	{
		var opts []grpc.ServerOption
		if conf.Conf.EnableTLS {
			tlsConfig, err := loadTLSCredentials()
			if err != nil {
				errCheck(err, "tls.LoadTLSCredentials")
			}

			opts = append(opts, grpc.Creds(tlsConfig))
		}

		interceptors := make([]grpc.UnaryServerInterceptor, 0, 3)

		interceptors = append(interceptors, grpcHandler.GrpcInterceptorLogger())

		opts = append(opts, grpc.ChainUnaryInterceptor(interceptors...))

		a.grpcServer = grpc.NewServer(opts...)

		grpcHandlers := grpcHandler.New(a.dataItemsUsecase, a.usersUsecase)
		gophkeeper.RegisterGophKeeperServiceServer(a.grpcServer, grpcHandlers)

		reflection.Register(a.grpcServer)
	}
}

// Start starts the application, initializing and running GRPC server.
func (a *App) Start() {
	slog.Info("Starting")

	// grpc server
	{
		lis, err := net.Listen("tcp", conf.Conf.GRPCPort)
		if err != nil {
			errCheck(err, "net.Listen")
		}
		go func() {
			err = a.grpcServer.Serve(lis)
			if err != nil {
				errCheck(err, "grpcServer.Serve")
			}
		}()

		slog.Info("GRPC-server started successfully " + lis.Addr().String())
	}
}

// Listen listens for signals to stop the application
func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer signalCtxCancel()

	// wait signal
	<-signalCtx.Done()
}

// Stop stops the application, shutting down the GRPC server.
func (a *App) Stop() {
	slog.Info("Shutting down...")

	// grpc server
	{
		a.grpcServer.GracefulStop()
	}
}

// Exit gracefully shuts down the application by logging the exit action
// and then terminating the program with the specified exit code.
func (a *App) Exit() {
	slog.Info("Exit")

	os.Exit(a.exitCode)
}

// loadTLSCredentials loads the TLS credentials for the server, including the server's
// certificate, private key, and the client's Certificate Authority (CA) certificate for mutual TLS.
// It returns the configured TransportCredentials and an error if any of the loading steps fail.
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := os.ReadFile(conf.Conf.CAFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to append client CA certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(conf.Conf.ServerCertFile, conf.Conf.ServerKeyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

// errCheck checks if an error occurred and logs it with the specified message.
// If an error is found, the function logs the error and terminates the program.
// If a message is provided, it is included in the logged output.
func errCheck(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}
