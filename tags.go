package ctd

import (
	"context"
	"fmt"
	"strings"
)

// TagsResponse represents the response structure for the tags API.
type TagsResponse struct {
	Status  string       `json:"status"`           // Status: Status of the response
	Data    []Tag        `json:"data"`             // Data: List of tags
	Meta    MetaResponse `json:"meta"`             // Meta: Metadata about the response
	Message string       `json:"message"`          // Message: Additional message from the API
	Errors  string       `json:"errors,omitempty"` // Errors: List of errors, if any
}

// MetaResponse contains metadata about the API response.
type TagResponse struct {
	Status  string `json:"status"`           // Status: Status of the response
	Data    Tag    `json:"data"`             // Data: The tag item
	Message string `json:"message"`          // Message: Additional message from the API
	Errors  string `json:"errors,omitempty"` // Errors: List of errors, if any
}

// MetaResponse contains metadata about the API response.
type Tag struct {
	ID          int    `json:"id"`          // ID: Unique identifier for the tag
	GroupID     int    `json:"group_id"`    // GroupID: Identifier for the group the tag belongs to
	GroupName   string `json:"group_name"`  // GroupName: Name of the group the tag belongs to
	Label       string `json:"label"`       // Label: Name of the tag
	Description string `json:"description"` // Description: Description of the tag
}

// GetTags retrieves a list of tags from the Chat2Desk API.
// It uses the APIGetTags method to fetch the tags and handles errors.
// It returns a pointer to a slice of Tag, which contains the tags.
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
func (dest *Ctd) APIGetTag(ctx context.Context, id int64) (*TagResponse, error) {
	url := fmt.Sprintf("%sv1/tags/%d", dest.Url, id)
	response := TagResponse{}
	_, err := dest.doRequest(ctx, "GET", url, nil, &response)
	if err != nil {
		dest.Error(ctx, "Failed to get tag: %v", err)
		return nil, err
	}
	return &response, nil
}

// APIAssignTag assigns tags to a client or request in the Chat2Desk API.
// It constructs the API endpoint URL, prepares the payload with tag IDs, mode, and assignee ID,
// sends a POST request to the API, and returns the response data as a BasicResponse struct.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a pointer to the BasicResponse struct.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - tag_ids: A slice of tag IDs to be assigned.
//   - mode: The mode of assignment ('client' or 'request').
//   - id: The ID of the client or request to which the tags will be assigned.
//
// Returns:
//   - A pointer to a BasicResponse struct containing the response data
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) APIAssignTag(ctx context.Context, tag_ids []int64, mode string, id int64) (*BasicResponse, error) {
	if len(tag_ids) == 0 {
		return nil, ErrorInvalidParameters
	}

	url := fmt.Sprintf("%sv1/tags/assign_to", dest.Url)
	if mode != "client" {
		mode = "request"
	}
	payload := map[string]any{
		"tag_ids":       tag_ids,
		"assignee_type": mode,
		"assignee_id":   id,
	}
	response := BasicResponse{}
	data, err := dest.doRequest(ctx, "POST", url, payload, &response)
	if err != nil {
		dest.Error(ctx, "Failed to assign tag: %v", err)
		return nil, err
	}

	if strings.Contains(string(data), "request does not belong") {
		return nil, ErrorInvalidRequestID
	}

	if strings.Contains(string(data), "client does not belong") {
		return nil, ErrorInvalidClientID
	}

	return &response, nil
}

// APIRemoveTagFrom removes a tag from a client or request in the Chat2Desk API.
// It constructs the API endpoint URL, prepares the payload with mode and assignee ID,
// sends a POST request to the API, and returns the response data as a BasicResponse struct.
// If an error occurs during the request, it logs the error and returns it.
// If the request is successful, it returns a pointer to the BasicResponse struct.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - tag_id: The ID of the tag to be removed.
//   - mode: The mode of removal ('client' or 'request').
//   - id: The ID of the client or request from which the tag will be removed.
//
// Returns:
//   - A pointer to a BasicResponse struct containing the response data
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) APIRemoveTagFrom(ctx context.Context, tag_id int64, mode string, id int64) (*BasicResponse, error) {
	url := fmt.Sprintf("%sv1/tags/%d/delete_from", dest.Url, tag_id)
	payload := map[string]int64{}
	if mode == "client" {
		payload["client_id"] = id
	} else {
		payload["request_id"] = id
	}
	response := BasicResponse{}
	data, err := dest.doRequest(ctx, "DELETE", url, payload, &response)
	if err != nil {
		dest.Error(ctx, "Failed to remove tag: %v", err)
		return nil, err
	}

	if strings.Contains(string(data), "tag does not exist") {
		return nil, ErrorInvalidTagID
	}

	if strings.Contains(string(data), "request not found") {
		return nil, ErrorInvalidRequestID
	}

	if strings.Contains(string(data), "client not found") {
		return nil, ErrorInvalidClientID
	}

	return &response, nil
}

