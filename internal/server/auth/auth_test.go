package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_jwtTokenCreate(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{
			name:   "create token",
			userID: 27,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			jwtToken := jwtToken{}
			_, err := jwtToken.Create(tt.userID)

			if !tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}
