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
	ClientID           int64           `json:"client_id,omitempty"`           // ClientID: optional client ID to send message to specific client
	ChannelID          int64           `json:"channel_id,omitempty"`          // ChannelID: optional channel ID to send message to specific channel
	OperatorID         int64           `json:"operator_id,omitempty"`         // OperatorID: optional operator ID to send message as specific operator
	Transport          string          `json:"transport,omitempty"`           // Transport: optional transport to send message via specific transport
	ExternalID         string          `json:"external_id,omitempty"`         // ExternalID: optional external ID to associate with the message
	ReplyMessageID     int64           `json:"reply_message_id,omitempty"`    // ReplyMessageID: optional ID of the message being replied to
	InlineButtons      []MessageButton `json:"inline_buttons,omitempty"`      // InlineButtons: optional array of inline buttons
	Keyboard           *MessageButtons `json:"keyboard,omitempty"`            // Keyboard: optional array of keyboard buttons
	Interactive        string          `json:"interactive,omitempty"`         // Interactive: optional interactive parameters for the 'list' and 'button' types. Only for wa_dialog.
}

type SendMessageResponse struct {
	Data    SendMessage `json:"data"` // Data: List of clients
	Message string      `json:"message"`
	Errors  any         `json:"errors,omitempty"` // Errors: List of errors,
	Status  string      `json:"status"`
}

type SendMessage struct {
	MessageID  int64  `json:"message_id"`  // MessageID: Unique message ID
	ChannelID  int64  `json:"channel_id"`  // ChannelID: Channel ID
	OperatorID int64  `json:"operator_id"` // OperatorID: Operator ID
	Transport  string `json:"transport"`   // Transport: Transport
	Type       string `json:"type"`        // Type: Message type ('to_client', 'autoreply', 'system', 'comment')
	ClientID   int64  `json:"client_id"`   // ClientID: Client ID
	DialogID   int64  `json:"dialog_id"`   // DialogID: Dialog ID
	RequestID  int64  `json:"request_id"`  // RequestID: Request ID
}

type MessageAttachment struct {
	Name string `json:"name"` // Name: Attachment name
	Link string `json:"link"` // Link: Attachment link
}

