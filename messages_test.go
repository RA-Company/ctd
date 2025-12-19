package ctd

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/ra-company/env"
	"github.com/stretchr/testify/require"
)

func TestCtd_Messages(t *testing.T) {
	ctx := context.Background()
	faker := gofakeit.New(0)

	url, token := getCredentials(t)

	clientID := int64(env.GetEnvInt("API_CLIENT_ID", 0))
	require.NotEqual(t, 0, clientID, "API_CLIENT_ID must be set in .env file and .settings")
	transport := env.GetEnvStr("API_MESSAGE_TRANSPORT", "")
	require.NotEmpty(t, transport, "API_MESSAGE_TRANSPORT must be set in .env file and .settings")
	channelID := int64(env.GetEnvInt("API_CHANNEL_ID", 0))
	require.NotEqual(t, 0, channelID, "API_CHANNEL_ID must be set in .env file and .settings")
	messageURL := env.GetEnvStr("API_MESSAGE_URL", "")
	require.NotEmpty(t, messageURL, "API_MESSAGE_URL must be set in .env file and .settings")

	tests := []struct {
		name    string
		token   string
		message MessagePayload
		isData  bool
		error   error
	}{
		{
			name:   "Incorrect token",
			token:  "incorrect_token",
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name:  "Correct token",
			token: token,
			message: MessagePayload{
				Text:      faker.Sentence(10),
				ChannelID: channelID,
				Transport: transport,
				ClientID:  clientID,
				Type:      faker.RandomString([]string{"to_client", "autoreply"}),
			},
			isData: true,
			error:  nil,
		},
		{
			name:  "Attachment incorrect",
			token: token,
			message: MessagePayload{
				Text:       faker.Sentence(10),
				Attachment: "http://incorrect.url/",
				ChannelID:  channelID,
				Transport:  transport,
				ClientID:   clientID,
				Type:       faker.RandomString([]string{"to_client", "autoreply"}),
			},
			isData: false,
			error:  ErrorInvalidParameters,
		},
		{
			name:  "Attachment correct",
			token: token,
			message: MessagePayload{
				Text:       faker.Sentence(10),
				Attachment: messageURL,
				ChannelID:  channelID,
				Transport:  transport,
				ClientID:   clientID,
				Type:       faker.RandomString([]string{"to_client", "autoreply"}),
			},
			isData: true,
			error:  nil,
		},
		{
			name:  "Correct button",
			token: token,
			message: MessagePayload{
				Text:      faker.Sentence(10),
				ChannelID: channelID,
				Transport: transport,
				ClientID:  clientID,
				Type:      faker.RandomString([]string{"to_client", "autoreply"}),
				Keyboard: &MessageButtons{
					Buttons: []MessageButton{
						{
							Type: "reply",
							Text: faker.Sentence(3),
						},
						{
							Type: "reply",
							Text: faker.Sentence(3),
						},
					},
				},
			},
			isData: true,
			error:  nil,
		},
		{
			name:  "Correct inline button",
			token: token,
			message: MessagePayload{
				Text:      faker.Sentence(10),
				ChannelID: channelID,
				Transport: transport,
				ClientID:  clientID,
				Type:      faker.RandomString([]string{"to_client", "autoreply"}),
				InlineButtons: []MessageButton{
					{
						Type: "reply",
						Text: faker.Sentence(3),
					},
					{
						Type: "reply",
						Text: faker.Sentence(3),
					},
				},
			},
			isData: true,
			error:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(url, tt.token)
			got, err := dst.SendMessage(ctx, &tt.message)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.SendMessage() error")
			}
			if tt.isData {
				require.NotNil(t, got, "dst.SendMessage() should return data")
			} else {
				require.Nil(t, got, "dst.SendMessage() should return nil data on error")
			}
		})
	}
}

