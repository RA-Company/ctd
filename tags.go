package ctd

import (
	"context"
	"fmt"
	"strings"
)

// TagsResponse represents the response structure for the tags API.
type TagsResponse struct {
	Status  string       `json:"status"`           // Status: Status of the response
	Data    []TagItem    `json:"data"`             // Data: List of tags
	Meta    MetaResponse `json:"meta"`             // Meta: Metadata about the response
	Message string       `json:"message"`          // Message: Additional message from the API
	Errors  string       `json:"errors,omitempty"` // Errors: List of errors, if any
}

// MetaResponse contains metadata about the API response.
type TagResponse struct {
	Status  string  `json:"status"`           // Status: Status of the response
	Data    TagItem `json:"data"`             // Data: The tag item
	Message string  `json:"message"`          // Message: Additional message from the API
	Errors  string  `json:"errors,omitempty"` // Errors: List of errors, if any
}

// MetaResponse contains metadata about the API response.
type TagItem struct {
	ID          int    `json:"id"`          // ID: Unique identifier for the tag
	GroupID     int    `json:"group_id"`    // GroupID: Identifier for the group the tag belongs to
	GroupName   string `json:"group_name"`  // GroupName: Name of the group the tag belongs to
	Label       string `json:"label"`       // Label: Name of the tag
	Description string `json:"description"` // Description: Description of the tag
}

// GetTags retrieves a list of tags from the Chat2Desk API.
// It uses the APIGetTags method to fetch the tags and handles errors.
// It returns a pointer to a slice of TagItem, which contains the tags.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination.
//   - limit: The maximum number of tags to retrieve.
//
// Returns:
//   - A pointer to a TagsResponse struct containing the list of tags
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) APIGetTags(ctx context.Context, offset, limit int) (*TagsResponse, error) {
	url := fmt.Sprintf("%sv1/tags?offset=%d&limit=%d", dest.Url, offset, limit)

	response := TagsResponse{}
	_, err := dest.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dest.Error(ctx, "Failed to get tags: %v", err)
		return nil, err
	}
	return &response, nil
}

// APIGetTag retrieves a specific tag by its ID from the Chat2Desk API.
// It uses the doRequest method to send a GET request to the API.
// If the request fails, it logs the error and returns nil.
// It returns a pointer to a TagResponse struct containing the tag data.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the tag to retrieve.
//
// Returns:
//   - A pointer to a TagResponse struct containing the tag data
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) APIGetTag(ctx context.Context, id int) (*TagResponse, error) {
	url := fmt.Sprintf("%sv1/tags/%d", dest.Url, id)
	response := TagResponse{}
	_, err := dest.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dest.Error(ctx, "Failed to get tag: %v", err)
		return nil, err
	}
	return &response, nil
}

// GetTags retrieves a list of tags from the Chat2Desk API.
// It uses the APIGetTags method to fetch the tags and handles errors.
// If the response status is not "success", it returns nil.
// It returns a slice of TagItem, which contains the tags.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination.
//   - limit: The maximum number of tags to retrieve.
//
// Returns:
//   - A slice of TagItem, which contains the tags
//   - The total number of tags available.
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) GetTags(ctx context.Context, offset, limit int) ([]TagItem, int, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 10 // Default limit if not specified
	}
	if limit > 200 {
		limit = 200 // Maximum limit enforced by the API
	}
	response, err := dest.APIGetTags(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	if response.Status != "success" {
		return nil, 0, nil
	}

	return response.Data, response.Meta.Total, nil
}

// GetTag retrieves a specific tag by its ID from the Chat2Desk API.
// It uses the APIGetTag method to fetch the tag and handles errors.
// If the response status is not "success", it returns nil.
// It returns a pointer to a TagItem, which contains the tag data.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the tag to retrieve.
//
// Returns:
//   - A pointer to a TagItem, which contains the tag data
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) GetTag(ctx context.Context, id int) (*TagItem, error) {
	response, err := dest.APIGetTag(ctx, id)
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

// GetAllTags retrieves all tags from the Chat2Desk API.
// It uses the GetTags method to fetch tags in a loop until all tags are retrieved.
// It returns a slice of TagItem, which contains all the tags.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A slice of TagItem, which contains all the tags.
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) GetAllTags(ctx context.Context) ([]TagItem, error) {
	var tags []TagItem
	offset := 0
	limit := 200

	for {
		response, total, err := dest.GetTags(ctx, offset, limit)
		if err != nil {
			return nil, err
		}

		if len(response) == 0 {
			break
		}

		tags = append(tags, response...)
		offset += limit
		if offset >= total {
			break
		}
	}

	return tags, nil
}
