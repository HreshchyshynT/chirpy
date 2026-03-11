package user_test

import (
	"github.com/hreshchyshynt/chirpy/internal/config"
	"github.com/hreshchyshynt/chirpy/internal/user"
	"net/http"
	"testing"
)

func TestHandlePolkaWebhook(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w         http.ResponseWriter
		r         *http.Request
		apiConfig *config.ApiConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user.HandlePolkaWebhook(tt.w, tt.r, tt.apiConfig)
		})
	}
}
