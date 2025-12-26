package ctd

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/ra-company/env"
	"github.com/stretchr/testify/require"
)

func TestCtd_Clients(t *testing.T) {
	ctx := context.Background()
	faker := gofakeit.New(0)

	url, token := getCredentials(t)

	clientID := env.GetEnvInt("API_CLIENT_ID", 0)
	require.NotEqual(t, 0, clientID, "API_CLIENT_ID must be set in .env file or .settings")

	tests := []struct {
		name      string
		token     string
		client_id int
		isData    bool
		error     error
	}{
		{
			name:   "Incorrect token",
			token:  "incorrect_token",
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name:      "Correct token",
			token:     token,
			client_id: clientID,
			isData:    true,
			error:     nil,
		},
		{
			name:      "Client not found",
			token:     token,
			client_id: 0, // Assuming 0 is an invalid ID for testing
			isData:    false,
			error:     ErrorInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(url, tt.token)
			got, err := dst.GetClient(ctx, tt.client_id)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.GetClient() error")

			} else {

			}
			if tt.isData {
				require.NotNil(t, got, "dst.GetClient() should return data")
			} else {
				require.Nil(t, got, "dst.GetClient() should return nil data on error")
			}

		})
	}

	t.Run("GetClientsList", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)

		offset := 0
		limit := 100

		got, total, err := dst.GetClientsList(ctx, offset, limit)
		require.NoError(t, err, "dst.GetClientsList() should not return an error")
		require.NotNil(t, got, "dst.GetClientsList() should return data")
		require.Greater(t, len(got), 0, "dst.GetClientsList() should return non-empty client list")
		require.Greater(t, total, 0, "dst.GetClientsList() should return non-zero total count")

		offset += limit
		got2, total2, err := dst.GetClientsList(ctx, offset, limit)
		require.NoError(t, err, "dst.GetClientsList() should not return an error on second call")
		require.NotNil(t, got2, "dst.GetClientsList() should return data on second call")
		require.Greater(t, len(got2), 0, "dst.GetClientsList() should return non-empty client list on second call")
		require.Greater(t, total2, 0, "dst.GetClientsList() should return non-zero total count on second call")
		require.Equal(t, total, total2, "dst.GetClientsList() should return the same total count on both calls")
		require.NotEqual(t, got, got2, "dst.GetClientsList() should return different data on different calls with different offsets")
	})

	type createClient struct {
		Phone         string `json:"phone"`
		Transport     string `json:"transport"`
		ChannelID     int    `json:"channel_id"`
		Nickname      string `json:"nickname"`
		AssignedPhone string `json:"assigned_phone"`
	}

	transport := env.GetEnvStr("API_TRANSPORT", "")
	if transport == "" {
		t.Skip("API_TRANSPORT is not set, skipping create client tests")
	}

	channel := env.GetEnvInt("API_CHANNEL_ID", 0)
	if channel == 0 {
		t.Skip("API_CHANNEL_ID is not set, skipping create client tests")
	}

	createTests := []struct {
		name   string
		token  string
		client createClient
		isData bool
		error  error
	}{
		{
			name:  "Create client with incorrect token",
			token: "incorrect_token",
			client: createClient{
				Phone:         fmt.Sprintf("%d", faker.IntRange(10000000, 20000000)),
				Transport:     transport,
				ChannelID:     channel,
				Nickname:      faker.Name(),
				AssignedPhone: faker.Phone(),
			},
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name:  "Create client with correct token",
			token: token,
			client: createClient{
				Phone:         fmt.Sprintf("%d", faker.IntRange(10000000, 20000000)),
				Transport:     transport,
				ChannelID:     channel,
				Nickname:      faker.Name(),
				AssignedPhone: faker.Phone(),
			},
			isData: true,
			error:  nil,
		},
	}

	for _, tt := range createTests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(url, tt.token)

			got, err := dst.CreateClient(ctx, tt.client.Phone, tt.client.Transport, tt.client.ChannelID, tt.client.Nickname, tt.client.AssignedPhone)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.CreateClient() error")
			} else {
				require.NoError(t, err, "dst.CreateClient() should not return an error")
			}
			if tt.isData {
				require.NotNil(t, got, "dst.CreateClient() should return data")
				require.Equal(t, tt.client.Phone, got.Phone, "dst.CreateClient() should return client with correct phone")
				require.Equal(t, tt.client.Nickname, got.AssignedName, "dst.CreateClient() should return client with correct nickname")
				require.Equal(t, tt.client.AssignedPhone, got.ClientPhone, "dst.CreateClient() should return client with correct assigned phone")
			} else {
				require.Nil(t, got, "dst.CreateClient() should return nil data on error")
			}
		})
	}
}

func TestCtd_CreateClient(t *testing.T) {
	ctx := t.Context()
	url, token := getCredentials(t)
	faker := gofakeit.New(0)

	transport := env.GetEnvStr("API_CREATE_CLIENT_TRANSPORT", "")
	if transport == "" {
		t.Skip("API_CREATE_CLIENT_TRANSPORT is not set, skipping create client tests")
	}

	channel := env.GetEnvInt("API_CREATE_CLIENT_CHANNEL_ID", 0)
	if channel == 0 {
		t.Skip("API_CREATE_CLIENT_CHANNEL_ID is not set, skipping create client tests")
	}

	phone := env.GetEnvStr("API_CREATE_CLIENT_PHONE", "")
	if phone == "" {
		t.Skip("API_CREATE_CLIENT_PHONE is not set, skipping create client tests")
	}

	id := env.GetEnvInt("API_CREATE_CLIENT_ID", 0)
	if id == 0 {
		t.Skip("API_CREATE_CLIENT_ID is not set, skipping create client tests")
	}

	service := &Ctd{}
	service.Init(url, token)

	t.Run("Create existing client", func(t *testing.T) {
		got, err := service.CreateClient(ctx, phone, transport, channel, faker.Name(), phone)
		require.ErrorIs(t, err, ErrorClieantAlreadyExists, "service.CreateClient() should return ErrorClieantAlreadyExists error")
		require.NotNil(t, got, "service.CreateClient() should return data")
		require.Equal(t, phone, got.Phone, "service.CreateClient() should return client with correct phone")
		require.Equal(t, id, got.ID, "service.CreateClient() should return client with correct ID")
	})

	t.Run("Incorrect channel", func(t *testing.T) {
		got, err := service.CreateClient(ctx, phone, transport, 0, faker.Name(), phone)
		require.ErrorIs(t, err, ErrorInvalidChannelID, "service.CreateClient() should return ErrorInvalidChannelID error")
		require.Nil(t, got, "service.CreateClient() should return data")
	})

	t.Run("Incorrect transport", func(t *testing.T) {
		got, err := service.CreateClient(ctx, phone, "invalid_transport", channel, faker.Name(), phone)
		require.ErrorIs(t, err, ErrorInvalidTransport, "service.CreateClient() should return ErrorInvalidTransport error")
		require.Nil(t, got, "service.CreateClient() should return data")
	})
}