type Message struct {
	ID              int64               `json:"id"`               // ID: Unique message ID
	Coordinates     string              `json:"coordinates"`      // Coordinates: Message coordinates (if any)
	Transport       string              `json:"transport"`        // Transport: Transport
	Type            string              `json:"type"`             // Type: Message type ('to_client', 'autoreply', 'system', 'comment')
	Read            int8                `json:"read"`             // Read: Read status (0 or 1)
	Created         string              `json:"created"`          // Created: Creation timestamp
	Pdf             string              `json:"pdf"`              // Pdf: PDF attachment (if any)
	RemoteID        string              `json:"remote_id"`        // RemoteID: Remote ID (if any)
	RecipientStatus string              `json:"recipient_status"` // RecipientStatus: Recipient status ('delivered', 'not delivered', etc.)
	AiTips          string              `json:"ai_tips"`          // AiTips: AI tips (if any)
	Attachments     []MessageAttachment `json:"attachments"`      // Attachments: List of attachments
	Photo           string              `json:"photo"`            // Photo: Photo attachment (if any)
	Video           string              `json:"video"`            // Video attachment (if any)
	Audio           string              `json:"audio"`            // Audio attachment (if any)
	OperatorID      int64               `json:"operator_id"`      // OperatorID: Operator ID
	ChannelID       int64               `json:"channel_id"`       // ChannelID: Channel ID
	DialogID        int64               `json:"dialog_id"`        // DialogID: Dialog ID
	ClientID        int64               `json:"client_id"`        // ClientID: Client ID
	RequestID       int64               `json:"request_id"`       // RequestID: Request ID
	ExtraData       any                 `json:"extra_data"`       // ExtraData: Extra data (if any)
	Status          string              `json:"status"`           // Status: Message status ('sent', 'failed', etc.)
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
func (dst *Ctd) APISendMessage(ctx context.Context, message *MessagePayload) (*SendMessageResponse, error) {
	url := fmt.Sprintf("%sv1/messages", dst.Url)
	response := SendMessageResponse{}

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

// APITransferToGroup transfers a message to a different group via the API.
// It takes a context, message ID, group ID, and force flag, and returns an error if the request fails.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - message_id (int64): The ID of the message to transfer.
//   - group_id (int64): The ID of the group to transfer the message to.
//   - force (bool): Whether to force the transfer.
//
// Returns:
//   - A pointer to a BasicResponse containing the response data.
//   - An error if the request fails.
func (dst *Ctd) APITransferToGroup(ctx context.Context, message_id, group_id int64, force bool) (*BasicResponse, error) {
	url := fmt.Sprintf("%sv1/messages/%d/transfer_to_group?group_id=%d&force=%t", dst.Url, message_id, group_id, force)
	response := BasicResponse{}

	data, err := dst.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dst.Error(ctx, "Failed transfer message to group: %v", err)
		return nil, err
	}

	str := strings.ToLower(string(data))
	if strings.Contains(str, "operator group") && strings.Contains(str, "not found") {
		return nil, ErrorInvalidOperatorGroupID
	}

	if strings.Contains(str, "message") && strings.Contains(str, "not found") {
		return nil, ErrorInvalidMesssageID
	}

	return &response, nil
}

// APITransferToOperator transfers a message to a different operator via the API.
// It takes a context, message ID, and operator ID, and returns an error if the request fails.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - message_id (int64): The ID of the message to transfer.
//   - operator_id (int64): The ID of the operator to transfer the message to.
//
// Returns:
//   - A pointer to a BasicResponse containing the response data.
//   - An error if the request fails.
func (dst *Ctd) APITransferToOperator(ctx context.Context, message_id, operator_id int64) (*BasicResponse, error) {
	url := fmt.Sprintf("%sv1/messages/%d/transfer?operator_id=%d", dst.Url, message_id, operator_id)
	response := BasicResponse{}

	data, err := dst.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dst.Error(ctx, "Failed transfer message to operator: %v", err)
		return nil, err
	}

	str := strings.ToLower(string(data))
	if strings.Contains(str, "operator") && strings.Contains(str, "not found") {
		return nil, ErrorInvalidOperatorID
	}

	if strings.Contains(str, "message") && strings.Contains(str, "not found") {
		return nil, ErrorInvalidMesssageID
	}

	return &response, nil
}

// SendMessage sends a message to the API.
// It takes a context and a MessagePayload, and returns a Message or an error.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - message (MessagePayload): The message payload to send.
//
// Returns:
//   - A pointer to a Message containing the response data.
//   - An error if the request fails.
func (dst *Ctd) SendMessage(ctx context.Context, message *MessagePayload) (*SendMessage, error) {
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

// TransferToGroup transfers a message to a different group.
// It takes a context, message ID, group ID, and force flag, and returns an error if the request fails.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - message_id (int64): The ID of the message to transfer.
//   - group_id (int64): The ID of the group to transfer the message to.
//   - force (bool): Whether to force the transfer.
//
// Returns:
//   - An error if the request fails.
func (dst *Ctd) TransferToGroup(ctx context.Context, message_id, group_id int64, force bool) error {
	data, err := dst.APITransferToGroup(ctx, message_id, group_id, force)
	if err != nil {
		return err
	}

	if data.Status != "success" {
		dst.Error(ctx, "Failed to transfer message to group: %s", data.Errors)
		return ErrorInvalidParameters
	}

	return nil
}

// TransferToOperator transfers a message to a different operator.
// It takes a context, message ID, and operator ID, and returns an error if the request fails.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - message_id (int64): The ID of the message to transfer.
//   - operator_id (int64): The ID of the operator to transfer the message to.
//
// Returns:
//   - An error if the request fails.
func (dst *Ctd) TransferToOperator(ctx context.Context, message_id, operator_id int64) error {
	data, err := dst.APITransferToOperator(ctx, message_id, operator_id)
	if err != nil {
		return err
	}

	if data.Status != "success" {
		dst.Error(ctx, "Failed to transfer message to operator: %s", data.Errors)
		return ErrorInvalidParameters
	}

	return nil
}
