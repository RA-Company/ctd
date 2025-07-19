package ctd

import (
	"context"
	"fmt"
	"strings"
)

// WebhooksResponse represents the response from the Chat2Desk API
// when fetching webhooks. It includes a list of WebhookItem objects
// that represent the webhooks and a status message.
// It is used to encapsulate the response structure for the Webhooks API endpoint.
type WebhooksResponse struct {
	Data   []WebhookItem `json:"data"` // Data: List of webhooks
	Status string        `json:"status"`
}

// CreateWebhookResponse represents the response from the Chat2Desk API
// when fetching a single webhook. It includes a WebhookItem object
// that represents the webhook and a status message.
// It is used to encapsulate the response structure for the CreateWebhook API endpoint.
type CreateWebhookResponse struct {
	Data    WebhookItem `json:"data"` // Data: Single webhook item
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"` // Message: Optional message providing additional information
	Errors  struct {
		Url    []string `json:"url,omitempty"`    // Url: List of errors related to the URL
		Order  []string `json:"order,omitempty"`  // Order: List of errors related to the order
		Events []string `json:"events,omitempty"` // Events: List of errors related to the events
	} `json:"errors"` // Errors: Optional field containing errors related to the request
}

func (dst *CreateWebhookResponse) Error() string {
	errs := []string{}
	if len(dst.Errors.Url) > 0 {
		errs = append(errs, fmt.Sprintf("URL: %s", strings.Join(dst.Errors.Url, ", ")))
	}
	if len(dst.Errors.Order) > 0 {
		errs = append(errs, fmt.Sprintf("Order: %s", strings.Join(dst.Errors.Order, ", ")))
	}
	if len(dst.Errors.Events) > 0 {
		errs = append(errs, fmt.Sprintf("Events: %s", strings.Join(dst.Errors.Events, ", ")))
	}
	if len(errs) > 0 {
		return strings.Join(errs, "; ")
	}

	return ""
}

// Postprocess processes the response from the CreateWebhook API endpoint.
// It checks the status of the response and returns an error if the status is not "success".
func (dst *CreateWebhookResponse) Postprocess() error {
	if dst.Status != "success" {
		if dst.Errors.Url != nil {
			for _, errMsg := range dst.Errors.Url {
				if strings.Contains(errMsg, "already used") {
					return ErrorWebhookUrlIsAlreadyUsed
				}
			}
		}
		return ErrorInvalidResponse
	}

	return nil
}

// DeleteWebhookResponse represents the response from the Chat2Desk API
// when deleting a webhook. It includes a status message and optional fields
// for additional information and errors.
type DeleteWebhookResponse struct {
	Status  string `json:"status"`            // Status: Status of the delete operation
	Message string `json:"message,omitempty"` // Message: Optional message providing additional information
	Errors  string `json:"errors,omitempty"`  // Errors: Optional field containing errors related to the delete operation
}

// WebhookItem represents a single webhook in the Chat2Desk API
// with its ID, name, URL, and transports.
// It is used in the WebhooksResponse to provide a list of webhooks.
type WebhookItem struct {
	ID       int             `json:"id"`       // ID: Unique identifier of the webhook
	Name     string          `json:"name"`     // Name: Name of the webhook
	URL      string          `json:"url"`      // URL: The URL to which the webhook
	Events   []string        `json:"events"`   // Events: List of events that trigger the webhook
	Status   string          `json:"status"`   // Status: Status of the webhook (e.g., enable, disable)
	Errors   []WebhookErrors `json:"errors"`   // Errors: List of errors related to the webhook
	Source   string          `json:"source"`   // Source: Source of the webhook (e.g., client, server)
	Channels []int           `json:"channels"` // Channels: List of channel IDs associated with the webhook
}

// WebhookErrors represents errors related to webhooks in the Chat2Desk API
// It contains a text message and a timestamp indicating when the error occurred.
// It is used to provide details about errors encountered during webhook operations.
type WebhookErrors struct {
	Text      string `json:"text"`       // Text: Error message text
	CreatedAt uint64 `json:"created_at"` // CreatedAt: Timestamp of when the error occurred
}

