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
	ErrorInvalidResponse = fmt.Errorf("invalid response")
	ErrorInvalidToken    = fmt.Errorf("invalid token")
)

type MetaResponse struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Ctd struct {
	Url     string
	Token   string
	Timeout uint
}

// Initialize Chat2Desk API
func (dst *Ctd) Init(url string, token string) {
	if url[len(url)-1:] != "/" {
		dst.Url = url + "/"
	} else {
		dst.Url = url
	}

	dst.Token = token
	dst.Timeout = 10
}

func (dst *Ctd) Get(ctx context.Context, path string) ([]byte, error) {
	url := dst.Url + path

	result, err := dst.doRequest(ctx, "GET", url, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			result, err = dst.doRequest(ctx, "GET", url, nil)
		}
	}
	return result, err
}

func (dst *Ctd) Post(ctx context.Context, path string, data interface{}) ([]byte, error) {
	url := dst.Url + path

	result, err := dst.doRequest(ctx, "POST", url, data)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			result, err = dst.doRequest(ctx, "POST", url, data)
		}
	}
	return result, err
}

func (dst *Ctd) Put(ctx context.Context, path string, data interface{}) ([]byte, error) {
	url := dst.Url + path

	result, err := dst.doRequest(ctx, "PUT", url, data)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			result, err = dst.doRequest(ctx, "PUT", url, data)
		}
	}
	return result, err
}

func (dst *Ctd) Delete(ctx context.Context, path string) ([]byte, error) {
	url := dst.Url + path

	result, err := dst.doRequest(ctx, "DELETE", url, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			result, err = dst.doRequest(ctx, "DELETE", url, nil)
		}
	}
	return result, err
}

func (dst *Ctd) doRequest(ctx context.Context, method string, url string, payload interface{}) ([]byte, error) {
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
		logging.Logs.Errorf(ctx, "%v", err)
		return nil, err
	}

	if dst.Token != "" {
		req.Header.Set("Authorization", dst.Token)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	logging.Logs.Debugf(ctx, "\033[1m\033[36mAPI %s (%.2f ms)\033[1m \033[35m%s\033[0m", method, float64(time.Since(start))/1000000, url)
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

	return body, nil
}
