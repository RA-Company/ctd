package ctd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtdApi_Operators(t *testing.T) {
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
			got, total, err := dst.Operators(ctx, 0, 10)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.Operators() error")
				require.Zero(t, total, "dst.Operators() should return zero total on error")
			} else {
				require.NoError(t, err, "dst.Operators() error")
				require.Greater(t, total, 0, "dst.Operators() should return some total")
			}
			if tt.isData {
				require.NotNil(t, got, "dst.Operators() should return data")
			} else {
				require.Nil(t, got, "dst.Operators() should return nil data on error")
			}

		})
	}
}
