package service

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Импортируем драйвер для работы с файлами
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gophKeeper/server/internal/domain/dataitems/model"
	dataItemsRepoPgP "gophKeeper/server/internal/domain/dataitems/repo/pg"
	dataItemsRepoS3P "gophKeeper/server/internal/domain/dataitems/repo/s3"
	"log"
	"reflect"
	"testing"
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

	err = migrateUp(dsn)
	if err != nil {
		return nil, err
	}

	conn, err := pgpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "INSERT INTO users (id, username, password_hash) VALUES ('999', 'Test User', 'test')")
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных: %w", err)
	}

	return pgpool, err
}

func migrateUp(dsn string) error {
	m, err := migrate.New(
		"file://../../../../migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("Ошибка при настройке миграций: %v\n", err)
	}

	// Применение всех миграций
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка при применении миграций: %v\n", err)
	}

	return err
}

func setupMinio(ctx context.Context) (string, string, string, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     "minioadmin",
			"MINIO_ROOT_PASSWORD": "minioadmin",
		},
		WaitingFor: wait.ForHTTP("/minio/health/live").WithPort("9000/tcp"),
		Cmd:        []string{"server", "/data"},
	}

	minioContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", "", "", "", fmt.Errorf("не удалось запустить контейнер: %w", err)
	}

	// Получаем порт, по которому доступен MinIO
	hostPort, err := minioContainer.MappedPort(ctx, "9000")
	if err != nil {
		return "", "", "", "", fmt.Errorf("не удалось получить порт: %w", err)
	}

	endpoint := fmt.Sprintf("localhost:%s", hostPort.Port())

	// Конфигурация доступа
	accessKey := "minioadmin"
	secretKey := "minioadmin"
	bucketName := "test-bucket"

	return endpoint, accessKey, secretKey, bucketName, nil
}

func TestNew(t *testing.T) {
	pgpool, err := getPgPoolTestContainer()
	if err != nil {
		t.Fatal(err)
	}

	endpoint, accessKey, secretKey, bucketName, err := setupMinio(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	dataItemsPgRepo := dataItemsRepoPgP.New(pgpool)
	dataItemsS3Repo, err := dataItemsRepoS3P.NewS3Repo(context.Background(), endpoint, accessKey, secretKey, bucketName)
	if err != nil {
		t.Fatal(err)
	}

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
	dataItemsPgRepo, dataItemsS3Repo, err := testRepos()
	if err != nil {
		t.Fatal(err)
	}

	testUserID := "999"
	bankCardType := model.CredentialsDataType
	credentialsType := model.CredentialsDataType
	binaryType := model.BinaryDataType
	data := []byte("test/test")
	meta := "binary"

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
		{
			name: "Create new data items service - bank card",
			fields: fields{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Edit{
					ID:     uuid.New().String(),
					UserID: &testUserID,
					Type:   &bankCardType,
					Data:   &data,
				},
			},
			wantErr: false,
		},
		{
			name: "Create new data items service - credentials",
			fields: fields{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Edit{
					ID:     uuid.New().String(),
					UserID: &testUserID,
					Type:   &credentialsType,
					Data:   &data,
				},
			},
			wantErr: false,
		},
		{
			name: "Create new data items service - binary",
			fields: fields{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Edit{
					ID:     uuid.New().String(),
					UserID: &testUserID,
					Type:   &binaryType,
					Data:   &data,
					Meta:   &meta,
				},
			},
			wantErr: false,
		},
		{
			name: "Create new data items service - binary error",
			fields: fields{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Edit{
					ID:     uuid.New().String(),
					UserID: &testUserID,
					Type:   &binaryType,
					Data:   &data,
				},
			},
			wantErr: true,
		},
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
	dataItemsPgRepo, dataItemsS3Repo, err := testRepos()
	if err != nil {
		t.Fatal(err)
	}

	testModel := testModelEdit()
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
		{
			name: "Delete new data items service",
			fields: fields{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			args: args{
				ctx: context.Background(),
				pars: &model.GetPars{
					ID: testModel.ID,
				},
				obj: testModel,
			},
			wantErr: false,
		},
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

			if err = s.Delete(tt.args.ctx, tt.args.pars); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	dataItemsPgRepo, dataItemsS3Repo, err := testRepos()
	if err != nil {
		t.Fatal(err)
	}

	testModel := testModelEdit()
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
		want    *model.Main
		want1   bool
		wantErr bool
	}{
		{
			name: "Get new data items service",
			fields: fields{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			args: args{
				ctx: context.Background(),
				pars: &model.GetPars{
					ID: testModel.ID,
				},
				obj: testModel,
			},
			want: &model.Main{
				ID: testModel.ID,
			},
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: tt.fields.repoDB,
				repoS3: tt.fields.repoS3,
			}

			if err = s.Create(tt.args.ctx, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, got1, err := s.Get(tt.args.ctx, tt.args.pars)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	dataItemsPgRepo, dataItemsS3Repo, err := testRepos()
	if err != nil {
		t.Fatal(err)
	}

	testModel := testModelEdit()
	type fields struct {
		repoDB RepoDBI
		repoS3 RepoS3
	}
	type args struct {
		ctx  context.Context
		pars *model.ListPars
		obj  *model.Edit
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want1   int64
		wantErr bool
	}{
		{
			name: "List new data items service",
			fields: fields{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			args: args{
				ctx: context.Background(),
				pars: &model.ListPars{
					UserID: testModel.UserID,
				},
				obj: testModel,
			},
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: tt.fields.repoDB,
				repoS3: tt.fields.repoS3,
			}

			if err = s.Create(tt.args.ctx, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, got1, err := s.List(tt.args.ctx, tt.args.pars)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got1 != tt.want1 {
				t.Errorf("List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	dataItemsPgRepo, dataItemsS3Repo, err := testRepos()
	if err != nil {
		t.Fatal(err)
	}

	testModel := testModelEdit()
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
		{
			name: "Update new data items service",
			fields: fields{
				repoDB: dataItemsPgRepo,
				repoS3: dataItemsS3Repo,
			},
			args: args{
				ctx: context.Background(),
				pars: &model.GetPars{
					ID: testModel.ID,
				},
				obj: testModel,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: tt.fields.repoDB,
				repoS3: tt.fields.repoS3,
			}

			if err = s.Create(tt.args.ctx, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			newValue := "test"
			tt.args.obj.URL = &newValue

			if err := s.Update(tt.args.ctx, tt.args.pars, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func testRepos() (*dataItemsRepoPgP.Repo, *dataItemsRepoS3P.S3Repo, error) {
	pgpool, err := getPgPoolTestContainer()
	if err != nil {
		return nil, nil, err
	}

	endpoint, accessKey, secretKey, bucketName, err := setupMinio(context.Background())
	if err != nil {
		return nil, nil, err
	}

	dataItemsPgRepo := dataItemsRepoPgP.New(pgpool)
	dataItemsS3Repo, err := dataItemsRepoS3P.NewS3Repo(context.Background(), endpoint, accessKey, secretKey, bucketName)
	if err != nil {
		return nil, nil, err
	}

	return dataItemsPgRepo, dataItemsS3Repo, nil
}

func testModelEdit() *model.Edit {
	id := uuid.New().String()
	testUserID := "999"
	binaryType := model.BinaryDataType
	data := []byte("test/test")
	meta := "binary"

	return &model.Edit{
		ID:     id,
		UserID: &testUserID,
		Type:   &binaryType,
		Data:   &data,
		Meta:   &meta,
	}
}
