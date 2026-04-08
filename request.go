package ctd

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

type RequestMessageExtraData struct {
	SystemType string `json:"system_type"`
}

type RequestMessage struct {
	ID            int64                   `json:"id"`             // ID: Message ID
	Text          string                  `json:"text"`           // Text: Message text
	Video         string                  `json:"video"`          // Video: URL of the video in the message
	Photo         string                  `json:"photo"`          // Photo: URL of the photo in the message
	Audio         string                  `json:"audio"`          // Audio: URL of the audio in the message
	PDF           string                  `json:"pdf"`            // PDF: URL of the PDF in the message
	Coordinates   string                  `json:"coordinates"`    // Coordinates: Coordinates in the message
	Transport     string                  `json:"transport"`      // Transport: Transport type of the message
	Type          string                  `json:"type"`           // Type: Type of the message ("in", or "out" or "system")
	Read          int8                    `json:"read"`           // Read: Read status of the message (0 or 1)
	Created       int64                   `json:"created"`        // Created: Creation time of the message in Unix timestamp
	ClientID      int64                   `json:"clientID"`       // ClientID: Client ID associated with the message
	OperatorID    int64                   `json:"operatorID"`     // OperatorID: Operator ID associated with the message
	ChannelID     int64                   `json:"channelID"`      // ChannelID: Channel ID associated with the message
	CompanyID     int64                   `json:"companyID"`      // CompanyID: Company ID associated with the message
	DialogID      int64                   `json:"dialogID"`       // DialogID: Dialog ID associated with the message
	GatewayStatus string                  `json:"gateway_status"` // GatewayStatus: Gateway status of the message
	ExtraData     RequestMessageExtraData `json:"extra_data"`     // ExtraData: Additional data for the message
}

// CreatedTime returns the creation time of the message as a time.Time object.
//
// Returns:
//   - time.Time: The creation time of the message.
func (dst *RequestMessage) CreatedTime() time.Time {
	return time.Unix(dst.Created, 0)
}

// Message returns the content of the message, prioritizing text, video, photo, audio, PDF, and coordinates in that order.
//
// Returns:
//   - string: The content of the message based on the available fields.
func (dst *RequestMessage) Message() string {
	if dst.Text != "" {
		return dst.Text
	}

	if dst.Video != "" {
		return dst.Video
	}

	if dst.Photo != "" {
		return dst.Photo
	}

	if dst.Audio != "" {
		return dst.Audio
	}

	if dst.PDF != "" {
		return dst.PDF
	}

	if dst.Coordinates != "" {
		return dst.Coordinates
	}

	return ""
}

// MessageFormat returns the format of the message based on the available fields, prioritizing text, video, photo, audio, PDF, and coordinates in that order.
//
// Returns:
//   - string: The format of the message based on the available fields ("text", "video", "photo", "audio", "pdf", "coordinates", or "unknown").
func (dst *RequestMessage) MessageFormat() string {
	if dst.Text != "" {
		return "text"
	}

	if dst.Video != "" {
		return "video"
	}

	if dst.Photo != "" {
		return "photo"
	}

	if dst.Audio != "" {
		return "audio"
	}

	if dst.PDF != "" {
		return "pdf"
	}

	if dst.Coordinates != "" {
		return "coordinates"
	}

	return "unknown"
}

// APIRequestMessages retrieves a list of messages for a specific request from the Chat2Desk API.
// It constructs the API endpoint URL using the provided request ID, sends a GET request to the API,
// and returns the response data as a slice of RequestMessage structs.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a slice of RequestMessage structs.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - request: The ID of the request for which to retrieve messages.
//
// Returns:
//   - A slice of RequestMessage structs containing the list of messages for the specified request or nil if an error occurs.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) APIRequestMessages(ctx context.Context, request int64) ([]RequestMessage, error) {
	url := fmt.Sprintf("%sv1/requests/%d/messages", dst.Url, request)

	response := []RequestMessage{}

	body, err := dst.doRequest(ctx, "GET", url, nil, &response)
	if err == ErrorInvalidResponse {
		if strings.Contains(string(body), "not_found") {
			dst.Error(ctx, "Invalid request ID: %d", request)
			return nil, ErrorInvalidRequestID
		}
	}

	if err != nil {
		dst.Error(ctx, "Failed to get request messages: %v", err)
		return nil, err
	}

	return response, nil
}

// RequestMessages is a wrapper around APIRequestMessages that retrieves messages for a specific request.
// It calls APIRequestMessages with the provided request ID and returns the result. Then its sort messages by creation time in ascending order.
// If creation time is the same, it sorts by ID in ascending order.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - request: The ID of the request for which to retrieve messages.
//
// Returns:
//   - A slice of RequestMessage structs containing the list of messages for the specified request or nil if an error occurs.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) RequestMessages(ctx context.Context, request int64) ([]RequestMessage, error) {
	messages, err := dst.APIRequestMessages(ctx, request)
	if err != nil {
		return nil, err
	}

	// Sort messages by creation time in ascending order. If creation time is the same, sort by ID in ascending order.
	sort.Slice(messages, func(i, j int) bool {
		if messages[i].Created == messages[j].Created {
			return messages[i].ID < messages[j].ID
		}
		return messages[i].Created < messages[j].Created
	})

	return messages, nil
}
