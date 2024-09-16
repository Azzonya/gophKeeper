package adaptor

import (
	"context"
	"google.golang.org/grpc/credentials"
	"gophKeeper/pkg/proto/gophkeeper"
	"reflect"
	"testing"
)

func TestGophKeeperClient_CreateData(t *testing.T) {
	type fields struct {
		client         gophkeeper.GophKeeperServiceClient
		enableTLS      bool
		serverAddress  string
		caFile         string
		clientCertFile string
		clientKeyFile  string
		BearerToken    string
	}
	type args struct {
		ctx context.Context
		req *gophkeeper.CreateDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gophkeeper.CreateDataResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GophKeeperClient{
				client:         tt.fields.client,
				enableTLS:      tt.fields.enableTLS,
				serverAddress:  tt.fields.serverAddress,
				caFile:         tt.fields.caFile,
				clientCertFile: tt.fields.clientCertFile,
				clientKeyFile:  tt.fields.clientKeyFile,
				BearerToken:    tt.fields.BearerToken,
			}
			got, err := c.CreateData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGophKeeperClient_DeleteData(t *testing.T) {
	type fields struct {
		client         gophkeeper.GophKeeperServiceClient
		enableTLS      bool
		serverAddress  string
		caFile         string
		clientCertFile string
		clientKeyFile  string
		BearerToken    string
	}
	type args struct {
		ctx context.Context
		req *gophkeeper.DeleteDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gophkeeper.DeleteDataResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GophKeeperClient{
				client:         tt.fields.client,
				enableTLS:      tt.fields.enableTLS,
				serverAddress:  tt.fields.serverAddress,
				caFile:         tt.fields.caFile,
				clientCertFile: tt.fields.clientCertFile,
				clientKeyFile:  tt.fields.clientKeyFile,
				BearerToken:    tt.fields.BearerToken,
			}
			got, err := c.DeleteData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGophKeeperClient_GetData(t *testing.T) {
	type fields struct {
		client         gophkeeper.GophKeeperServiceClient
		enableTLS      bool
		serverAddress  string
		caFile         string
		clientCertFile string
		clientKeyFile  string
		BearerToken    string
	}
	type args struct {
		ctx context.Context
		req *gophkeeper.GetDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gophkeeper.GetDataResponse
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GophKeeperClient{
				client:         tt.fields.client,
				enableTLS:      tt.fields.enableTLS,
				serverAddress:  tt.fields.serverAddress,
				caFile:         tt.fields.caFile,
				clientCertFile: tt.fields.clientCertFile,
				clientKeyFile:  tt.fields.clientKeyFile,
				BearerToken:    tt.fields.BearerToken,
			}
			got, err := c.GetData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGophKeeperClient_Login(t *testing.T) {
	type fields struct {
		client         gophkeeper.GophKeeperServiceClient
		enableTLS      bool
		serverAddress  string
		caFile         string
		clientCertFile string
		clientKeyFile  string
		BearerToken    string
	}
	type args struct {
		ctx context.Context
		req *gophkeeper.LoginRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gophkeeper.LoginResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GophKeeperClient{
				client:         tt.fields.client,
				enableTLS:      tt.fields.enableTLS,
				serverAddress:  tt.fields.serverAddress,
				caFile:         tt.fields.caFile,
				clientCertFile: tt.fields.clientCertFile,
				clientKeyFile:  tt.fields.clientKeyFile,
				BearerToken:    tt.fields.BearerToken,
			}
			got, err := c.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGophKeeperClient_Register(t *testing.T) {
	type fields struct {
		client         gophkeeper.GophKeeperServiceClient
		enableTLS      bool
		serverAddress  string
		caFile         string
		clientCertFile string
		clientKeyFile  string
		BearerToken    string
	}
	type args struct {
		ctx context.Context
		req *gophkeeper.RegisterRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gophkeeper.RegisterResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GophKeeperClient{
				client:         tt.fields.client,
				enableTLS:      tt.fields.enableTLS,
				serverAddress:  tt.fields.serverAddress,
				caFile:         tt.fields.caFile,
				clientCertFile: tt.fields.clientCertFile,
				clientKeyFile:  tt.fields.clientKeyFile,
				BearerToken:    tt.fields.BearerToken,
			}
			got, err := c.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGophKeeperClient_SyncData(t *testing.T) {
	type fields struct {
		client         gophkeeper.GophKeeperServiceClient
		enableTLS      bool
		serverAddress  string
		caFile         string
		clientCertFile string
		clientKeyFile  string
		BearerToken    string
	}
	type args struct {
		ctx context.Context
		req *gophkeeper.SyncDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gophkeeper.SyncDataResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GophKeeperClient{
				client:         tt.fields.client,
				enableTLS:      tt.fields.enableTLS,
				serverAddress:  tt.fields.serverAddress,
				caFile:         tt.fields.caFile,
				clientCertFile: tt.fields.clientCertFile,
				clientKeyFile:  tt.fields.clientKeyFile,
				BearerToken:    tt.fields.BearerToken,
			}
			got, err := c.SyncData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGophKeeperClient_UpdateData(t *testing.T) {
	type fields struct {
		client         gophkeeper.GophKeeperServiceClient
		enableTLS      bool
		serverAddress  string
		caFile         string
		clientCertFile string
		clientKeyFile  string
		BearerToken    string
	}
	type args struct {
		ctx context.Context
		req *gophkeeper.UpdateDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gophkeeper.UpdateDataResponse
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GophKeeperClient{
				client:         tt.fields.client,
				enableTLS:      tt.fields.enableTLS,
				serverAddress:  tt.fields.serverAddress,
				caFile:         tt.fields.caFile,
				clientCertFile: tt.fields.clientCertFile,
				clientKeyFile:  tt.fields.clientKeyFile,
				BearerToken:    tt.fields.BearerToken,
			}
			got, err := c.UpdateData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGophKeeperClient(t *testing.T) {
	type args struct {
		enableTLS      bool
		serverAddress  string
		caFile         string
		clientCertFile string
		clientKeyFile  string
	}
	tests := []struct {
		name    string
		args    args
		want    *GophKeeperClient
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGophKeeperClient(tt.args.enableTLS, tt.args.serverAddress, tt.args.caFile, tt.args.clientCertFile, tt.args.clientKeyFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGophKeeperClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGophKeeperClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadTLSCredentials(t *testing.T) {
	type args struct {
		caFile         string
		clientCertFile string
		clientKeyFile  string
	}
	tests := []struct {
		name    string
		args    args
		want    credentials.TransportCredentials
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadTLSCredentials(tt.args.caFile, tt.args.clientCertFile, tt.args.clientKeyFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadTLSCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadTLSCredentials() got = %v, want %v", got, tt.want)
			}
		})
	}
}
