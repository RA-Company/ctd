package ctd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// newTestCtd returns a Ctd connected to a test server that answers every request with the given handler.
func newTestCtd(t *testing.T, handler http.HandlerFunc) *Ctd {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	ctd := &Ctd{}
	ctd.Init(srv.URL, "token")

	return ctd
}

func respond(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}
}

func Test_tooFastDelay(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want int
	}{
		{"seconds", "please, try again after 42 seconds", 42},
		{"one second", "please, try again after 1 second", 1},
		{"no after", "something else", 0},
		{"no second", "please, try again after soon", 0},
		{"not a number", "please, try again after abc seconds", 0},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tooFastDelay(tt.str))
		})
	}
}

func TestCtd_Login_Unit(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    LoginResult
		wantErr error
	}{
		{
			name: "success returns auth key",
			body: `{"status":"success","auth_key":"K.user"}`,
			want: LoginResult{AuthKey: "K.user"},
		},
		{
			name:    "new API timeout",
			body:    `{"status":"error","errors":{"error":["timeout"]},"login_attempts_info":{"a":1}}`,
			want:    LoginResult{RetryDelay: 10},
			wantErr: ErrorTooFast,
		},
		{
			name:    "old API too fast",
			body:    `{"status":"error","errors":{"brute_force":["Please, try again after 42 seconds"]},"login_attempts_info":{"a":1}}`,
			want:    LoginResult{RetryDelay: 42},
			wantErr: ErrorTooFast,
		},
		{
			name:    "too fast with unexpected text does not panic",
			body:    `{"status":"error","errors":{"brute_force":["Please, try again after soon"]},"login_attempts_info":{"a":1}}`,
			want:    LoginResult{},
			wantErr: ErrorTooFast,
		},
		{
			name:    "too fast in plain text does not panic",
			body:    `{"status":"error","errors":{"x":["please, try again after"]},"login_attempts_info":{"a":1}}`,
			want:    LoginResult{},
			wantErr: ErrorTooFast,
		},
		{
			name:    "OTP required",
			body:    `{"status":"error","errors":{"error":["incorrect_otp"]},"login_attempts_info":{"a":1}}`,
			wantErr: ErrorOTPRequired,
		},
		{
			name:    "unknown response returns raw body",
			body:    `{"status":"error","errors":{"foo":["bar"]},"login_attempts_info":{"a":1}}`,
			want:    LoginResult{Error: `{"status":"error","errors":{"foo":["bar"]},"login_attempts_info":{"a":1}}`},
			wantErr: ErrorUnknownError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctd := newTestCtd(t, respond(tt.body))

			got, err := ctd.Login(t.Context(), "login", "password", "", "", "", "")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestClient_UnmarshalJSON_ExternalID(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"string", `{"id":1,"external_id":"ab"}`, "ab"},
		{"number", `{"id":1,"external_id":5}`, "5"},
		{"big number", `{"id":1,"external_id":9007199254740993}`, "9007199254740993"},
		{"null", `{"id":1,"external_id":null}`, ""},
		{"missing", `{"id":1}`, ""},
		{"empty string", `{"id":1,"external_id":""}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client Client
			require.NoError(t, json.Unmarshal([]byte(tt.json), &client))
			require.Equal(t, tt.want, client.ExtrnalID)
			require.Equal(t, 1, client.ID)
		})
	}

	t.Run("nested in response", func(t *testing.T) {
		var response ClientResponse
		require.NoError(t, json.Unmarshal([]byte(`{"data":{"id":2,"name":"n","external_id":7},"status":"success"}`), &response))
		require.Equal(t, "7", response.Data.ExtrnalID)
		require.Equal(t, "n", response.Data.Name)
	})

	t.Run("in list", func(t *testing.T) {
		var response ClientsResponse
		require.NoError(t, json.Unmarshal([]byte(`{"data":[{"id":1,"external_id":3},{"id":2,"external_id":"x"}],"status":"success"}`), &response))
		require.Equal(t, "3", response.Data[0].ExtrnalID)
		require.Equal(t, "x", response.Data[1].ExtrnalID)
	})

	t.Run("invalid type", func(t *testing.T) {
		var client Client
		require.Error(t, json.Unmarshal([]byte(`{"id":1,"external_id":{"a":1}}`), &client))
	})
}

func TestCtd_createClientExists(t *testing.T) {
	ctd := &Ctd{}

	tests := []struct {
		name    string
		errors  string
		want    int64
		wantErr bool
	}{
		{"valid", `{"client":["client already exists","{\"id\":123}"]}`, 123, false},
		{"one element does not panic", `{"client":["client already exists"]}`, 0, true},
		{"missing key does not panic", `{"other":["x"]}`, 0, true},
		{"not json", `not json`, 0, true},
		{"inner not json", `{"client":["a","b"]}`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctd.createClientExists(t.Context(), json.RawMessage(tt.errors))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCtd_CreateClient_AlreadyExists(t *testing.T) {
	ctd := newTestCtd(t, respond(`{"status":"error","errors":{"client":["Client already exists","{\"id\":55}"]}}`))

	client, err := ctd.CreateClient(t.Context(), "+100", "whatsapp", 1, "nick", "")
	require.ErrorIs(t, err, ErrorClieantAlreadyExists)
	require.NotNil(t, client)
	require.Equal(t, 55, client.ID)
}

// blockFirst makes the first request hang until the client gives up, the next ones are answered immediately.
func blockFirst(calls *atomic.Int32) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The body must be read, otherwise the server does not notice that the client has gone.
		_, _ = io.Copy(io.Discard, r.Body)
		if calls.Add(1) == 1 {
			<-r.Context().Done()
			return
		}
		_, _ = w.Write([]byte(`{}`))
	}
}

func TestCtd_Retry(t *testing.T) {
	tests := []struct {
		name      string
		retry     bool
		wantCalls int32
		wantErr   bool
	}{
		{"disabled by default", false, 1, true},
		{"enabled", true, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls atomic.Int32
			ctd := newTestCtd(t, blockFirst(&calls))
			ctd.Timeout = 1
			ctd.Retry = tt.retry

			var out map[string]any
			_, err := ctd.Post(t.Context(), "path", `{}`, &out)

			require.Equal(t, tt.wantCalls, calls.Load())
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	t.Run("default is disabled", func(t *testing.T) {
		ctd := &Ctd{}
		ctd.Init("http://localhost", "")
		require.False(t, ctd.Retry)
	})
}

func TestCtd_Retry_NotOnCancelledContext(t *testing.T) {
	var calls atomic.Int32
	ctd := newTestCtd(t, blockFirst(&calls))
	ctd.Retry = true

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	_, err := ctd.Get(ctx, "path", nil)
	require.Error(t, err)
	require.Equal(t, int32(1), calls.Load(), "request must not be repeated after the context is done")
}

func TestCtd_ContextCancellation(t *testing.T) {
	ctd := newTestCtd(t, func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	})

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := ctd.Get(ctx, "path", nil)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(start), 5*time.Second, "request must be aborted by the context, not by the client timeout")
}

func TestCtd_Init(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://api.test", "https://api.test/"},
		{"https://api.test/", "https://api.test/"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			ctd := &Ctd{}
			require.NotPanics(t, func() { ctd.Init(tt.url, "token") })
			require.Equal(t, tt.want, ctd.Url)
			require.Equal(t, "token", ctd.Token)
			require.Equal(t, uint(10), ctd.Timeout)
		})
	}
}

func TestAPIErrors_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"string", `{"errors":"client not found"}`, "client not found"},
		{"object", `{"errors":{"client":["not found"]}}`, `{"client":["not found"]}`},
		{"array", `{"errors":["a","b"]}`, `["a","b"]`},
		{"null", `{"errors":null}`, ""},
		{"missing", `{}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tags TagResponse
			require.NoError(t, json.Unmarshal([]byte(tt.json), &tags))
			require.Equal(t, tt.want, tags.Errors.String())

			var clients ClientsResponse
			require.NoError(t, json.Unmarshal([]byte(tt.json), &clients))
			require.Equal(t, tt.want, clients.Errors.String())
		})
	}
}

