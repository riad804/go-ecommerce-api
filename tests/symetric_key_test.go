package tests

import (
	"encoding/base64"
	"testing"

	"github.com/riad804/go_ecommerce_api/token"
)

func TestSymmetricKey(t *testing.T) {
	validKeyBase64 := "kbb3ObClKKHAZClXy7z3KfBaqJLey2ydGLfV2heOUz8="
	shortKeyBase64 := "c2hvcnRrZXk="
	invalidBase64 := "!@#$$%"

	tests := []struct {
		name      string
		base64Key string
		wantErr   bool
	}{
		{"valid key", validKeyBase64, false},
		{"too short", shortKeyBase64, true},
		{"invalid base64", invalidBase64, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := base64.StdEncoding.DecodeString(tt.base64Key)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected decode error: %v", err)
				return
			}
			if err == nil {
				maker, err := token.NewPasetoMaker(tt.base64Key)
				if (err != nil) != tt.wantErr {
					t.Errorf("NewPasetoMaker() error = %v, wantErr = %v", err, tt.wantErr)
				}
				if err == nil && maker == nil {
					t.Errorf("expected a valid PasetoMaker, got nil")
				}
			}
		})
	}
}
