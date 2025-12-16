package ctd

import (
	"context"
	"testing"

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
