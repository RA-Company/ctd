package ctd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtdApi_CompaniesApiInfo(t *testing.T) {
	ctx := t.Context()

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
			got, err := dst.CompaniesApiInfo(ctx)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.CompaniesApiInfo() error")

			} else {

			}
			if tt.isData {
				require.NotNil(t, got, "dst.CompaniesApiInfo() should return data")
				require.NotNil(t, got.CompanyName, "got.CompanyName should not be nil")
				require.NotNil(t, got.AdminEmail, "got.AdminEmail should not be nil")
				require.NotEqual(t, 0, got.CompanyID, "got.CompanyID should not be 0")
				require.NotEqual(t, 0, got.PartnerID, "got.PartnerID should not be 0")
			} else {
				require.Nil(t, got, "dst.CompaniesApiInfo() should return nil data on error")
			}
		})
	}
}
