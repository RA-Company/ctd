package ctd

import (
	"context"
)

type CompaniesApiInfoData struct {
	CompanyID   int64  `json:"companyID"`
	PartnerID   int64  `json:"partnerID"`
	CompanyName string `json:"company_name"`
	AdminEmail  string `json:"admin_email"`
}

type CompaniesApiInfoResponse struct {
	Status  string               `json:"status"`           // Status: Status of the response
	Data    CompaniesApiInfoData `json:"data"`             // Data: Information about the company
	Message string               `json:"message"`          // Message: Additional message from the API
	Errors  string               `json:"errors,omitempty"` // Errors: List of errors, if any
}

// APIGetCompaniesApiInfo retrieves information about the company using the Chat2Desk API.
// It constructs the API endpoint URL, sends a GET request to the API,
// and returns the response data as a CompaniesApiInfoResponse struct.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a pointer to the CompaniesApiInfoResponse struct.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
func (dst *Ctd) APICompaniesApiInfo(ctx context.Context) (*CompaniesApiInfoResponse, error) {
	var data CompaniesApiInfoResponse

	_, err := dst.Get(ctx, "v1/companies/api_info", &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetCompaniesApiInfo retrieves information about the company using the Chat2Desk API.
// It uses the APICompaniesApiInfo method to fetch the company information and handles errors.
// If the response status is not "success", it returns nil.
// It returns a pointer to a CompaniesApiInfoData struct, which contains the company information.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
func (dst *Ctd) CompaniesApiInfo(ctx context.Context) (*CompaniesApiInfoData, error) {
	data, err := dst.APICompaniesApiInfo(ctx)
	if err != nil {
		return nil, err
	}

	if data.Status != "success" {
		return nil, nil
	}

	return &data.Data, nil
}
