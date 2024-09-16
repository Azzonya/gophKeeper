// Package adaptor provides a gRPC adaptor implementation for interacting with the GophKeeper server,
// offering methods to register, log in, and perform CRUD operations on data items.
package adaptor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "gophKeeper/pkg/proto/gophkeeper"
	"log/slog"
	"os"
	"sync"
	"time"
)

// GophKeeperClient represents the gRPC adaptor for interacting with the GophKeeper service.
// It handles both secure (TLS) and insecure connections and manages the Bearer token
// for authenticated requests.
type GophKeeperClient struct {
	client         pb.GophKeeperServiceClient
	wg             sync.WaitGroup
	enableTLS      bool
	serverAddress  string
	caFile         string
	clientCertFile string
	clientKeyFile  string

	ServerAvailable bool
	BearerToken     string
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
// adaptor certificate, and private key, and returns the configured TransportCredentials.
func loadTLSCredentials(caFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed adaptor's certificate
	pemServerCA, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	// Load adaptor's certificate and private key
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

func (c *GophKeeperClient) CreateContextWithMetadata(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	md := metadata.Pairs("token", "Bearer "+c.BearerToken)
	return metadata.NewOutgoingContext(ctx, md), cancel
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

// ListData sends a request to retrieve a data items from GophKeeper server.
func (c *GophKeeperClient) ListData(ctx context.Context, req *emptypb.Empty) (*pb.ListDataResponse, error) {
	return c.client.ListData(ctx, req)
}

// UpdateData sends a request to update an existing data item in the GophKeeper server.
func (c *GophKeeperClient) UpdateData(ctx context.Context, req *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	return c.client.UpdateData(ctx, req)
}

// DeleteData sends a request to delete a data item from the GophKeeper server.
func (c *GophKeeperClient) DeleteData(ctx context.Context, req *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	return c.client.DeleteData(ctx, req)
}

// SyncData sends a request to synchronize data between the adaptor and the GophKeeper server.
func (c *GophKeeperClient) SyncData(ctx context.Context, req *pb.SyncDataRequest) (*pb.SyncDataResponse, error) {
	return c.client.SyncData(ctx, req)
}

// IsServerAvailable pings server
func (c *GophKeeperClient) IsServerAvailable(ctx context.Context, req *emptypb.Empty, preStartHook bool) {
	if !preStartHook {
		defer c.wg.Done()
	}

	resp, err := c.client.Ping(ctx, req)
	c.ServerAvailable = !(err != nil || resp == nil)
}

// Start checks is server available for some period
func (c *GophKeeperClient) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.wg.Add(1)
				c.IsServerAvailable(ctx, &emptypb.Empty{}, false)
			}
		}
	}()
}

// Wait blocks until the WaitGroup counter is zero.
func (c *GophKeeperClient) Wait() {
	c.wg.Wait()
}
