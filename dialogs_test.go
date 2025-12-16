package ctd

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/ra-company/env"
	"github.com/stretchr/testify/require"
)

func TestCtd_Dialogs(t *testing.T) {
	ctx := context.Background()
	faker := gofakeit.New(0)

	url, token := getCredentials(t)

	t.Run("01 GetDialogs Incorrect token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, "incorrect_token")
		params := &GetDialogsParams{
			Limit:  10,
			Offset: 0,
			State:  "any",
			Order:  "desc",
		}
		got, total, err := dst.GetDialogs(ctx, params)
		require.ErrorIs(t, err, ErrorInvalidToken, "dst.APIGetDialogs() error")
		require.Nil(t, got, "dst.APIGetDialogs() should return nil data on error")
		require.Equal(t, 0, total, "dst.APIGetDialogs() should return zero total on error")
	})

	var dialog_id = int64(0)
	t.Run("02 GetDialogs Correct token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		params := &GetDialogsParams{
			Limit:  10,
			Offset: 0,
			State:  "any",
			Order:  "desc",
		}
		got, total, err := dst.GetDialogs(ctx, params)
		require.NoError(t, err, "dst.APIGetDialogs() error")
		require.NotNil(t, got, "dst.APIGetDialogs() should return data")
		require.Greater(t, len(got), 0, "dst.APIGetDialogs() should return some dialogs")
		require.Greater(t, total, 0, "dst.APIGetDialogs() should return total dialogs count")
		dialog_id = got[0].ID
	})

	t.Run("03 GetDialog Incorrect ID", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		got, err := dst.GetDialog(ctx, 0)
		require.ErrorIs(t, err, ErrorInvalidID, "dst.APIGetDialog() error")
		require.Nil(t, got, "dst.APIGetDialog() should return nil data on error")
	})

	t.Run("04 GetDialog Correct ID", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		got, err := dst.GetDialog(ctx, dialog_id)
		require.NoError(t, err, "dst.GetDialog() error")
		require.NotNil(t, got, "dst.GetDialog() should return data")
		require.Equal(t, dialog_id, got.ID, "dst.GetDialog() should return correct dialog ID")
	})

	clientID := env.GetEnvInt("API_CLIENT_ID", 0)
	require.NotEqual(t, 0, clientID, "API_CLIENT_ID must be set in .env file and .settings")
	transport := env.GetEnvStr("API_MESSAGE_TRANSPORT", "")
	require.NotEmpty(t, transport, "API_MESSAGE_TRANSPORT must be set in .env file and .settings")
	channelID := env.GetEnvInt("API_CHANNEL_ID", 0)
	require.NotEqual(t, 0, channelID, "API_CHANNEL_ID must be set in .env file and .settings")
	operatorID := env.GetEnvInt("API_OPERATOR_ID", 0)
	require.NotEqual(t, 0, operatorID, "API_OPERATOR_ID must be set in .env file and .settings")

	t.Run("05 CloseDialog Incorrect token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, "incorrect_token")

		err := dst.CloseDialog(ctx, 2934870, 0, 0)
		require.ErrorIs(t, err, ErrorInvalidToken, "dst.CloseDialog() error")
	})

	t.Run("06 CloseDialog correct token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)

		message, err := dst.SendMessage(ctx, &MessagePayload{
			Text:      faker.Sentence(10),
			ChannelID: int64(channelID),
			Transport: transport,
			ClientID:  int64(clientID),
			Type:      "to_client",
		})
		require.NoError(t, err, "dst.SendMessage() error")
		require.NotNil(t, message, "dst.SendMessage() should return data")
		require.Greater(t, message.DialogID, int64(0), "dst.SendMessage() should return valid message DialogID")

		time.Sleep(1 * time.Second)
		err = dst.CloseDialog(ctx, int64(message.DialogID), int64(operatorID), 0)
		require.NoError(t, err, "dst.CloseDialog() error")

		time.Sleep(1 * time.Second)
		err = dst.CloseDialog(ctx, int64(message.DialogID), int64(operatorID), 0)
		require.ErrorIs(t, err, ErrorDialogClosed, "dst.CloseDialog() error")
	})
}
