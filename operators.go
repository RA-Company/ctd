package ctd

import (
	"context"
	"fmt"
)

type Operator struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Avatar        string `json:"avatar"`
	Role          string `json:"role"`
	StatusID      uint8  `json:"status_id"`
	OpenedDialogs int64  `json:"opened_dialogs"`
	Online        uint8  `json:"online"`
}

type OperatorsResponse struct {
	Data []Operator   `json:"data"`
	Meta MetaResponse `json:"meta"`
	BasicResponse
}

// APIOperators retrieves a list of operators from the Chat2Desk API.
// It constructs the API endpoint URL with the provided offset and limit,
// sends a GET request to the API, and returns the response data as a OperatorsResponse struct.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a pointer to the OperatorsResponse struct.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination, indicating where to start fetching operators.
//   - limit: The maximum number of operators to return.
//
// Returns:
//   - A pointer to a CtdOperatorsResponse struct containing the list of operators and metadata.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) APIOperators(ctx context.Context, offset int, limit int) (*OperatorsResponse, error) {
	url := fmt.Sprintf("%sv1/operators?offset=%d&limit=%d", dst.Url, offset, limit)

	response := OperatorsResponse{}

	if _, err := dst.doRequest(ctx, "GET", url, nil, &response); err != nil {
		dst.Error(ctx, "Failed to get operators: %v", err)
		return nil, err
	}

	return &response, nil
}

// Operators retrieves a list of operators from the Chat2Desk API.
// It uses the APIOperators method to fetch the operators and handles errors.
// If the response status is not "success", it returns nil.
// It returns a slice of Operator, which contains the operators.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination, indicating where to start fetching operators.
//   - limit: The maximum number of operators to return.
//
// Returns:
//   - A a slice of Operator containing the list of operators.
//   - The total number of operators available (for pagination).
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) Operators(ctx context.Context, offset int, limit int) ([]Operator, int, error) {
	data, err := dst.APIOperators(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	if data.Status != "success" {
		dst.Error(ctx, "Failed to get operators: %s", data.Errors)
		return nil, 0, fmt.Errorf("failed to get operators: %s", data.Errors)
	}

	return data.Data, data.Meta.Total, nil
}

// AllOperators retrieves all operators from the Chat2Desk API by handling pagination.
// It repeatedly calls the Operators method with increasing offsets until all operators are fetched.
// It returns a slice of Operator, which contains all the operators.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A slice of Operator containing all the operators.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) AllOperators(ctx context.Context) ([]Operator, error) {
	operators := []Operator{}
	offset := 0
	limit := 200
	for {
		data, _, err := dst.Operators(ctx, offset, limit)
		if err != nil {
			return nil, err
		}

		operators = append(operators, data...)
		if len(data) < limit {
			break
		}

		offset += limit
	}

	return operators, nil
}
