package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gophKeeper/server/internal/domain/data_items/model"
	dataItemsRepoPgP "gophKeeper/server/internal/domain/data_items/repo/pg"
	dataItemsRepoS3P "gophKeeper/server/internal/domain/data_items/repo/s3"
	"log"
	"reflect"
	"testing"
	"time"
)

func getPgPoolTestContainer() (*pgxpool.Pool, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port())
	pgpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	log.Print(dsn)

	return pgpool, err
}

func minioContainerStart() (string, string, string, string) {
	ctx := context.Background()

	// Запрос на запуск контейнера с MinIO
	req := testcontainers.ContainerRequest{
		Image:        "minio/minio",
		ExposedPorts: []string{"9005/tcp"},
		Cmd:          []string{"server", "/data"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     "minioadmin",
			"MINIO_ROOT_PASSWORD": "minioadmin",
		},
		WaitingFor: wait.ForLog("API: http://0.0.0.0:9005").WithStartupTimeout(60 * time.Second), // Ожидание появления лога
	}

	minioContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Ошибка при создании контейнера: %v", err)
	}
	defer minioContainer.Terminate(ctx)

	host, err := minioContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Ошибка при получении хоста: %v", err)
	}

	port, err := minioContainer.MappedPort(ctx, "9005")
	if err != nil {
		log.Fatalf("Ошибка при маппинге порта: %v", err)
	}

	endpoint := fmt.Sprintf("%s:%s", host, port.Port())
	accessKey := "minioadmin"
	secretKey := "minioadmin"

	// Инициализация клиента MinIO
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Ошибка при создании MinIO клиента: %v", err)
	}

	bucketName := "my-bucket"
	location := "us-east-1"

	// Ожидание перед созданием bucket (можно убрать при использовании wait.ForLog)
	time.Sleep(5 * time.Second)

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("Bucket %s уже существует\n", bucketName)
		} else {
			log.Fatalf("Ошибка при создании bucket: %v", err)
		}
	} else {
		fmt.Printf("Успешно создан bucket %s\n", bucketName)
	}

	fmt.Printf("MinIO запущен на: %s\n", endpoint)
	fmt.Printf("AccessKey: %s\n", accessKey)
	fmt.Printf("SecretKey: %s\n", secretKey)
	fmt.Printf("Bucket: %s\n", bucketName)

	return endpoint, accessKey, secretKey, bucketName
}

func TestNew(t *testing.T) {
	pgpool, err := getPgPoolTestContainer()
	if err != nil {
		t.Fatal(err)
	}

	endpoint, accessKey, secretKey, bucketName := minioContainerStart()

	dataItemsPgRepo := dataItemsRepoPgP.New(pgpool)
	dataItemsS3Repo, err := dataItemsRepoS3P.NewS3Repo(context.Background(), endpoint, accessKey, secretKey, bucketName)

	type args struct {
		repoDB RepoDBI
		repoS3 RepoS3
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Create new data items service",
			args: args{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			want: &Service{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.repoDB, tt.args.repoS3); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	type fields struct {
		repoDB RepoDBI
		repoS3 RepoS3
	}
	type args struct {
		ctx context.Context
		obj *model.Edit
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: tt.fields.repoDB,
				repoS3: tt.fields.repoS3,
			}
			if err := s.Create(tt.args.ctx, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	type fields struct {
		repoDB RepoDBI
		repoS3 RepoS3
	}
	type args struct {
		ctx  context.Context
		pars *model.GetPars
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: tt.fields.repoDB,
				repoS3: tt.fields.repoS3,
			}
			if err := s.Delete(tt.args.ctx, tt.args.pars); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	type fields struct {
		repoDB RepoDBI
		repoS3 RepoS3
	}
	type args struct {
		ctx  context.Context
		pars *model.GetPars
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Main
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: tt.fields.repoDB,
				repoS3: tt.fields.repoS3,
			}
			got, got1, err := s.Get(tt.args.ctx, tt.args.pars)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	type fields struct {
		repoDB RepoDBI
		repoS3 RepoS3
	}
	type args struct {
		ctx  context.Context
		pars *model.ListPars
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Main
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: tt.fields.repoDB,
				repoS3: tt.fields.repoS3,
			}
			got, got1, err := s.List(tt.args.ctx, tt.args.pars)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	type fields struct {
		repoDB RepoDBI
		repoS3 RepoS3
	}
	type args struct {
		ctx  context.Context
		pars *model.GetPars
		obj  *model.Edit
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: tt.fields.repoDB,
				repoS3: tt.fields.repoS3,
			}
			if err := s.Update(tt.args.ctx, tt.args.pars, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
