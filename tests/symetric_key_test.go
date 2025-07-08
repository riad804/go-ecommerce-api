package tests

import (
	"encoding/base64"
	"testing"

	"github.com/riad804/go_ecommerce_api/token"
)

func TestSymmetricKey(t *testing.T) {
	accessKeyBase64 := "u+JGOLu6JG7q8apsVV1UXe74agrSHqVCJIeD/ad3tSQ="
	refreshKeyBase64 := "GFTR4lalgqPCUKpkuAKC7VUwE6yM0d4Yqe3Fg3hRn9I="
	invalidBase64 := "!@#$$%"

	tests := []struct {
		name      string
		base64Key string
		wantErr   bool
	}{
		{"access key", accessKeyBase64, false},
		{"refresh key", refreshKeyBase64, false},
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
				maker, err := token.NewPasetoMaker(tt.base64Key, tt.base64Key)
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