func TestCtd_MessagesTransferToGroup(t *testing.T) {
	ctx := context.Background()
	faker := gofakeit.New(0)

	url, token := getCredentials(t)

	clientID := int64(env.GetEnvInt("API_CLIENT_ID", 0))
	require.NotEqual(t, 0, clientID, "API_CLIENT_ID must be set in .env file and .settings")
	transport := env.GetEnvStr("API_MESSAGE_TRANSPORT", "")
	require.NotEmpty(t, transport, "API_MESSAGE_TRANSPORT must be set in .env file and .settings")
	channelID := int64(env.GetEnvInt("API_CHANNEL_ID", 0))
	require.NotEqual(t, 0, channelID, "API_CHANNEL_ID must be set in .env file and .settings")
	operatorID := int64(env.GetEnvInt("API_OPERATOR_ID", 0))
	require.NotEqual(t, 0, operatorID, "API_OPERATOR_ID must be set in .env file and .settings")
	groupID := int64(env.GetEnvInt("API_GROUP_ID", 0))
	require.NotEqual(t, 0, groupID, "API_GROUP_ID must be set in .env file and .settings")

	service := Ctd{}
	service.Init(url, token)

	message, err := service.SendMessage(ctx, &MessagePayload{
		Text:      faker.Sentence(10),
		ChannelID: channelID,
		Transport: transport,
		ClientID:  clientID,
		Type:      "to_client",
	})
	require.NoError(t, err, "service.SendMessage() error")
	require.NotNil(t, message, "service.SendMessage() should return data")
	require.Greater(t, message.MessageID, int64(0), "service.SendMessage() should return valid message ID")
	require.Greater(t, message.DialogID, int64(0), "service.SendMessage() should return valid dialog ID")

	defer func() {
		time.Sleep(1 * time.Second)
		service.CloseDialog(ctx, message.DialogID, operatorID, 0)
	}()

	t.Run("01 Incorrect token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, "incorrect_token")
		err := dst.TransferToGroup(ctx, message.MessageID, groupID, false)
		require.ErrorIs(t, err, ErrorInvalidToken, "dst.TransferDialogToGroup() error")
	})

	t.Run("02 Correct token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		err := dst.TransferToGroup(ctx, message.MessageID, groupID, false)
		require.NoError(t, err, "dst.TransferDialogToGroup() error")
	})

	t.Run("03 Incorrect group ID", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		err := dst.TransferToGroup(ctx, message.MessageID, 0, false)
		require.ErrorIs(t, err, ErrorInvalidOperatorGroupID, "dst.TransferDialogToGroup() error")
	})

	t.Run("04 Incorrect message ID", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		err := dst.TransferToGroup(ctx, 0, groupID, false)
		require.ErrorIs(t, err, ErrorInvalidMesssageID, "dst.TransferDialogToGroup() error")
	})
}

func TestCtd_MessagesTransferToOperator(t *testing.T) {
	ctx := context.Background()
	faker := gofakeit.New(0)

	url, token := getCredentials(t)

	clientID := int64(env.GetEnvInt("API_CLIENT_ID", 0))
	require.NotEqual(t, 0, clientID, "API_CLIENT_ID must be set in .env file and .settings")
	transport := env.GetEnvStr("API_MESSAGE_TRANSPORT", "")
	require.NotEmpty(t, transport, "API_MESSAGE_TRANSPORT must be set in .env file and .settings")
	channelID := int64(env.GetEnvInt("API_CHANNEL_ID", 0))
	require.NotEqual(t, 0, channelID, "API_CHANNEL_ID must be set in .env file and .settings")
	operatorID := int64(env.GetEnvInt("API_OPERATOR_ID", 0))
	require.NotEqual(t, 0, operatorID, "API_OPERATOR_ID must be set in .env file and .settings")
	groupID := int64(env.GetEnvInt("API_GROUP_ID", 0))
	require.NotEqual(t, 0, groupID, "API_GROUP_ID must be set in .env file and .settings")

	service := Ctd{}
	service.Init(url, token)

	message, err := service.SendMessage(ctx, &MessagePayload{
		Text:      faker.Sentence(10),
		ChannelID: channelID,
		Transport: transport,
		ClientID:  clientID,
		Type:      "to_client",
	})
	require.NoError(t, err, "service.SendMessage() error")
	require.NotNil(t, message, "service.SendMessage() should return data")
	require.Greater(t, message.MessageID, int64(0), "service.SendMessage() should return valid message ID")
	require.Greater(t, message.DialogID, int64(0), "service.SendMessage() should return valid dialog ID")

	err = service.TransferToGroup(ctx, message.MessageID, groupID, true)
	require.NoError(t, err, "service.TransferToGroup() error")

	defer func() {
		time.Sleep(1 * time.Second)
		service.CloseDialog(ctx, message.DialogID, operatorID, 0)
	}()

	t.Run("01 Incorrect token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, "incorrect_token")
		err := dst.TransferToOperator(ctx, message.MessageID, operatorID)
		require.ErrorIs(t, err, ErrorInvalidToken, "dst.TransferToOperator() error")
	})

	t.Run("02 Correct token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		err := dst.TransferToOperator(ctx, message.MessageID, operatorID)
		require.NoError(t, err, "dst.TransferToOperator() error")
	})

	t.Run("03 Incorrect operator ID", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		err := dst.TransferToOperator(ctx, message.MessageID, 0)
		require.ErrorIs(t, err, ErrorInvalidOperatorID, "dst.TransferToOperator() error")
	})

	t.Run("04 Incorrect message ID", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		err := dst.TransferToOperator(ctx, 0, operatorID)
		require.ErrorIs(t, err, ErrorInvalidMesssageID, "dst.TransferToOperator() error")
	})
}
