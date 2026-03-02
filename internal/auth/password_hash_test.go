package auth_test

import (
	"testing"

	"github.com/hreshchyshynt/chirpy/internal/auth"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		password string
	}{
		{
			password: "12345Kja@asdf.am",
		},
		{
			password: "asldkfjasdkjfcm,vm",
		},
		{
			password: "1238172348123",
		},
		{
			password: "@*#&@*&#@&#@)#()@()",
		},
	}
	t.Run("check hashing password match", func(t *testing.T) {
		for _, tt := range tests {
			got, gotErr := auth.HashPassword(tt.password)
			if gotErr != nil {
				t.Errorf("Error received during hashing password: %v", gotErr)
			}
			match, err := auth.CheckPasswordHash(tt.password, got)
			if err != nil {
				t.Errorf("Error received during checking password: %v", err)
			}
			if !match {
				t.Errorf("HashPassword() = %v, match: %v", got, match)
			}

		}

	})

}
