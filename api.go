package ctd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ra-company/logging"
)

var (
	ErrorInvalidResponse         = fmt.Errorf("invalid response")
	ErrorInvalidToken            = fmt.Errorf("invalid token")
	ErrorWebhookUrlIsAlreadyUsed = fmt.Errorf("webhook URL is already used")
	ErrorInvalidID               = fmt.Errorf("invalid ID")
	ErrorInvalidParameters       = fmt.Errorf("invalid parameters")
	ErrorUnknownError            = fmt.Errorf("unknown error")
	ErrorDialogClosed            = fmt.Errorf("dialog is closed")
	ErrorInvalidRequestID        = fmt.Errorf("invalid request ID")
	ErrorInvalidClientID         = fmt.Errorf("invalid client ID")
	ErrorInvalidTagID            = fmt.Errorf("invalid tag ID")
	ErrorInvalidOperatorGroupID  = fmt.Errorf("invalid operator group ID")
	ErrorInvalidOperatorID       = fmt.Errorf("invalid operator ID")
	ErrorInvalidMesssageID       = fmt.Errorf("invalid message ID")
)

// MetaResponse provides metadata about the response,
// including total count, limit, and offset.
// It is used to provide pagination information.
// This struct is typically included in API responses to indicate
// how many items are available, how many items are returned,
// and the offset for the next set of items.
type MetaResponse struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type BasicResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}
type Ctd struct {
	logging.CustomLogger
	Url       string
	Token     string
	Timeout   uint
	lastError any // Last error encountered during API requests
}

// Init initializes the Ctd instance with the provided URL and token.
// It sets the URL to ensure it ends with a slash and assigns the token.
// The timeout is set to 10 seconds by default.
// This method is typically called before making any API requests to ensure
// that the Ctd instance is properly configured with the necessary
// URL and authentication token.
//
// Parameters:
//   - url: The base URL of the Chat2Desk API, which should end with a slash.
//   - token: The authentication token for the Chat2Desk API, used
func (dst *Ctd) Init(url string, token string) {
	if url[len(url)-1:] != "/" {
		dst.Url = url + "/"
	} else {
		dst.Url = url
	}

	dst.Token = token
	dst.Timeout = 10
}

// Get retrieves data from the specified path using a GET request.
// It constructs the full URL by appending the path to the base URL.
// The method sends a GET request to the API and returns the response data as a byte slice.
// If an error occurs during the request, it logs the error and returns it.
// If the request times out, it retries the request once.
// This method is typically used to fetch data from the Chat2Desk API.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - path: The path to the specific API endpoint to retrieve data from.
//   - response: A pointer to a struct where the response data will be unmarshaled.
//
// Returns:
//   - A byte slice containing the response data from the API.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) Get(ctx context.Context, path string, response any) ([]byte, error) {
	url := dst.Url + path

	result, err := dst.doRequest(ctx, "GET", url, nil, response)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			result, err = dst.doRequest(ctx, "GET", url, nil, response)
		}
	}
	return result, err
}

// Post sends data to the specified path using a POST request.
// It constructs the full URL by appending the path to the base URL.
// The method sends a POST request to the API with the provided data and returns the response data as a byte slice.
// If an error occurs during the request, it logs the error and returns it.
// If the request times out, it retries the request once.
// This method is typically used to send data to the Chat2Desk API.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - path: The path to the specific API endpoint to send data to.
//   - data: The data to be sent in the request body, which can be of any type (string, byte slice, or struct).
//   - response: A pointer to a struct where the response data will be unmarshaled.
//
// Returns:
//   - A byte slice containing the response data from the API.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) Post(ctx context.Context, path string, data any, response any) ([]byte, error) {
	url := dst.Url + path

	result, err := dst.doRequest(ctx, "POST", url, data, response)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			result, err = dst.doRequest(ctx, "POST", url, data, response)
		}
	}
	return result, err
}

