package ctd

import (
	"context"
	"fmt"
)

type CustomClientFieldResponse struct {
	Status  string                  `json:"status"`           // Status: Status of the response
	Data    []CustomClientFieldItem `json:"data"`             // Data: List of custom client fields
	Message string                  `json:"message"`          // Message: Additional message from the API
	Errors  string                  `json:"errors,omitempty"` // Errors: List of errors, if any
}

// CustomClientFieldItem represents a single custom client field in the Chat2Desk API.
// It contains various fields that describe the custom client field, such as ID, name, type,
// value, and whether it is editable, viewable, or a tracking field.
// This struct is used to represent custom fields that can be associated with clients in the Chat2Desk system.
type CustomClientFieldItem struct {
	ID            int    `json:"id"`             // ID: Unique identifier for the custom client field
	Name          string `json:"name"`           // Name: Name of the custom client field
	Type          string `json:"type"`           // Type: Type of the custom client field (e.g., text, number, date)
	Value         string `json:"value"`          // Value: Value of the custom client field
	Editable      bool   `json:"editable"`       // Editable: Indicates if the custom client field is editable
	Viewable      bool   `json:"viewable"`       // Viewable: Indicates if the custom client field is viewable
	TrackingField bool   `json:"tracking_field"` // TrackingField: Indicates if the custom client field is a tracking field
}

// APICustomClientFields retrieves a list of custom client fields from the Chat2Desk API.
// It constructs the API endpoint URL, sends a GET request to the API,
// and returns the response data as a CustomClientFieldResponse struct.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a pointer to the CustomClientFieldResponse struct.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A pointer to a CustomClientFieldResponse struct containing the list of custom client fields
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) APICustomClientFields(ctx context.Context) (*CustomClientFieldResponse, error) {
	url := fmt.Sprintf("%sv1/custom_client_fields", dst.Url)

	response := CustomClientFieldResponse{}
	_, err := dst.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dst.Error(ctx, "Failed to get custom client fields: %v", err)
		return nil, err
	}
	return &response, nil
}

// GetCustomClientFields retrieves a list of custom client fields from the Chat2Desk API.
// It uses the APICustomClientFields method to fetch the custom client fields and handles errors.
// If the response status is not "success", it returns nil.
// It returns a pointer to a slice of CustomClientFieldItem, which contains the custom client fields.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A slice of CustomClientFieldItem containing the custom client fields.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) GetCustomClientFields(ctx context.Context) ([]CustomClientFieldItem, error) {
	response, err := dst.APICustomClientFields(ctx)
	if err != nil {
		return nil, err
	}

	if response.Status != "ok" {
		dst.Error(ctx, "Failed to get custom client fields: %s", response.Errors)
		return nil, ErrorInvalidResponse
	}

	return response.Data, nil
}
