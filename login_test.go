package ctd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtd_Login(t *testing.T) {
	ctx := t.Context()

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

func TestCtd_newLoginAPIParsing(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name    string
		result  string
		want    string
		wantErr error
	}{
		{
			name: "Password has expired",
			result: `{
  				"status": "error",
  				"errors": {
    			"password": [
	      			"{\"access\":false,\"reason\":\"password\",\"message\":\"change_password_required\",\"confirm_hash\":\"d8f3187f13c69aa0f0395f640701fa33\",\"password_expired\":\"Your password has expired. Please change your password. You will be logged in immediately after changing your password.\"}"
    			],
    			"password_expired": [
	    	  		"Your password has expired. Please change your password. You will be logged in immediately after changing your password."
    			]
  				},
  				"login_attempts_info": {
    				"max_attempts_number": 5,
    				"failed_attempts_number": 17,
    				"failed_login_attempt_date": 1690804029
  				}
			}`,
			want:    "",
			wantErr: ErrorPasswordHasExpired,
		},
		{
			name:    "Incorrect OTP",
			result:  `{"status":"error","errors":{"error":["incorrect_otp"],"status_code":[401],"otp_required":["Enter one time password"]},"login_attempts_info":{"max_attempts_number":5,"failed_attempts_number":1,"failed_login_attempt_date":1769423465}}`,
			want:    "",
			wantErr: ErrorOTPRequired,
		},
		{
			name:    "Successful login",
			result:  `{"status":"success","session_code":"617.PqOCydEPe7p2fkexKSgb","auth_key":"PqOCydEPe7p2fkexKSgb.user"}`,
			want:    "",
			wantErr: nil,
		},
		{
			name:    "Incorrect login or password",
			result:  `{"status":"error","errors":{"error":["incorrect_password"],"status_code":[401],"password":["Wrong login or password"]},"login_attempts_info":{"max_attempts_number":5,"failed_attempts_number":3,"failed_login_attempt_date":1769433935}}`,
			want:    "",
			wantErr: ErrorInvalidLoginOrPassword,
		},
		{
			name:    "User not found",
			result:  `{"status":"error","errors":{"error":["user_does_not_exist"],"status_code":[401],"password":["Wrong login or password"]},"login_attempts_info":null}`,
			want:    "",
			wantErr: ErrorUserNotFound,
		},
	}
	ctd := Ctd{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctd.newLoginAPIParsing(ctx, []byte(tt.result))
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr, "Error mismatch")
				require.Equal(t, tt.want, got, "Result mismatch")
			} else {
				require.NoError(t, err, "newLoginAPIParsing() error")
				require.Equal(t, tt.want, got, "Result mismatch")
			}
		})
	}
}
