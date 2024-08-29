package service

import "testing"

func TestService_HashPassword(t *testing.T) {

	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test hash password",
			args: args{
				"mysecretpassword",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: nil,
			}
			got, err := s.HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got1, err := s.HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == got1 {
				t.Errorf("expected different hashes for the same password, but got the same")
			}
		})
	}
}

func TestService_IsValidPassword(t *testing.T) {
	type args struct {
		password      string
		plainPassword string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test not valid password",
			args: args{
				password:      "mysecretpassword",
				plainPassword: "mysecretpassword",
			},
			want: false,
		},
		{
			name: "test valid password",
			args: args{
				password: "123",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repoDB: nil,
			}

			if tt.want {
				got, err := s.HashPassword(tt.args.password)
				if err != nil {
					t.Errorf("HashPassword() error = %v", err)
					return
				}
				tt.args.plainPassword = got
			}

			if got := s.IsValidPassword(tt.args.plainPassword, tt.args.password); got != tt.want {
				t.Errorf("IsValidPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