// WebhookPayload represents the payload structure for creating or updating a webhook
// in the Chat2Desk API. It includes fields for the webhook name, URL, events,
// channels, status, and order.
// It is used to send data when creating or updating webhooks via the API.
type WebhookPayload struct {
	Name     string   `json:"name"`            // Name: Name of the webhook
	URL      string   `json:"url"`             // URL: The URL to which the webhook will be sent
	Events   []string `json:"events"`          // Events: List of events that will trigger the webhook
	Channels []int    `json:"channels"`        // Channels: List of channel IDs associated with the webhook
	Status   string   `json:"status"`          // Status: Status of the webhook (e.g., enable, disable)
	Order    int      `json:"order,omitempty"` // Order: Order of the webhook in the list
}

// Prepare normalizes the status field of the WebhookPayload.
// It ensures that the status is set to either "enable" or "disable".
// If the status is not one of these values, it defaults to "enable".
// This method is typically used to ensure that the status field is in a valid format
// before sending the payload to the API.
// It is called before creating or updating a webhook to ensure consistency.
// It is used to prepare the payload for API requests.
func (dst *WebhookPayload) Prepare() {
	dst.Status = strings.ToLower(dst.Status)
	if dst.Status != "enable" && dst.Status != "disable" {
		dst.Status = "enable"
	}
}

// GetWebhooks retrieves a list of webhooks from the Chat2Desk API.
// It takes a context as a parameter and constructs the API endpoint URL.
// It sends a GET request to the API and returns the response data as a byte slice.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it unmarshals the response data into a WebhooksResponse
// struct and returns it.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A pointer to a WebhooksResponse struct containing the list of webhooks and status.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) Webhooks(ctx context.Context) (*WebhooksResponse, error) {
	url := fmt.Sprintf("%sv1/webhooks", dst.Url)

	response := WebhooksResponse{}

	_, err := dst.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dst.Error(ctx, "Failed to get webhooks: %v", err)
		return nil, err
	}

	return &response, nil
}

// PostWebhook creates a new webhook in the Chat2Desk API.
// It takes a context and a WebhookPayload as parameters.
// It constructs the API endpoint URL and sends a POST request with the payload.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it unmarshals the response data into a WebhookResponse
// struct and returns it.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - payload: The WebhookPayload containing the details of the webhook to be created.
//
// Returns:
//   - A pointer to a WebhookResponse struct containing the created webhook and status.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) PostWebhooks(ctx context.Context, payload WebhookPayload) (*CreateWebhookResponse, error) {
	url := fmt.Sprintf("%sv1/webhooks", dst.Url)

	payload.Prepare() // Ensure the payload is prepared before sending

	var response CreateWebhookResponse

	_, err := dst.doRequest(ctx, "POST", url, payload, &response)
	if err != nil {
		dst.Error(ctx, "Failed to create webhook: %v", err)
		return nil, err
	}

	return &response, nil
}

// PutWebhooks updates an existing webhook in the Chat2Desk API.
// It takes a context, the webhook ID, and a WebhookPayload as parameters.
// It constructs the API endpoint URL with the webhook ID and sends a PUT request with the payload.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it unmarshals the response data into a CreateWebhookResponse
// struct and returns it.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the webhook to be updated.
//   - payload: The WebhookPayload containing the updated details of the webhook.
//
// Returns:
//   - A pointer to a CreateWebhookResponse struct containing the updated webhook and status.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) PutWebhooks(ctx context.Context, id int, payload WebhookPayload) (*CreateWebhookResponse, error) {
	url := fmt.Sprintf("%sv1/webhooks/%d", dst.Url, id)

	payload.Prepare() // Ensure the payload is prepared before sending

	var response CreateWebhookResponse

	_, err := dst.doRequest(ctx, "PUT", url, payload, &response)
	if err != nil {
		dst.Error(ctx, "Failed to update webhook: %v", err)
		return nil, err
	}

	return &response, nil
}

