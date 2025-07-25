package ctd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtd_CustomClientFields(t *testing.T) {
	ctx := context.Background()

	url, token := getCredentials(t)

	tests := []struct {
		name   string
		token  string
		isData bool
		error  error
	}{
		{
			name:   "Incorrect token",
			token:  "incorrect_token",
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name:   "Correct token",
			token:  token,
			isData: true,
			error:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(url, tt.token)
			got, err := dst.GetCustomClientFields(ctx)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.GetClient() error")

			} else {

			}
			if tt.isData {
				require.NotNil(t, got, "dst.GetClient() should return data")
			} else {
				require.Nil(t, got, "dst.GetClient() should return nil data on error")
			}

		})
	}
}