// GetTags retrieves a list of tags from the Chat2Desk API.
// It uses the APIGetTags method to fetch the tags and handles errors.
// If the response status is not "success", it returns nil.
// It returns a slice of Tag, which contains the tags.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - offset: The offset for pagination.
//   - limit: The maximum number of tags to retrieve.
//
// Returns:
//   - A slice of Tag, which contains the tags
//   - The total number of tags available.
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) GetTags(ctx context.Context, offset, limit int) ([]Tag, int, error) {
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
// It returns a pointer to a Tag, which contains the tag data.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - id: The ID of the tag to retrieve.
//
// Returns:
//   - A pointer to a Tag, which contains the tag data
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) GetTag(ctx context.Context, id int64) (*Tag, error) {
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
// It returns a slice of Tag, which contains all the tags.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//
// Returns:
//   - A slice of Tag, which contains all the tags.
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) GetAllTags(ctx context.Context) ([]Tag, error) {
	var tags []Tag
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

// AddTagToRequest assigns tags to a specific request in the Chat2Desk API.
// It uses the APIAssignTag method to assign the tags and handles errors.
// If the response status is not "success", it returns an error.
// It returns nil if the tags are successfully assigned.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - tag_ids: A slice of tag IDs to be assigned.
//   - id: The ID of the request to which the tags will be assigned.
//
// Returns:
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) AddTagToRequest(ctx context.Context, tag_ids []int64, id int64) error {
	response, err := dest.APIAssignTag(ctx, tag_ids, "request", id)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return ErrorInvalidResponse
	}

	return nil
}

// AddTagToClient assigns tags to a specific client in the Chat2Desk API.
// It uses the APIAssignTag method to assign the tags and handles errors.
// If the response status is not "success", it returns an error.
// It returns nil if the tags are successfully assigned.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - tag_ids: A slice of tag IDs to be assigned.
//   - id: The ID of the client to which the tags will be assigned.
//
// Returns:
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) AddTagToClient(ctx context.Context, tag_ids []int64, id int64) error {
	response, err := dest.APIAssignTag(ctx, tag_ids, "client", id)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return ErrorInvalidResponse
	}

	return nil
}

// RemoveTagFromRequest removes a tag from a specific request in the Chat2Desk API.
// It uses the APIRemoveTagFrom method to remove the tag and handles errors.
// If the response status is not "success", it returns an error.
// It returns nil if the tag is successfully removed.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - tag_id: The ID of the tag to be removed.
//   - id: The ID of the request from which the tag will be removed.
//
// Returns:
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) RemoveTagFromRequest(ctx context.Context, tag_id int64, id int64) error {
	response, err := dest.APIRemoveTagFrom(ctx, tag_id, "request", id)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return ErrorInvalidResponse
	}

	return nil
}

// RemoveTagFromClient removes a tag from a specific client in the Chat2Desk API.
// It uses the APIRemoveTagFrom method to remove the tag and handles errors.
// If the response status is not "success", it returns an error.
// It returns nil if the tag is successfully removed.
//
// Parameters:
//   - ctx: The context for the request, allowing for cancellation and timeouts.
//   - tag_id: The ID of the tag to be removed.
//   - id: The ID of the client from which the tag will be removed.
//
// Returns:
//   - An error if the request fails or if the response is invalid.
func (dest *Ctd) RemoveTagFromClient(ctx context.Context, tag_id int64, id int64) error {
	response, err := dest.APIRemoveTagFrom(ctx, tag_id, "client", id)
	if err != nil {
		return err
	}

	if response.Status != "success" {
		return ErrorInvalidResponse
	}

	return nil
}
