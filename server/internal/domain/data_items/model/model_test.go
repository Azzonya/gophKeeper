package model

import "testing"

func TestGetPars_IsValid(t *testing.T) {
	type fields struct {
		ID     string
		UserID string
		Type   string
		Meta   string
		URL    string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid struct",
			fields: fields{
				ID: "32",
			},
			want: true,
		},
		{
			name:   "invalid struct",
			fields: fields{},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GetPars{
				ID:     tt.fields.ID,
				UserID: tt.fields.UserID,
				Type:   tt.fields.Type,
				Meta:   tt.fields.Meta,
				URL:    tt.fields.URL,
			}
			if got := m.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