// DeleteWebhooks deletes a webhook in the Chat2Desk API.
// It takes a context and the webhook ID as parameters.
// It constructs the API endpoint URL with the webhook ID and sends a DELETE request.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it unmarshals the response data into a DeleteWebhookResponse
// struct and returns it.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the webhook to be deleted.
//
// Returns:
//   - A pointer to a DeleteWebhookResponse struct containing the status of the delete operation.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) DeleteWebhooks(ctx context.Context, id int) (*DeleteWebhookResponse, error) {
	url := fmt.Sprintf("%sv1/webhooks/%d", dst.Url, id)

	response := DeleteWebhookResponse{}
	_, err := dst.doRequest(ctx, "DELETE", url, nil, &response)
	if err != nil {
		dst.Error(ctx, "Failed to delete webhook: %v", err)
		return nil, err
	}

	return &response, nil
}

// GetWebhooks retrieves a list of webhooks from the Chat2Desk API.
// It takes a context as a parameter and calls the Webhooks method.
// If the response status is not "success", it logs an error and returns nil.
// It returns a pointer to a slice of WebhookItem, which contains the webhooks.
// If an error occurs during the request, it returns nil and the error.
// If the request is successful, it returns a pointer to a slice of WebhookItem.
// This method is typically used to fetch webhooks from the Chat2Desk API.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A pointer to a slice of WebhookItem containing the webhooks.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) GetWebhooks(ctx context.Context) (*[]WebhookItem, error) {
	response, err := dst.Webhooks(ctx)
	if err != nil {
		return nil, err
	}

	if response.Status != "success" {
		dst.Error(ctx, "Failed to get webhooks: %s", response.Status)
		return nil, ErrorInvalidResponse
	}

	return &response.Data, nil
}

// CreateWebhook creates a new webhook in the Chat2Desk API.
// It takes a context and a WebhookPayload as parameters.
// It calls the PostWebhook method to send the request.
// If the response status is not "success", it logs an error and returns nil.
// If the URL is already used, it returns an error indicating that the URL is already used.
// If the request is successful, it returns a pointer to the created WebhookItem.
// This method is typically used to create new webhooks in the Chat2Desk API.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - payload: The WebhookPayload containing the details of the webhook to be created.
//
// Returns:
//   - A pointer to a WebhookItem containing the created webhook.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) CreateWebhook(ctx context.Context, payload WebhookPayload) (*WebhookItem, error) {
	response, err := dst.PostWebhooks(ctx, payload)
	if err != nil {
		return nil, err
	}

	if err := response.Postprocess(); err != nil {
		dst.Error(ctx, "Failed to create webhook: %+v", response.Error())
		return nil, err
	}

	return &response.Data, nil
}

// UpdateWebhook updates an existing webhook in the Chat2Desk API.
// It takes a context, the webhook ID, and a WebhookPayload as parameters.
// It calls the PutWebhooks method to send the request.
// If the response status is not "success", it logs an error and returns nil.
// If the URL is already used, it returns an error indicating that the URL is already used.
// If the request is successful, it returns a pointer to the updated WebhookItem.
// This method is typically used to update existing webhooks in the Chat2Desk API.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the webhook to be updated.
//   - payload: The WebhookPayload containing the updated details of the webhook.
//
// Returns:
//   - A pointer to a WebhookItem containing the updated webhook.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) UpdateWebhook(ctx context.Context, id int, payload WebhookPayload) (*WebhookItem, error) {
	response, err := dst.PutWebhooks(ctx, id, payload)
	if err != nil {
		return nil, err
	}

	if err := response.Postprocess(); err != nil {
		dst.Error(ctx, "Failed to update webhook: %v", response.Error())
		return nil, err
	}

	return &response.Data, nil
}

// DeleteWebhook deletes a webhook in the Chat2Desk API.
// It takes a context and the webhook ID as parameters.
// It calls the DeleteWebhooks method to send the request.
// If the response status is not "success", it logs an error and returns an error.
// If the request is successful, it returns nil.
// This method is typically used to delete webhooks in the Chat2Desk API.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the webhook to be deleted.
//
// Returns:
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) DeleteWebhook(ctx context.Context, id int) error {
	response, err := dst.DeleteWebhooks(ctx, id)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		dst.Error(ctx, "Failed to delete webhook: %s", response.Errors)
		return ErrorInvalidID
	}

	return nil
}
