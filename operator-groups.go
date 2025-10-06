package ctd

import (
	"context"
	"fmt"
)

type OperatorGroup struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Operators []int64 `json:"operator_ids"`
}

type OperatorGroupsResponse struct {
	BasicResponse
	Data []OperatorGroup `json:"data"`
}

// APIOperatorGroups retrieves a list of operator groups from the Chat2Desk API.
// It constructs the API endpoint URL, sends a GET request to the API,
// and returns the response data as an OperatorGroupsResponse struct.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a pointer to the OperatorGroupsResponse struct.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A pointer to an OperatorGroupsResponse struct containing the list of operator groups
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) APIOperatorGroups(ctx context.Context) (*OperatorGroupsResponse, error) {
	url := fmt.Sprintf("%sv1/operators_groups", dst.Url)
	response := OperatorGroupsResponse{}

	if _, err := dst.doRequest(ctx, "GET", url, nil, &response); err != nil {
		dst.Error(ctx, "Failed send message: %v", err)
		return nil, err
	}
	return &response, nil
}

// OperatorGroups retrieves a list of operator groups from the Chat2Desk API.
// It uses the APIOperatorGroups method to fetch the operator groups and handles errors.
// If the response status is not "success", it returns nil.
// It returns a pointer to a slice of OperatorGroup, which contains the operator groups.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A slice of OperatorGroup containing the list of operator groups.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) OperatorGroups(ctx context.Context) ([]OperatorGroup, error) {
	data, err := dst.APIOperatorGroups(ctx)
	if err != nil {
		return nil, err
	}

	if data.Status != "success" {
		dst.Error(ctx, "Failed to get operator groups: %s", data.Errors)
		return nil, ErrorInvalidParameters
	}

	return data.Data, nil
}
