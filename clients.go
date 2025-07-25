package ctd

import (
	"context"
	"fmt"
	"strings"
)

// ClientResponse represents the response structure for client-related API calls.
type ClientResponse struct {
	Data    ClientItem `json:"data"` // Data: List of clients
	Message string     `json:"message"`
	Errors  string     `json:"errors,omitempty"` // Errors: List of errors,
	Status  string     `json:"status"`
}

type ClientsResponse struct {
	Data    []ClientItem `json:"data"` // Data: List of clients
	Meta    MetaResponse `json:"meta"`
	Message string       `json:"message"`
	Errors  string       `json:"errors,omitempty"` // Errors: List of errors,
	Status  string       `json:"status"`
}

// ClientItem represents a single client in the Chat2Desk API.
// It contains various fields that describe the client, such as ID, name, phone number,
// avatar, region, country, messages, comments, custom fields, and associated channels and tags.
type ClientItem struct {
	ID                    int               `json:"id"`                   // ID: Unique identifier of the client
	Name                  string            `json:"name"`                 // Name: Name of the client
	Username              string            `json:"username"`             // Username: Username of the client
	Comment               string            `json:"comment"`              // Comment: Comment associated with the client
	AssignedName          string            `json:"assigned_name"`        // AssignedName: Name of the assigned user
	Phone                 string            `json:"phone"`                // Phone: Phone number of the client
	ClientPhone           string            `json:"client_phone"`         // ClientPhone: Client's phone number
	Avatar                string            `json:"avatar"`               // Avatar: URL of the client's avatar
	RegionID              int               `json:"region_id"`            // RegionID: ID of the region associated with the client
	CountryID             int               `json:"country_id"`           // CountryID: ID of the country associated with the client
	FirstClientMessageStr string            `json:"first_client_message"` // FirstClientMessageStr: String representation of the first client message
	LastClientMessageStr  string            `json:"last_client_message"`  // LastClientMessageStr: String representation of the last client message
	ExtraComment1         string            `json:"extra_comment_1"`      // ExtraComment1: First extra comment associated with the client
	ExtraComment2         string            `json:"extra_comment_2"`      // ExtraComment2: Second extra comment associated with the client
	ExtraComment3         string            `json:"extra_comment_3"`      // ExtraComment3: Third extra comment associated with the client
	CustomFields          map[string]string `json:"custom_fields"`        // CustomFields: Map of custom fields associated with the client
	ClientExternalID      string            `json:"client_external_id"`   // ClientExternalID: External ID of the client
	ExtrnalID             int               `json:"external_id"`          // ExternalID: External ID of the client
	ExtrnalIDs            map[string]int    `json:"external_ids"`         // ExternalIDs: Map of external IDs associated with the client
	Channels              []struct {
		ID         int      `json:"id"`         // ID: Unique identifier of the channel
		Transports []string `json:"transports"` // Transports: List of transports used by the channel
	} `json:"channels"` // Channels: List of channels associated with the client
	Tags []struct {
		ID          int    `json:"id"`          // ID: Unique identifier of the tag
		Label       string `json:"label"`       // Label: Label of the tag
		Description string `json:"description"` // Description: Description of the tag
		GroupID     int    `json:"group_id"`    // GroupID: ID of the group associated with the tag
		GroupName   string `json:"group_name"`  // GroupName: Name of the group associated with the tag
	} `json:"tags"` // Tags: List of tags associated with the client
}

// APIGetClient retrieves a client by its ID from the Chat2Desk API.
// It takes a context and the client ID as parameters.
// It constructs the API endpoint URL with the provided client ID,
// sends a GET request to the API, and returns the response data as a byte slice.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns the response data.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the client to retrieve.
//
// Returns:
//   - A pointer to a ClientsResponse struct containing the client data and metadata.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) APIGetClient(ctx context.Context, id int) (*ClientResponse, error) {
	url := fmt.Sprintf("%sv1/clients/%d", dst.Url, id)
	response := ClientResponse{}
	_, err := dst.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dst.Error(ctx, "Failed to get client: %v", err)
		return nil, err
	}
	if response.Status != "success" {
		dst.Error(ctx, "Failed to get client: %s", response.Errors)
		return nil, ErrorInvalidID
	}
	return &response, nil
}

// APIGetClients retrieves a list of clients from the Chat2Desk API.
// It takes a context, an offset, a limit, an order, and additional parameters as strings.
// The offset is used for pagination, the limit specifies the maximum number of clients to return,
// the order specifies the sorting order, and params can include additional query parameters.
// It constructs the API endpoint URL with the provided parameters,
// sends a GET request to the API, and returns the response data as a byte slice.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns the response data.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination, indicating where to start fetching clients.
//   - limit: The maximum number of clients to return.
//   - order: The sorting order for the clients (e.g., "asc", "desc").
//   - params: Additional query parameters to include in the request.
//
// Returns:
//   - A pointer to a ClientResponse struct containing the list of clients and metadata.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) APIGetClients(ctx context.Context, offset, limit int, order, params string) (*ClientsResponse, error) {
	order = strings.ToLower(order)
	if order != "desc" {
		order = "asc"
	}
	if params != "" {
		params = "&" + params
	}
	url := fmt.Sprintf("%sv1/clients?offset=%d&limit=%d&order=%s%s", dst.Url, offset, limit, order, params)
	response := ClientsResponse{}
	_, err := dst.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dst.Error(ctx, "Failed to get clients: %v", err)
		return nil, err
	}
	return &response, nil
}

// GetClient retrieves a client by its ID from the Chat2Desk API.
// It takes a context and the client ID as parameters.
// It calls the APIGetClient method to fetch the client data.
// If the response contains an error or if no client data is found, it returns an error.
// If the client is found, it returns a pointer to the ClientItem struct containing the client details.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the client to retrieve.
//
// Returns:
//   - A pointer to a ClientItem struct containing the client details.
//   - An error if the request fails, if the response is invalid, or if no client data is found.
func (dst *Ctd) GetClient(ctx context.Context, id int) (*ClientItem, error) {
	response, err := dst.APIGetClient(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.Contains(response.Errors, " not found") {
		return nil, ErrorInvalidID
	}

	if response.Status != "success" {
		return nil, ErrorInvalidResponse
	}

	return &response.Data, nil
}

// GetClients retrieves a list of clients from the Chat2Desk API.
// It uses the APIGetClients method to fetch the clients and handles errors.
// If the response status is not "success", it returns nil.
// It returns a pointer to a slice of ClientItem, which contains the clients.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination, indicating where to start fetching clients.
//   - limit: The maximum number of clients to return.
//
// Returns:
//   - A pointer to a slice of ClientItem containing the clients.
//   - The total number of clients available (for pagination).
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) GetClientsList(ctx context.Context, offset, limit int) (*[]ClientItem, int, error) {
	response, err := dst.APIGetClients(ctx, offset, limit, "asc", "")
	if err != nil {
		return nil, 0, err
	}

	if response.Status != "success" {
		return nil, 0, ErrorInvalidResponse
	}

	if len(response.Data) == 0 {
		return nil, 0, nil // No clients found
	}

	return &response.Data, response.Meta.Total, nil
}
