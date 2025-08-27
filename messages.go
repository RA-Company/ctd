package ctd

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

type MessageButton struct {
	Type  string `json:"type,omitempty"`  // Type: button type ('reply', 'location', 'phone', 'email', 'url')
	Text  string `json:"text,omitempty"`  // Text: button label
	Color string `json:"color,omitempty"` // Color: button color ('red', 'green', 'blue', 'yellow') (for online chat only)
	Url   string `json:"url,omitempty"`   // Url: button URL (for 'url' type)
}

type MessageButtons struct {
	Buttons []MessageButton `json:"buttons,omitempty"` // Buttons: optional array of keyboard buttons
}

type MessagePayload struct {
	Text               string          `json:"text"`                          // Text: message text
	Attachment         string          `json:"attachment,omitempty"`          // Attachment: optional message attachment
	AttachmentFilename string          `json:"attachment_filename,omitempty"` // AttachmentFilename: optional attachment filename
	Type               string          `json:"type,omitempty"`                // Type: optional message type ('to_client' (default), 'autoreply', 'system', 'comment')
	ClientID           int             `json:"client_id,omitempty"`           // ClientID: optional client ID to send message to specific client
	ChannelID          int             `json:"channel_id,omitempty"`          // ChannelID: optional channel ID to send message to specific channel
	OperatorID         int             `json:"operator_id,omitempty"`         // OperatorID: optional operator ID to send message as specific operator
	Transport          string          `json:"transport,omitempty"`           // Transport: optional transport to send message via specific transport
	ExternalID         string          `json:"external_id,omitempty"`         // ExternalID: optional external ID to associate with the message
	ReplyMessageID     int             `json:"reply_message_id,omitempty"`    // ReplyMessageID: optional ID of the message being replied to
	InlineButtons      []MessageButton `json:"inline_buttons,omitempty"`      // InlineButtons: optional array of inline buttons
	Keyboard           *MessageButtons `json:"keyboard,omitempty"`            // Keyboard: optional array of keyboard buttons
	Interactive        string          `json:"interactive,omitempty"`         // Interactive: optional interactive parameters for the 'list' and 'button' types. Only for wa_dialog.
}

type MessageResponse struct {
	Data    MessageItem `json:"data"` // Data: List of clients
	Message string      `json:"message"`
	Errors  any         `json:"errors,omitempty"` // Errors: List of errors,
	Status  string      `json:"status"`
}

type MessageItem struct {
	MessageID  int    `json:"message_id"`  // MessageID: Unique message ID
	ChannelID  int    `json:"channel_id"`  // ChannelID: Channel ID
	OperatorID int    `json:"operator_id"` // OperatorID: Operator ID
	Transport  string `json:"transport"`   // Transport: Transport
	Type       string `json:"type"`        // Type: Message type ('to_client', 'autoreply', 'system', 'comment')
	ClientID   int    `json:"client_id"`   // ClientID: Client ID
	DialogID   int    `json:"dialog_id"`   // DialogID: Dialog ID
	RequestID  int    `json:"request_id"`  // RequestID: Request ID
}

// APISendMessage sends a message via the API.
// It takes a context and a MessagePayload, and returns a MessageResponse or an error.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - message (*MessagePayload): The message payload to send.
//
// Returns:
//   - A pointer to a MessageResponse containing the response data.
//   - An error if the request fails.
func (dst *Ctd) APISendMessage(ctx context.Context, message *MessagePayload) (*MessageResponse, error) {
	url := fmt.Sprintf("%sv1/messages", dst.Url)
	response := MessageResponse{}

	message.Type = strings.ToLower(message.Type)
	if !slices.Contains([]string{"to_client", "autoreply", "system", "comment"}, message.Type) {
		message.Type = "to_client"
	}

	if _, err := dst.doRequest(ctx, "POST", url, message, &response); err != nil {
		dst.Error(ctx, "Failed send message: %v", err)
		return nil, err
	}
	return &response, nil
}

// SendMessage sends a message to the API.
// It takes a context and a MessagePayload, and returns a MessageItem or an error.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - message (MessagePayload): The message payload to send.
//
// Returns:
//   - A pointer to a MessageItem containing the response data.
//   - An error if the request fails.
func (dst *Ctd) SendMessage(ctx context.Context, message *MessagePayload) (*MessageItem, error) {
	data, err := dst.APISendMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	if data.Status != "success" {
		dst.Error(ctx, "Failed to send message: %s", data.Errors)
		return nil, ErrorInvalidParameters
	}

	return &data.Data, nil
}