func TestCtd_GetTag_ObjectErrors(t *testing.T) {
	ctd := newTestCtd(t, respond(`{"status":"error","errors":{"tag":["tag not found"]}}`))

	_, err := ctd.GetTag(t.Context(), 1)
	require.ErrorIs(t, err, ErrorInvalidID)
}

func TestCtd_HTTPStatus(t *testing.T) {
	t.Run("server error with HTML body", func(t *testing.T) {
		ctd := newTestCtd(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`<html>unavailable</html>`))
		})

		_, err := ctd.GetClient(t.Context(), 1)
		require.ErrorIs(t, err, ErrorInvalidResponse)
		require.ErrorContains(t, err, "503")
	})

	t.Run("server error without expected response", func(t *testing.T) {
		ctd := newTestCtd(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		})

		_, err := ctd.Login(t.Context(), "login", "password", "", "", "", "")
		require.ErrorIs(t, err, ErrorInvalidResponse)
		require.ErrorContains(t, err, "502")
	})

	t.Run("client error with JSON body is not an error", func(t *testing.T) {
		ctd := newTestCtd(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"status":"error","errors":"client not found"}`))
		})

		_, err := ctd.GetClient(t.Context(), 1)
		require.ErrorIs(t, err, ErrorInvalidID)
	})

	t.Run("invalid JSON with 200", func(t *testing.T) {
		ctd := newTestCtd(t, respond(`not json`))

		_, err := ctd.GetClient(t.Context(), 1)
		require.ErrorIs(t, err, ErrorInvalidResponse)
		require.NotContains(t, err.Error(), "HTTP")
	})
}

func TestCtd_Post_MarshalError(t *testing.T) {
	var calls atomic.Int32
	ctd := newTestCtd(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
	})

	_, err := ctd.Post(t.Context(), "path", make(chan int), nil)
	require.Error(t, err)
	require.Equal(t, int32(0), calls.Load(), "request must not be sent")
}

func TestClient_UnmarshalJSON_ExternalIDs(t *testing.T) {
	var client Client
	require.NoError(t, json.Unmarshal([]byte(`{"id":1,"external_ids":{"a":1,"b":"x","c":null}}`), &client))
	require.Equal(t, map[string]string{"a": "1", "b": "x", "c": ""}, client.ExtrnalIDs)

	require.NoError(t, json.Unmarshal([]byte(`{"id":1,"external_ids":{}}`), &client))
	require.NotNil(t, client.ExtrnalIDs)
	require.Empty(t, client.ExtrnalIDs)

	require.NoError(t, json.Unmarshal([]byte(`{"id":1,"external_ids":null}`), &client))
	require.Nil(t, client.ExtrnalIDs)

	require.Error(t, json.Unmarshal([]byte(`{"id":1,"external_ids":{"a":[1]}}`), &client))
}

func TestCtd_APIGetClients_Params(t *testing.T) {
	tests := []struct {
		name    string
		params  string
		want    string
		wantErr error
	}{
		{"empty", "", "limit=10&offset=0&order=asc", nil},
		{"simple", "phone=123", "limit=10&offset=0&order=asc&phone=123", nil},
		{"leading separators", "?&phone=123", "limit=10&offset=0&order=asc&phone=123", nil},
		{"unescaped space", "name=John Doe", "limit=10&name=John+Doe&offset=0&order=asc", nil},
		{"escaped plus sign in phone", "phone=%2B7999", "limit=10&offset=0&order=asc&phone=%2B7999", nil},
		{"invalid escape", "phone=%zz", "", ErrorInvalidParameters},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var query string
			ctd := newTestCtd(t, func(w http.ResponseWriter, r *http.Request) {
				query = r.URL.Query().Encode()
				_, _ = w.Write([]byte(`{"status":"success","data":[]}`))
			})

			_, err := ctd.APIGetClients(t.Context(), 0, 10, "asc", tt.params)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, query)
		})
	}
}

func TestCtd_GetTags_ErrorStatus(t *testing.T) {
	ctd := newTestCtd(t, respond(`{"status":"error","errors":"forbidden"}`))

	tags, total, err := ctd.GetTags(t.Context(), 0, 10)
	require.ErrorIs(t, err, ErrorInvalidResponse)
	require.Nil(t, tags)
	require.Zero(t, total)
}

func TestCtd_CompaniesApiInfo_ErrorStatus(t *testing.T) {
	ctd := newTestCtd(t, respond(`{"status":"error","errors":"forbidden"}`))

	info, err := ctd.CompaniesApiInfo(t.Context())
	require.ErrorIs(t, err, ErrorInvalidResponse)
	require.Nil(t, info)

	ctd = newTestCtd(t, respond(`{"status":"success","data":{"companyID":5,"company_name":"c"}}`))
	info, err = ctd.CompaniesApiInfo(t.Context())
	require.NoError(t, err)
	require.Equal(t, int64(5), info.CompanyID)
}

func TestCtd_CreateClient_EmailUnit(t *testing.T) {
	var payload map[string]any
	ctd := newTestCtd(t, func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		_, _ = w.Write([]byte(`{"status":"success","data":{"id":77,"phone":"user@example.com","assigned_name":"John"}}`))
	})

	client, err := ctd.CreateClient(t.Context(), "user@example.com", "email", 5, "John", "")
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, 77, client.ID)
	require.Equal(t, "user@example.com", client.Phone)

	require.Equal(t, "user@example.com", payload["phone"], "email is sent in the phone field")
	require.Equal(t, "email", payload["transport"])
	require.EqualValues(t, 5, payload["channel_id"])
	require.Equal(t, "John", payload["nickname"])
	require.NotContains(t, payload, "assigned_phone", "empty assigned phone must not be sent")
}

func TestCtd_CreateClient_EmailAlreadyExistsUnit(t *testing.T) {
	ctd := newTestCtd(t, respond(`{"status":"error","errors":{"client":["Client already exists","{\"id\":88}"]}}`))

	client, err := ctd.CreateClient(t.Context(), "user@example.com", "email", 5, "John", "")
	require.ErrorIs(t, err, ErrorClieantAlreadyExists)
	require.NotNil(t, client)
	require.Equal(t, 88, client.ID)
	require.Equal(t, "user@example.com", client.Phone)
}
