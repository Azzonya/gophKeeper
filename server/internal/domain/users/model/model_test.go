package model

import "testing"

func TestGetPars_IsValid(t *testing.T) {
	type fields struct {
		UserID   string
		Username string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "test user model valid",
			fields: fields{
				UserID: "111",
			},
			want: true,
		},
		{
			name:   "test user model not valid",
			fields: fields{},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GetPars{
				UserID:   tt.fields.UserID,
				Username: tt.fields.Username,
			}
			if got := m.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
