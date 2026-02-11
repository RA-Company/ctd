package ctd

import (
	"testing"

	"github.com/ra-company/env"
	"github.com/stretchr/testify/require"
)

func TestCtd_RequestMessages(t *testing.T) {
	ctx := t.Context()
	url, token := getCredentials(t)

	requestID := env.GetEnvInt("API_REQUEST_ID", 0)
	require.NotEqual(t, 0, requestID, "API_REQUEST_ID must be set in .env file or .settings")

	tests := []struct {
		name     string
		token    string
		request  int64
		wantData bool
		wantErr  error
	}{
		{
			name:     "Invalid token",
			token:    "invalid_token",
			request:  12345,
			wantData: false,
			wantErr:  ErrorInvalidToken,
		},
		{
			name:     "Valid token but request not found",
			token:    token,
			request:  1,
			wantData: false,
			wantErr:  ErrorInvalidRequestID,
		},
		{
			name:     "Valid token and request",
			token:    token,
			request:  int64(requestID),
			wantData: true,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(url, tt.token)
			got, err := dst.RequestMessages(ctx, tt.request)
			if tt.wantErr == nil {
				require.NoError(t, err, "dst.RequestMessages() error")
			} else {
				require.ErrorIs(t, err, tt.wantErr, "dst.RequestMessages() error")
			}
			if tt.wantData {
				require.NotNil(t, got, "dst.RequestMessages() should return data")
			} else {
				require.Nil(t, got, "dst.RequestMessages() should return nil data on error")
			}
		})
	}
}
