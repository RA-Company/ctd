package ctd

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtd_Login(t *testing.T) {
	ctx := context.Background()

	login := os.Getenv("TEST_CTD_LOGIN")
	require.NotEqual(t, "", login, "Environment variable TEST_CTD_LOGIN isn't set")

	url := os.Getenv("TEST_WEB_URL")
	require.NotEqual(t, "", url, "Environment variable TEST_WEB_URL isn't set")

	pwd := os.Getenv("TEST_CTD_PWD")
	require.NotEqual(t, "", pwd, "Environment variable TEST_CTD_PWD isn't set")

	ctd := Ctd{
		Url: url,
	}

	t.Run("01 - Login with incorrect username", func(t *testing.T) {
		_, err := ctd.Login(ctx, "incorrect", "incorrect", "", "", "", "")
		require.ErrorIs(t, err, ErrorUserNotFound)
	})

	t.Run("02 - Login with incorrect password", func(t *testing.T) {
		_, err := ctd.Login(ctx, login, "incorrect", "", "", "", "")
		require.ErrorIs(t, err, ErrorInvalidLoginOrPassword)
	})

	t.Run("03 - Login with correct credentials and OTP", func(t *testing.T) {
		_, err := ctd.Login(ctx, login, pwd, "", "", "", "")
		require.ErrorIs(t, err, ErrorOTPRequired)
	})
}
