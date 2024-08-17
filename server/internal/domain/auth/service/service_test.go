package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"gophKeeper/server/internal/domain/users/model"
	"reflect"
	"testing"
)

func TestAuth_CreateToken(t *testing.T) {
	type fields struct {
		JwtSecret string
	}
	type args struct {
		u *model.Main
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test for CreateToken",
			fields: fields{
				JwtSecret: "secret",
			},
			args: args{
				u: &model.Main{
					UserID:   "999",
					Username: "Test",
				},
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				JwtSecret: tt.fields.JwtSecret,
			}
			got, err := a.CreateToken(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.want {
				t.Errorf("CreateToken() got = %v, want not empty", got)
			}
		})
	}
}

func TestAuth_GetUserIDFromContext(t *testing.T) {
	invalidToken := "invalid-jwt-token"

	u := &model.Main{
		UserID:   "999",
		Username: "Test",
	}
	a := &Auth{
		JwtSecret: "testsecret",
	}
	got, err := a.NewToken(u)
	if err != nil {
		t.Errorf("NewToken() error = %v", err)
	}

	tests := []struct {
		name       string
		jwtSecret  string
		metadata   metadata.MD
		expectedID string
		expectErr  bool
	}{
		{
			name:      "valid token",
			jwtSecret: "",
			metadata: metadata.Pairs(
				"token", "Bearer "+got,
			),
			expectedID: "999", // Предполагаемый UserID, который вы ожидаете получить из валидного токена
			expectErr:  false,
		},
		{
			name:      "missing metadata",
			jwtSecret: "your-secret-key",
			metadata:  nil,
			expectErr: true,
		},
		{
			name:      "missing token in metadata",
			jwtSecret: "your-secret-key",
			metadata:  metadata.Pairs(),
			expectErr: true,
		},
		{
			name:      "invalid token format",
			jwtSecret: "your-secret-key",
			metadata: metadata.Pairs(
				"token", "Bearer invalid-format-token",
			),
			expectErr: true,
		},
		{
			name:      "invalid token signature",
			jwtSecret: "your-secret-key",
			metadata: metadata.Pairs(
				"token", "Bearer "+invalidToken,
			),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.metadata != nil {
				ctx = metadata.NewIncomingContext(ctx, tt.metadata)
			}

			userID, err := a.GetUserIDFromContext(ctx)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}

func TestAuth_NewToken(t *testing.T) {
	type fields struct {
		JwtSecret string
	}
	type args struct {
		u *model.Main
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test for NewToken",
			fields: fields{
				JwtSecret: "secret",
			},
			args: args{
				u: &model.Main{
					UserID:   "999",
					Username: "Test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				JwtSecret: tt.fields.JwtSecret,
			}
			got, err := a.NewToken(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("NewToken() got = %v", got)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		jwtSecret string
	}
	tests := []struct {
		name string
		args args
		want *Auth
	}{
		{
			name: "Create new auth service",
			args: args{
				jwtSecret: "secret",
			},
			want: &Auth{
				JwtSecret: "secret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.jwtSecret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
