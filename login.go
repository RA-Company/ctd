package ctd

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
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
	ErrorPasswordHasExpired     = fmt.Errorf("password has expired")
)

// LoginResult is the result of Login.
type LoginResult struct {
	AuthKey    string // AuthKey: authorization key, set on successful login
	RetryDelay int    // RetryDelay: number of seconds to wait before the next attempt, set together with ErrorTooFast
	Error      string // Error: raw response body, set together with ErrorUnknownError
}

// Login signs in to Chat2Desk.
//
// Returns:
//   - LoginResult with AuthKey on success, with RetryDelay (seconds) when ErrorTooFast is returned,
//     or with the raw response body in Error when ErrorUnknownError is returned.
//   - An error if the login fails.
func (dst *Ctd) Login(ctx context.Context, login, password, personal_login, personal_password, otp, captcha string) (res LoginResult, err error) {
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

	url := fmt.Sprintf("%sapi/user/sign_in?lang=en", dst.Url)

	if personal_login != "" && personal_password != "" {
		url = fmt.Sprintf("%sapi/user/master_sign_in?lang=en", dst.Url)
	}

	result, err := dst.doRequest(ctx, "POST", url, data, nil)
	if err != nil {
		dst.Errorf(ctx, "Failed login: %v", err)
		return LoginResult{}, err
	}

	// The response of a failed login contains no auth key, so it is safe to log it.
	defer func() {
		if err != nil {
			logging.Logs.Debugf(ctx, "Response: %s", string(result))
		}
	}()

	response, err := dst.newLoginAPIParsing(ctx, result)
	if err == nil {
		return response, nil
	}

	if err != ErrorUnknownError {
		return response, err
	}

	// Old version of Chat2Desk API
	str := strings.ToLower(string(result))

	if strings.Contains(str, "wrong login or password") {
		return LoginResult{}, ErrorInvalidLoginOrPassword
	}

	if strings.Contains(str, "this account is blocked") {
		return LoginResult{}, ErrorAccountIsBlocked
	}

	if strings.Contains(str, "\"master_login\":[true]") {
		return LoginResult{}, ErrorMasterPasswordRequired
	}

	if strings.Contains(str, "login with master-password is not permitted using this method") {
		return LoginResult{}, ErrorMasterPasswordRequired
	}

	if strings.Contains(str, "access under master password is not allowed by the account administrator") {
		return LoginResult{}, ErrorMasterNotAllowed
	}

	if strings.Contains(str, "enter one time password") {
		return LoginResult{}, ErrorOTPRequired
	}

	if strings.Contains(str, "please, enter captcha to log in") {
		return LoginResult{}, ErrorCaptchRequired
	}

	if strings.Contains(str, "user_does_not_exist") {
		return LoginResult{}, ErrorUserNotFound
	}

	if strings.Contains(str, "e-mail is not a valid email address") {
		return LoginResult{}, ErrorUserNotFound
	}

	if strings.Contains(str, "password_expired") {
		return LoginResult{}, ErrorPasswordHasExpired
	}

	if strings.Contains(str, "please, try again after") {
		return LoginResult{RetryDelay: tooFastDelay(str)}, ErrorTooFast
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
		return LoginResult{}, ErrorInvalidResponse
	}

	for key, value := range responseOld.Errors {
		switch key {
		case "password":
			for _, v := range value {
				str := strings.ToLower(fmt.Sprintf("%v", v))
				if strings.Contains(str, "wrong login or password") {
					return LoginResult{}, ErrorInvalidLoginOrPassword
				}
				if strings.Contains(str, "this account is blocked") {
					return LoginResult{}, ErrorAccountIsBlocked
				}
				if strings.Contains(str, "login with master-password is not permitted using this method") {
					return LoginResult{}, ErrorMasterPasswordRequired
				}
				if strings.Contains(str, "access under master password is not allowed by the account administrator") {
					return LoginResult{}, ErrorMasterNotAllowed
				}
			}
		case "brute_force":
			for _, v := range value {
				str := strings.ToLower(fmt.Sprintf("%v", v))
				if strings.Contains(str, "please, try again after") {
					return LoginResult{RetryDelay: tooFastDelay(str)}, ErrorTooFast
				}
				if strings.Contains(str, "please, enter captcha to log in") {
					return LoginResult{}, ErrorCaptchRequired
				}
			}
		case "one_time_password":
			for _, v := range value {
				str := strings.ToLower(fmt.Sprintf("%v", v))
				if strings.Contains(str, "must be filled") {
					return LoginResult{}, ErrorOTPRequired
				}
			}
		}
	}

	return LoginResult{Error: string(result)}, ErrorUnknownError
}

func (dst *Ctd) newLoginAPIParsing(ctx context.Context, result []byte) (LoginResult, error) {
	type ctdAuthResponse struct {
		Status  string `json:"status"`
		AuthKey string `json:"auth_key"`
		Errors  struct {
			Error           json.RawMessage `json:"error"`
			StatusCode      []any           `json:"status_code"`
			PasswordExpired []string        `json:"password_expired"`
		} `json:"errors"`
		Attempts any `json:"login_attempts_info"`
	}

	var response ctdAuthResponse

	if err := json.Unmarshal(result, &response); err != nil {
		logging.Logs.Errorf(ctx, "Ctd.newLoginAPIParsing->json.Unmarshal() error: %v", err)
		return LoginResult{}, ErrorInvalidResponse
	}

	if strings.ToLower(response.Status) == "success" {
		return LoginResult{AuthKey: response.AuthKey}, nil
	}

	if delay, err := dst.parseLoginErrorSrings(response.Errors.Error); err != nil {
		if err != ErrorInvalidResponse {
			return LoginResult{RetryDelay: delay}, err
		}
	}

	if len(response.Errors.PasswordExpired) > 0 {
		return LoginResult{}, ErrorPasswordHasExpired
	}

	if response.Attempts == nil {
		return LoginResult{}, ErrorUserNotFound
	}

	return LoginResult{}, ErrorUnknownError
}

func (dst *Ctd) parseLoginErrorSrings(data json.RawMessage) (int, error) {
	results := []string{}

	if err := json.Unmarshal(data, &results); err != nil {
		return 0, ErrorInvalidResponse
	}

	// New version of Chat2Desk API
	if slices.Index(results, "user_does_not_exist") != -1 {
		return 0, ErrorUserNotFound
	}

	if slices.Index(results, "incorrect_otp") != -1 {
		return 0, ErrorOTPRequired
	}

	if slices.Index(results, "captcha") != -1 {
		return 0, ErrorCaptchRequired
	}

	if slices.Index(results, "incorrect_password") != -1 {
		return 0, ErrorInvalidLoginOrPassword
	}

	if slices.Index(results, "timeout") != -1 {
		return 10, ErrorTooFast
	}

	return 0, nil
}

// tooFastDelay extracts the delay from a "please, try again after N seconds" message.
// It returns 0 if the message has an unexpected format.
func tooFastDelay(str string) int {
	_, rest, found := strings.Cut(str, "after ")
	if !found {
		return 0
	}

	delay, _, found := strings.Cut(rest, " second")
	if !found {
		return 0
	}

	seconds, err := strconv.Atoi(strings.TrimSpace(delay))
	if err != nil {
		return 0
	}

	return seconds
}
