package ctd

import (
	"context"
	"testing"

	"github.com/ra-company/env"
	"github.com/stretchr/testify/require"
)

func TestCtd_Clients(t *testing.T) {
	ctx := context.Background()

	url, token := getCredentials(t)

	clientID := env.GetEnvInt("API_CLIENT_ID", 0)
	require.NotEqual(t, 0, clientID, "API_CLIENT_ID must be set in .env file or .settings")

	tests := []struct {
		name      string
		token     string
		client_id int
		isData    bool
		error     error
	}{
		{
			name:   "Incorrect token",
			token:  "incorrect_token",
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name:      "Correct token",
			token:     token,
			client_id: clientID,
			isData:    true,
			error:     nil,
		},
		{
			name:      "Client not found",
			token:     token,
			client_id: 0, // Assuming 0 is an invalid ID for testing
			isData:    false,
			error:     ErrorInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(url, tt.token)
			got, err := dst.GetClient(ctx, tt.client_id)
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

	t.Run("GetClientsList", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)

		offset := 0
		limit := 100

		got, total, err := dst.GetClientsList(ctx, offset, limit)
		require.NoError(t, err, "dst.GetClientsList() should not return an error")
		require.NotNil(t, got, "dst.GetClientsList() should return data")
		require.Greater(t, len(*got), 0, "dst.GetClientsList() should return non-empty client list")
		require.Greater(t, total, 0, "dst.GetClientsList() should return non-zero total count")

		offset += limit
		got2, total2, err := dst.GetClientsList(ctx, offset, limit)
		require.NoError(t, err, "dst.GetClientsList() should not return an error on second call")
		require.NotNil(t, got2, "dst.GetClientsList() should return data on second call")
		require.Greater(t, len(*got2), 0, "dst.GetClientsList() should return non-empty client list on second call")
		require.Greater(t, total2, 0, "dst.GetClientsList() should return non-zero total count on second call")
		require.Equal(t, total, total2, "dst.GetClientsList() should return the same total count on both calls")
		require.NotEqual(t, got, got2, "dst.GetClientsList() should return different data on different calls with different offsets")
	})
}