// Put sends data to the specified path using a PUT request.
// It constructs the full URL by appending the path to the base URL.
// The method sends a PUT request to the API with the provided data and returns the response data as a byte slice.
// If an error occurs during the request, it logs the error and returns it.
// If the request times out, it retries the request once.
// This method is typically used to update data in the Chat2Desk API.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - path: The path to the specific API endpoint to send data to.
//   - data: The data to be sent in the request body, which can be of any type (string, byte slice, or struct).
//   - response: A pointer to a struct where the response data will be unmarshaled.
//
// Returns:
//   - A byte slice containing the response data from the API.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) Put(ctx context.Context, path string, data any, response any) ([]byte, error) {
	url := dst.Url + path

	result, err := dst.doRequest(ctx, "PUT", url, data, response)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			result, err = dst.doRequest(ctx, "PUT", url, data, response)
		}
	}
	return result, err
}

// Delete sends a DELETE request to the specified path.
// It constructs the full URL by appending the path to the base URL.
// The method sends a DELETE request to the API and returns the response data as a byte slice.
// If an error occurs during the request, it logs the error and returns it.
// If the request times out, it retries the request once.
// This method is typically used to delete data from the Chat2Desk API.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - path: The path to the specific API endpoint to delete data from.
//   - response: A pointer to a struct where the response data will be unmarshaled.
//
// Returns:
//   - A byte slice containing the response data from the API.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) Delete(ctx context.Context, path string, response any) ([]byte, error) {
	url := dst.Url + path

	result, err := dst.doRequest(ctx, "DELETE", url, nil, response)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			result, err = dst.doRequest(ctx, "DELETE", url, nil, response)
		}
	}
	return result, err
}

// doRequest performs an HTTP request with the specified method, URL, and payload.
// It handles the request creation, sending, and response reading.
// The method supports GET, POST, PUT, and DELETE requests.
// It sets the appropriate headers, including the Authorization header if a token is provided.
// It also measures the time taken for the request and logs debug information.
// If the response body contains an error message indicating an invalid token,
// it returns an ErrorInvalidToken error.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - method: The HTTP method to use for the request (e.g., "GET", "POST", "PUT", "DELETE").
//   - url: The full URL for the request, including the base URL and any specific path.
//   - payload: The data to be sent in the request body, which can be of any type (string, byte slice, or struct).
//
// Returns:
//   - A byte slice containing the response data from the API.
//   - An error if the request fails, if the response is invalid, or if the response indicates an invalid token.
func (dst *Ctd) doRequest(ctx context.Context, method string, url string, payload any, response any) ([]byte, error) {
	start := time.Now()
	client := &http.Client{
		Timeout: time.Duration(dst.Timeout) * time.Second,
	}

	var req *http.Request
	var err error
	if payload == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		var data []byte
		switch v := payload.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			data, _ = json.Marshal(v)
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	}
	if err != nil {
		dst.Error(ctx, fmt.Sprintf("%v", err))
		return nil, err
	}

	if dst.Token != "" {
		req.Header.Set("Authorization", dst.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	dst.Debug(ctx, fmt.Sprintf("\033[1m\033[36mAPI %s (%.2f ms)\033[1m \033[35m%s\033[0m", method, float64(time.Since(start))/1000000, url))
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(body), "Token is not correct") {
		return nil, ErrorInvalidToken
	}

	if response != nil {
		err = json.Unmarshal(body, response)
		if err != nil {
			dst.Error(ctx, "Failed to unmarshal response (%s): %v", body, err)
			return body, ErrorInvalidResponse
		}
	}

	return body, nil
}

// LastError returns the last error encountered during API requests.
// This method is useful for retrieving the last error that occurred,
// allowing for error handling or logging in the application.
//
// Returns:
//   - The last error encountered during API requests, or nil if no error occurred.
func (dst *Ctd) LastError() any {
	return dst.lastError
}
