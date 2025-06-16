package ctd

import (
	"context"
	"testing"

	"github.com/ra-company/env"
	"github.com/stretchr/testify/require"
)

func TestCtd_Channels(t *testing.T) {
	ctx := context.Background()

	url, token := getCredentials(t)

	type args struct {
		url   string
		token string
	}
	tests := []struct {
		name   string
		args   args
		isData bool
		error  error
	}{
		{
			name: "Incorrect token",
			args: args{
				url:   url,
				token: "incorrect_token",
			},
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name: "Correct token",
			args: args{
				url:   url,
				token: token,
			},
			isData: true,
			error:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(tt.args.url, tt.args.token)
			got, err := dst.Channels(ctx, 0, 100)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.Channel() error")

			} else {

			}
			if tt.isData {
				require.NotNil(t, got, "dst.Channel() should return data")
				require.NotEqual(t, 0, got.Meta.Total, "dst.Channel() should return non-empty total count")
				require.Equal(t, 100, got.Meta.Limit, "dst.Channel() should return correct limit")
			} else {
				require.Nil(t, got, "dst.Channel() should return nil data on error")
			}

		})
	}
}

func getCredentials(t *testing.T) (string, string) {
	url := env.GetEnvUrl("API_URL", "")
	require.NotEqual(t, "", url, "API_URL must be set in .env file or .settings")

	token := env.GetEnvStr("API_TOKEN", "")
	require.NotEqual(t, "", token, "API_TOKEN must be set in .env file or .settings")

	return url, token
}
