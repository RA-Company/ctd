package ctd

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/ra-company/logging"
)

var (
	ErrorInvalidLoginOrPassword = fmt.Errorf("invalid login or password")
	ErrorTooFast                = fmt.Errorf("too fast")
	ErrorAccountIsBlocked       = fmt.Errorf("account is blocked")
	ErrorMasterPasswordRequired = fmt.Errorf("master password required")
	ErrorMasterNotAllowed       = fmt.Errorf("master not allowed")
	ErrorOTPRequired            = fmt.Errorf("otp required")
	ErrorCaptchRequired         = fmt.Errorf("captcha required")
	ErrorUserNotFound           = fmt.Errorf("user not found")
)

func (dst *Ctd) Login(ctx context.Context, login, password, personal_login, personal_password, otp, captcha string) (string, error) {
	data := struct {
		Email            string `json:"email"`              // Chat2Desk login
		Password         string `json:"password"`           // Chat2Desk password
		PersonalEmail    string `json:"personal_email"`     // Personal email for master login
		PersonalPassword string `json:"personal_password"`  // Personal password for master login
		OTP              string `json:"oneTimePassword"`    // One-time password
		Captcha          string `json:"grecaptchaResponse"` // Captcha response
	}{
		Email:            login,
		Password:         password,
		PersonalEmail:    personal_login,
		PersonalPassword: personal_password,
		OTP:              otp,
		Captcha:          captcha,
	}

	type ctdAuthResponse struct {
		Status  string `json:"status"`
		AuthKey string `json:"auth_key"`
		Errors  struct {
			Error      []string `json:"error"`
			StatusCode []any    `json:"status_code"`
		} `json:"errors"`
		Attempts any `json:"login_attempts_info"`
	}

	url := fmt.Sprintf("%sapi/user/sign_in?lang=en", dst.Url)

	if personal_login != "" && personal_password != "" {
		url = fmt.Sprintf("%sapi/user/master_sign_in?lang=en", dst.Url)
	}

	result, err := dst.doRequest(ctx, "POST", url, data, nil)
	if err != nil {
		dst.Error(ctx, "Failed login: %v", err)
		return "", err
	}

	logging.Logs.Debugf(ctx, "Response: %s", string(result))

	var response ctdAuthResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		logging.Logs.Errorf(ctx, "%v", err)
		return "", ErrorInvalidResponse
	}

	if strings.ToLower(response.Status) == "success" {
		return "", nil
	}

	// New version of Chat2Desk API
	if slices.Index(response.Errors.Error, "user_does_not_exist") != -1 {
		return "", ErrorUserNotFound
	}

	if slices.Index(response.Errors.Error, "incorrect_otp") != -1 {
		return "", ErrorOTPRequired
	}

	if slices.Index(response.Errors.Error, "captcha") != -1 {
		return "", ErrorCaptchRequired
	}

	if slices.Index(response.Errors.Error, "incorrect_password") != -1 {
		return "", ErrorInvalidLoginOrPassword
	}

	if slices.Index(response.Errors.Error, "timeout") != -1 {
		return "10", ErrorTooFast
	}

	// Old version of Chat2Desk API
	if response.Attempts == nil {
		return "", ErrorUserNotFound
	}

	str := strings.ToLower(string(result))

	if strings.Contains(str, "wrong login or password") {
		return "", ErrorInvalidLoginOrPassword
	}

	if strings.Contains(str, "this account is blocked") {
		return "", ErrorAccountIsBlocked
	}

	if strings.Contains(str, "\"master_login\":[true]") {
		return "", ErrorMasterPasswordRequired
	}

	if strings.Contains(str, "login with master-password is not permitted using this method") {
		return "", ErrorMasterPasswordRequired
	}

	if strings.Contains(str, "access under master password is not allowed by the account administrator") {
		return "", ErrorMasterNotAllowed
	}

	if strings.Contains(str, "enter one time password") {
		return "", ErrorOTPRequired
	}

	if strings.Contains(str, "please, enter captcha to log in") {
		return "", ErrorCaptchRequired
	}

	if strings.Contains(str, "user_does_not_exist") {
		return "", ErrorUserNotFound
	}

	if strings.Contains(str, "e-mail is not a valid email address") {
		return "", ErrorUserNotFound
	}

	if strings.Contains(str, "please, try again after") {
		i := strings.Index(str, "after")
		j := strings.Index(str, "second")
		return str[(i + 6):(j - 1)], ErrorTooFast
	}

	// Very Old version of Chat2Desk API
	type ctdAuthResponseOld struct {
		Status   string           `json:"status"`
		AuthKey  string           `json:"auth_key"`
		Errors   map[string][]any `json:"errors"`
		Attempts any              `json:"login_attempts_info"`
	}

	var responseOld ctdAuthResponseOld
	err = json.Unmarshal(result, &responseOld)
	if err != nil {
		logging.Logs.Errorf(ctx, "%v", err)
		return "", ErrorInvalidResponse
	}

	for key, value := range responseOld.Errors {
		switch key {
		case "password":
			for _, v := range value {
				str := strings.ToLower(fmt.Sprintf("%v", v))
				if strings.Contains(str, "wrong login or password") {
					return "", ErrorInvalidLoginOrPassword
				}
				if strings.Contains(str, "this account is blocked") {
					return "", ErrorAccountIsBlocked
				}
				if strings.Contains(str, "login with master-password is not permitted using this method") {
					return "", ErrorMasterPasswordRequired
				}
				if strings.Contains(str, "access under master password is not allowed by the account administrator") {
					return "", ErrorMasterNotAllowed
				}
			}
		case "brute_force":
			for _, v := range value {
				str := strings.ToLower(fmt.Sprintf("%v", v))
				if strings.Contains(str, "please, try again after") {
					i := strings.Index(str, "after")
					j := strings.Index(str, "second")
					return str[(i + 6):(j - 1)], ErrorTooFast
				}
				if strings.Contains(str, "please, enter captcha to log in") {
					return "", ErrorCaptchRequired
				}
			}
		case "one_time_password":
			for _, v := range value {
				str := strings.ToLower(fmt.Sprintf("%v", v))
				if strings.Contains(str, "must be filled") {
					return "", ErrorOTPRequired
				}
			}
		}
	}

	return fmt.Sprintf("%v", response.Errors), ErrorUnknownError
}
