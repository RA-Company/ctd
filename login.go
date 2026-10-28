package ctd

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrorInvalidLoginOrPassword = fmt.Errorf("invalid login or password")
	ErrorTooFast                = fmt.Errorf("too fast")
	ErrorAccountIsBlocked       = fmt.Errorf("account is blocked")
	ErrorMasterPasswordRequired = fmt.Errorf("master password required")
	ErrorMasterNotAllowed       = fmt.Errorf("master not allowed")
	ErrorOTPRequired            = fmt.Errorf("otp required")
	ErrorCaptchaRequired        = fmt.Errorf("captcha required")
	ErrorUserNotFound           = fmt.Errorf("user not found")
	ErrorPasswordHasExpired     = fmt.Errorf("password has expired")
)

// legacyLoginRules maps the texts of the responses of the old versions of Chat2Desk API to errors.
// The texts are lower case and matched against the whole response in the given order, the first match wins.
var legacyLoginRules = []struct {
	text string
	err  error
}{
	{"wrong login or password", ErrorInvalidLoginOrPassword},
	{"this account is blocked", ErrorAccountIsBlocked},
	{`"master_login":[true]`, ErrorMasterPasswordRequired},
	{"login with master-password is not permitted using this method", ErrorMasterPasswordRequired},
	{"access under master password is not allowed by the account administrator", ErrorMasterNotAllowed},
	{"enter one time password", ErrorOTPRequired},
	{"please, enter captcha to log in", ErrorCaptchaRequired},
	{"user_does_not_exist", ErrorUserNotFound},
	{"e-mail is not a valid email address", ErrorUserNotFound},
	{"password_expired", ErrorPasswordHasExpired},
}

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
func (dst *Ctd) Login(ctx context.Context, login, password, personalLogin, personalPassword, otp, captcha string) (res LoginResult, err error) {
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
		PersonalEmail:    personalLogin,
		PersonalPassword: personalPassword,
		OTP:              otp,
		Captcha:          captcha,
	}

	url := fmt.Sprintf("%sapi/user/sign_in?lang=en", dst.Url)

	if personalLogin != "" && personalPassword != "" {
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
			dst.Debugf(ctx, "Response: %s", string(result))
		}
	}()

	response, err := dst.newLoginAPIParsing(ctx, result)
	if err == nil {
		return response, nil
	}

	if err != ErrorUnknownError {
		return response, err
	}

	// Old versions of Chat2Desk API
	str := strings.ToLower(string(result))

	for _, rule := range legacyLoginRules {
		if strings.Contains(str, rule.text) {
			return LoginResult{}, rule.err
		}
	}

	if strings.Contains(str, "please, try again after") {
		return LoginResult{RetryDelay: tooFastDelay(str)}, ErrorTooFast
	}

	var responseOld struct {
		Errors map[string][]any `json:"errors"`
	}

	if err = json.Unmarshal(result, &responseOld); err != nil {
		dst.Errorf(ctx, "%v", err)
		return LoginResult{}, ErrorInvalidResponse
	}

	for _, v := range responseOld.Errors["one_time_password"] {
		if strings.Contains(strings.ToLower(fmt.Sprintf("%v", v)), "must be filled") {
			return LoginResult{}, ErrorOTPRequired
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
		dst.Errorf(ctx, "Ctd.newLoginAPIParsing->json.Unmarshal() error: %v", err)
		return LoginResult{}, ErrorInvalidResponse
	}

	if strings.ToLower(response.Status) == "success" {
		return LoginResult{AuthKey: response.AuthKey}, nil
	}

	if delay, err := dst.parseLoginErrorStrings(response.Errors.Error); err != nil {
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

func (dst *Ctd) parseLoginErrorStrings(data json.RawMessage) (int, error) {
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
		return 0, ErrorCaptchaRequired
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
