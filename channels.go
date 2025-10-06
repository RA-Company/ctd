package ctd

import (
	"context"
	"fmt"
)

// ChannelItem represents a single channel in the Chat2Desk API
// with its ID, name, phone number, and transports.
// It is used in the ChannelsResponse to provide a list of channels.
type ChannelItem struct {
	ID         int      `json:"id"`              // ID: Unique identifier of the channel
	Name       string   `json:"name,omitempty"`  // Name: Name of the channel
	Phone      string   `json:"phone,omitempty"` // Phone: Phone number associated with the channel
	Transports []string `json:"transports"`      // Transports: List of transports used by the channel
}

// ChannelsResponse represents the response from the Chat2Desk API
// when fetching channels. It includes metadata about the response
// and a status message.
// It contains a list of ChannelItem objects that represent the channels.
// The MetaResponse provides pagination information such as total count,
// limit, and offset for the channels returned.
// It is used to encapsulate the response structure for the Channels API endpoint.
type ChannelsResponse struct {
	Data   []ChannelItem `json:"data"` // Data: List of channels
	Meta   MetaResponse  `json:"meta"`
	Status string        `json:"status"`
}

// Channels retrieves a list of channels from the Chat2Desk API.
// It takes a context, an offset, and a limit as parameters.
// The offset is used for pagination, and the limit specifies the maximum
// number of channels to return.
// It constructs the API endpoint URL with the provided offset and limit,
// sends a GET request to the API, and returns the response data as a byte slice.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns the response data.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination, indicating where to start fetching channels.
//   - limit: The maximum number of channels to return.
//
// Returns:
//   - A pointer to a ChannelsResponse struct containing the list of channels and metadata.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) Channels(ctx context.Context, offset, limit int) (*ChannelsResponse, error) {
	url := fmt.Sprintf("%sv1/channels?offset=%d&limit=%d", dst.Url, offset, limit)

	response := ChannelsResponse{}

	_, err := dst.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dst.Error(ctx, "Failed to get channels: %v", err)
		return nil, err
	}

	/*err = json.Unmarshal(data, &response)
	if err != nil {
		dst.Error(ctx, "Failed to unmarshal channels response: %v", err)
		return nil, ErrorInvalidResponse
	}*/

	return &response, nil
}

// GetChannels retrieves a list of channels from the Chat2Desk API.
// It uses the Channels method to fetch the channels and handles errors.
// If the response status is not "success", it logs an error and returns nil.
// It returns a slice of ChannelItem, which contains the channels.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination, indicating where to start fetching channels.
//   - limit: The maximum number of channels to return.
//
// Returns:
//   - A slice of ChannelItem containing the channels.
//   - The total number of channels available (for pagination).
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) GetChannels(ctx context.Context, offset, limit int) ([]ChannelItem, int, error) {
	response, err := dst.Channels(ctx, offset, limit)
	if err != nil {
		dst.Error(ctx, "Failed to get channels: %v", err)
		return nil, 0, err
	}

	if response.Status != "success" {
		dst.Error(ctx, "Invalid response status: %s", response.Status)
		return nil, 0, ErrorInvalidResponse
	}

	return response.Data, response.Meta.Total, nil
}
