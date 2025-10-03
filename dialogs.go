package ctd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ra-company/jsons"
)

type DialogResponse struct {
	BasicResponse
	Data DialogItem `json:"data"` // Data: List of clients
}

type DialogsResponse struct {
	BasicResponse
	Data []DialogItem `json:"data"` // Data: List of clients
	Meta MetaResponse `json:"meta"` // Meta: Meta information
}

type DialogItem struct {
	ID            int64       `json:"id"`              // ID: Dialog ID
	State         string      `json:"state"`           // State: Dialog state
	Begin         jsons.Time  `json:"begin"`           // Begin: Dialog begin time
	End           jsons.Time  `json:"end"`             // End: Dialog end time
	LastMessage   MessageItem `json:"last_message"`    // LastMessage: Last message in dialog
	LastRequestID int64       `json:"last_request_id"` // LastRequestID: Last request ID
	Messages      json.Number `json:"messages"`        // Messages: Number of messages in dialog
	OperatorID    int64       `json:"operator_id"`     // OperatorID: Operator ID
}

type GetDialogsParams struct {
	Limit      int    `json:"limit,omitempty"`       // Limit: Optional limit of dialogs to retrieve (default: 100, max: 1000)
	Offset     int    `json:"offset,omitempty"`      // Offset: Optional offset for pagination (default: 0)
	State      string `json:"state,omitempty"`       // State: Optional filter by dialog state ('open', 'closed', '' (default: ''))
	OperatorID int    `json:"operator_id,omitempty"` // OperatorID: Optional filter by operator ID
	Order      string `json:"order,omitempty"`       // Order: Optional order of results ('asc' or 'desc', default: '')
}

func (p *GetDialogsParams) Params() string {
	params := []string{}

	if p.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", p.Limit))
	}
	if p.Offset > 0 {
		params = append(params, fmt.Sprintf("offset=%d", p.Offset))
	}
	if p.State != "" {
		if p.State == "open" || p.State == "closed" {
			params = append(params, fmt.Sprintf("state=%s", p.State))
		}
	}
	if p.OperatorID > 0 {
		params = append(params, fmt.Sprintf("operator_id=%d", p.OperatorID))
	}
	if p.Order != "" {
		if p.Order == "asc" || p.Order == "desc" {
			params = append(params, fmt.Sprintf("order=%s", p.Order))
		}
	}

	if len(params) > 0 {
		return "?" + strings.Join(params, "&")
	}
	return ""
}

// APIGetDialogs retrieves a list of dialogs from the API.
// It takes a context and GetDialogsParams, and returns a DialogsResponse or an error.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - params (*GetDialogsParams): The parameters for filtering and pagination.
//
// Returns:
//   - A pointer to a DialogsResponse containing the response data.
//   - An error if the request fails.
func (dst *Ctd) APIGetDialogs(ctx context.Context, params *GetDialogsParams) (*DialogsResponse, error) {
	url := fmt.Sprintf("%sv1/dialogs%s", dst.Url, params.Params())
	response := DialogsResponse{}

	if _, err := dst.doRequest(ctx, "GET", url, nil, &response); err != nil {
		dst.Error(ctx, "Failed get dialogs: %v", err)
		return nil, err
	}
	return &response, nil
}

// APIGetDialog retrieves a dialog by its ID from the API.
// It takes a context and a dialog ID, and returns a DialogResponse or an error.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - dialog_id (int64): The ID of the dialog to retrieve.
//
// Returns:
//   - A pointer to a DialogResponse containing the response data.
//   - An error if the request fails.
func (dst *Ctd) APIGetDialog(ctx context.Context, dialog_id int64) (*DialogResponse, error) {
	url := fmt.Sprintf("%sv1/dialogs/%d", dst.Url, dialog_id)
	response := DialogResponse{}

	if _, err := dst.doRequest(ctx, "GET", url, nil, &response); err != nil {
		dst.Error(ctx, "Failed get dialog by ID: %v", err)
		return nil, err
	}
	return &response, nil
}

// GetDialogs retrieves a list of dialogs.
// It takes a context and GetDialogsParams, and returns a slice of DialogItem or an error.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - params (*GetDialogsParams): The parameters for filtering and pagination.
//
// Returns:
//   - A pointer to a slice of DialogItem containing the dialogs.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) GetDialogs(ctx context.Context, params *GetDialogsParams) (*[]DialogItem, error) {
	data, err := dst.APIGetDialogs(ctx, params)
	if err != nil {
		return nil, err
	}

	if data.Status != "success" {
		dst.Error(ctx, "Failed to get dialogs: %s", data.Errors)
		return nil, ErrorInvalidParameters
	}

	return &data.Data, nil
}

// GetDialog retrieves a dialog by its ID.
// It takes a context and a dialog ID, and returns a DialogItem or an error.
//
// Parameters:
//   - ctx (context.Context): The context for the request.
//   - dialog_id (int64): The ID of the dialog to retrieve.
//
// Returns:
//   - A pointer to a DialogItem containing the dialog data.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) GetDialog(ctx context.Context, dialog_id int64) (*DialogItem, error) {
	data, err := dst.APIGetDialog(ctx, dialog_id)
	if err != nil {
		return nil, err
	}

	if data.Status != "success" {
		dst.Error(ctx, "Failed to get dialog by ID: %s", data.Errors)
		if data.Message == "not_found" {
			return nil, ErrorInvalidID
		}
		return nil, ErrorInvalidParameters
	}

	return &data.Data, nil
}
