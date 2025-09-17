package ctd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtd_Webhook(t *testing.T) {
	ctx := context.Background()

	url, token := getCredentials(t)

	var id int

	t.Run("CreateWebhook", func(t *testing.T) {
		ctx := context.Background()

		dst := &Ctd{}
		dst.Init(url, token)

		payload := WebhookPayload{
			Name:     "Test Webhook",
			URL:      "https://localhost/",
			Events:   []string{"outbox", "outbox_status"},
			Channels: []int{},
			Status:   "enable",
		}

		got, err := dst.CreateWebhook(ctx, &payload)
		require.NoError(t, err, "dst.CreateWebhook() should not return an error")
		require.NotNil(t, got, "dst.CreateWebhook() should return data")
		require.Equal(t, payload.Name, got.Name, "dst.CreateWebhook() should return the correct webhook name")
		require.Equal(t, payload.URL, got.URL, "dst.CreateWebhook() should return the correct webhook URL")
		require.Equal(t, payload.Events, got.Events, "dst.CreateWebhook() should return the correct webhook events")
		require.NotEqual(t, 0, got.ID, "dst.CreateWebhook() should return a valid webhook ID")
		require.Equal(t, payload.Status, got.Status, "dst.CreateWebhook() should return the correct webhook status")

		id = got.ID

		t.Run("CreateWebhook with existing URL", func(t *testing.T) {
			got, err := dst.CreateWebhook(ctx, &payload)
			require.ErrorIs(t, err, ErrorWebhookUrlIsAlreadyUsed, "dst.CreateWebhook() should return an error for existing URL")
			require.Nil(t, got, "dst.CreateWebhook() should return nil data on error")
		})
	})

	dst := &Ctd{}
	dst.Init(url, token)
	defer dst.DeleteWebhook(ctx, id)

	tests := []struct {
		name   string
		url    string
		token  string
		isData bool
		error  error
	}{
		{
			name:   "Incorrect token",
			url:    url,
			token:  "incorrect_token",
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name:   "Correct token",
			url:    url,
			token:  token,
			isData: true,
			error:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(tt.url, tt.token)
			got, err := dst.GetWebhooks(ctx)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.GetWebhook() error")
			} else {
				require.NoError(t, err, "dst.GetWebhook() should not return an error")
			}
			if tt.isData {
				require.NotNil(t, got, "dst.GetWebhook() should return data")
			} else {
				require.Nil(t, got, "dst.GetWebhook() should return nil data on error")
			}

		})
	}

	t.Run("Update Webhook", func(t *testing.T) {
		ctx := context.Background()

		dst := &Ctd{}
		dst.Init(url, token)

		payload := WebhookPayload{
			Name:     "Updated Test Webhook",
			URL:      "https://localhost/updated",
			Events:   []string{"inbox"},
			Channels: []int{},
			Status:   "disable",
		}

		got, err := dst.UpdateWebhook(ctx, id, &payload)
		require.NoError(t, err, "dst.UpdateWebhook() should not return an error")
		require.NotNil(t, got, "dst.UpdateWebhook() should return data")
		require.Equal(t, payload.Name, got.Name, "dst.UpdateWebhook() should return the correct webhook name")
		require.Equal(t, payload.URL, got.URL, "dst.UpdateWebhook() should return the correct webhook URL")
		require.Equal(t, payload.Events, got.Events, "dst.UpdateWebhook() should return the correct webhook events")
		require.Equal(t, payload.Status, got.Status, "dst.UpdateWebhook() should return the correct webhook status")
		require.Equal(t, id, got.ID, "dst.UpdateWebhook() should return the correct webhook ID")
	})

	t.Run("Delete Webhook", func(t *testing.T) {
		ctx := context.Background()

		dst := &Ctd{}
		dst.Init(url, token)

		err := dst.DeleteWebhook(ctx, id)
		require.NoError(t, err, "dst.DeleteWebhooks() should not return an error")

		t.Run("Delete non-existing Webhook", func(t *testing.T) {
			err := dst.DeleteWebhook(ctx, id)
			require.ErrorIs(t, err, ErrorInvalidID, "dst.DeleteWebhooks() should return an error for non-existing webhook")
		})
	})
}
