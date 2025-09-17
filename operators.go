package ctd

import (
	"context"
	"fmt"
)

type CtdOperator struct {
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

type CtdOperatorsResponse struct {
	Data []CtdOperator `json:"data"`
	Meta MetaResponse  `json:"meta"`
	BasicResponse
}

// APIOperators retrieves a list of operators from the Chat2Desk API.
// It constructs the API endpoint URL with the provided offset and limit,
// sends a GET request to the API, and returns the response data as a CtdOperatorsResponse struct.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a pointer to the CtdOperatorsResponse struct.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination, indicating where to start fetching operators.
//   - limit: The maximum number of operators to return.
//
// Returns:
//   - A pointer to a CtdOperatorsResponse struct containing the list of operators and metadata.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) APIOperators(ctx context.Context, offset int, limit int) (*CtdOperatorsResponse, error) {
	url := fmt.Sprintf("%sv1/operators?offset=%d&limit=%d", dst.Url, offset, limit)

	response := CtdOperatorsResponse{}

	if _, err := dst.doRequest(ctx, "GET", url, nil, &response); err != nil {
		dst.Error(ctx, "Failed to get operators: %v", err)
		return nil, err
	}

	return &response, nil
}

// Operators retrieves a list of operators from the Chat2Desk API.
// It uses the APIOperators method to fetch the operators and handles errors.
// If the response status is not "success", it returns nil.
// It returns a pointer to a slice of CtdOperator, which contains the operators.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination, indicating where to start fetching operators.
//   - limit: The maximum number of operators to return.
//
// Returns:
//   - A pointer to a slice of CtdOperator containing the list of operators.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) Operators(ctx context.Context, offset int, limit int) (*[]CtdOperator, error) {
	data, err := dst.APIOperators(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	if data.Status != "success" {
		dst.Error(ctx, "Failed to get operators: %s", data.Errors)
		return nil, fmt.Errorf("failed to get operators: %s", data.Errors)
	}

	return &data.Data, nil
}

// AllOperators retrieves all operators from the Chat2Desk API by handling pagination.
// It repeatedly calls the Operators method with increasing offsets until all operators are fetched.
// It returns a pointer to a slice of CtdOperator, which contains all the operators.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A pointer to a slice of CtdOperator containing all the operators.
//   - An error if the request fails or if the response is invalid.
func (dst *Ctd) AllOperators(ctx context.Context) (*[]CtdOperator, error) {
	operators := []CtdOperator{}
	offset := 0
	limit := 200
	for {
		data, err := dst.Operators(ctx, offset, limit)
		if err != nil {
			return nil, err
		}

		operators = append(operators, *data...)
		if len(*data) < limit {
			break
		}

		offset += limit
	}

	return &operators, nil
}
