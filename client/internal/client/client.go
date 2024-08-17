// Package client provides a gRPC client implementation for interacting with the GophKeeper server,
// offering methods to register, log in, and perform CRUD operations on data items.
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	pb "gophKeeper/pkg/proto/gophkeeper"
	"log/slog"
	"os"
)

// GophKeeperClient represents the gRPC client for interacting with the GophKeeper service.
// It handles both secure (TLS) and insecure connections and manages the Bearer token
// for authenticated requests.
type GophKeeperClient struct {
	client         pb.GophKeeperServiceClient
	enableTLS      bool
	serverAddress  string
	caFile         string
	clientCertFile string
	clientKeyFile  string

	BearerToken string
}

// NewGophKeeperClient creates a new GophKeeperClient instance, setting up the gRPC connection
// with either secure (TLS) or insecure credentials based on the provided configuration.
func NewGophKeeperClient(enableTLS bool, serverAddress, caFile, clientCertFile, clientKeyFile string) (*GophKeeperClient, error) {
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if enableTLS {
		tlsConfig, err := loadTLSCredentials(caFile, clientCertFile, clientKeyFile)
		if err != nil {
			return nil, err
		}

		transportOption = grpc.WithTransportCredentials(tlsConfig)
	}

	conn, err := grpc.NewClient(serverAddress, transportOption)
	if err != nil {
		slog.Error("NewGophKeeperClient error", slog.String("error", err.Error()))

	}

	return &GophKeeperClient{
		client:        pb.NewGophKeeperServiceClient(conn),
		serverAddress: serverAddress,
	}, nil
}

// loadTLSCredentials loads the necessary TLS credentials, including the CA certificate,
// client certificate, and private key, and returns the configured TransportCredentials.
func loadTLSCredentials(caFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemServerCA, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

// Register sends a registration request to the GophKeeper server.
func (c *GophKeeperClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return c.client.Register(ctx, req)
}

// Login sends a login request to the GophKeeper server and returns a response containing a token.
func (c *GophKeeperClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return c.client.Login(ctx, req)
}

// CreateData sends a request to create a new data item in the GophKeeper server.
func (c *GophKeeperClient) CreateData(ctx context.Context, req *pb.CreateDataRequest) (*pb.CreateDataResponse, error) {
	return c.client.CreateData(ctx, req)
}

// GetData sends a request to retrieve a data item from the GophKeeper server.
func (c *GophKeeperClient) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	return c.client.GetData(ctx, req)
}

// UpdateData sends a request to update an existing data item in the GophKeeper server.
func (c *GophKeeperClient) UpdateData(ctx context.Context, req *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	return c.client.UpdateData(ctx, req)
}

// DeleteData sends a request to delete a data item from the GophKeeper server.
func (c *GophKeeperClient) DeleteData(ctx context.Context, req *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	return c.client.DeleteData(ctx, req)
}

// SyncData sends a request to synchronize data between the client and the GophKeeper server.
func (c *GophKeeperClient) SyncData(ctx context.Context, req *pb.SyncDataRequest) (*pb.SyncDataResponse, error) {
	return c.client.SyncData(ctx, req)
}
