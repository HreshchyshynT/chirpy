package auth_test

import (
	"github.com/hreshchyshynt/chirpy/internal/auth"
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name: "No Authorization header fails",
			headers: http.Header{
				"Content-type": []string{"Application/json"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "No bearer prefix fails",
			headers: http.Header{
				"Content-type":  []string{"Application/json"},
				"Authorization": []string{"token"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Token returns without Bearer",
			headers: http.Header{
				"Content-type":  []string{"Application/json"},
				"Authorization": []string{"Bearer my-token"},
			},
			want:    "my-token",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}

			if tt.want != got {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
