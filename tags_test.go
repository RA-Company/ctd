package ctd

import (
	"context"
	"testing"

	"github.com/ra-company/env"
	"github.com/stretchr/testify/require"
)

func TestCtd_Tags(t *testing.T) {
	ctx := context.Background()

	url, token := getCredentials(t)

	tagID := int64(env.GetEnvInt("API_TAG_ID", 0))
	require.NotEqual(t, 0, tagID, "API_TAG_ID must be set in .env file or .settings")

	tests := []struct {
		name   string
		token  string
		tsg_id int64
		isData bool
		error  error
	}{
		{
			name:   "Incorrect token",
			token:  "incorrect_token",
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name:   "Correct token",
			token:  token,
			tsg_id: tagID,
			isData: true,
			error:  nil,
		},
		{
			name:   "Tag not found",
			token:  token,
			tsg_id: 0, // Assuming 0 is an invalid ID for testing
			isData: false,
			error:  ErrorInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(url, tt.token)
			got, err := dst.GetTag(ctx, tt.tsg_id)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.GetTag() error")

			} else {

			}
			if tt.isData {
				require.NotNil(t, got, "dst.GetTag() should return data")
			} else {
				require.Nil(t, got, "dst.GetTag() should return nil data on error")
			}

		})
	}

	t.Run("01 GetTagsList", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)

		offset := 0
		limit := 10
		total := 0
		var (
			got []Tag
			err error
		)

		t.Run("GetTagsList first page", func(t *testing.T) {
			got, total, err = dst.GetTags(ctx, offset, limit)
			require.NoError(t, err, "dst.GetTags() should not return an error")
			require.NotNil(t, got, "dst.GetTags() should return data")
			require.Greater(t, len(got), 0, "dst.GetTags() should return non-empty tag list")
			require.Greater(t, total, 0, "dst.GetTags() should return non-zero total count")
		})

		offset += limit
		t.Run("GetTagsList second page", func(t *testing.T) {
			got2, total2, err := dst.GetTags(ctx, offset, limit)
			require.NoError(t, err, "dst.GetTags() should not return an error on second call")
			require.NotNil(t, got2, "dst.GetTags() should return data on second call")
			require.Greater(t, len(got2), 0, "dst.GetTags() should return non-empty tag list on second call")
			require.Greater(t, total2, 0, "dst.GetTags() should return non-zero total count on second call")
			require.Equal(t, total, total2, "dst.GetTags() should return the same total count on both calls")
			require.NotEqual(t, got, got2, "dst.GetTags() should return different data on different calls with different offsets")
		})

		t.Run("GetTagsList with invalid offset", func(t *testing.T) {
			got, _, err = dst.GetTags(ctx, 100000, limit)
			require.NoError(t, err, "dst.GetTags() should return an error for invalid offset")
			require.Equal(t, len(got), 0, "dst.GetTags() should return nil data for invalid offset")
		})

		t.Run("GetAllTags", func(t *testing.T) {
			allTags, err := dst.GetAllTags(ctx)
			require.NoError(t, err, "dst.GetAllTags() should not return an error")
			require.NotNil(t, allTags, "dst.GetAllTags() should return data")
			require.Greater(t, len(allTags), 0, "dst.GetAllTags() should return non-empty tag list")
		})
	})

	requestID := env.GetEnvInt("API_REQUEST_ID", 0)
	require.NotEqual(t, 0, requestID, "API_REQUEST_ID must be set in .env file or .settings")

	t.Run("02 Assign tags to request", func(t *testing.T) {
		dst := &Ctd{}
		t.Run("Invalid token", func(t *testing.T) {
			dst.Init(url, "invalid token")
			err := dst.AddTagToRequest(ctx, []int64{tagID}, int64(requestID))
			require.ErrorIs(t, err, ErrorInvalidToken, "dst.AddTagToRequest() should return an error for invalid token")
		})

		t.Run("Valid token", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToRequest(ctx, []int64{tagID}, int64(requestID))
			require.NoError(t, err, "dst.AddTagToRequest() should not return an error for valid token and data")
		})

		t.Run("Set same tags", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToRequest(ctx, []int64{tagID}, int64(requestID))
			require.NoError(t, err, "dst.AddTagToRequest() should not return an error when setting the same tags again")
		})

		t.Run("Empty tag ids", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToRequest(ctx, []int64{}, int64(requestID))
			require.ErrorIs(t, err, ErrorInvalidParameters, "dst.AddTagToRequest() should return an error for empty tag IDs")
		})

		t.Run("Invalid tag IDs", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToRequest(ctx, []int64{0, -1}, int64(requestID))
			require.NoError(t, err, "dst.AddTagToRequest() should not return an error for valid invalid tag ids")
		})

		t.Run("Invalid request ID", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToRequest(ctx, []int64{tagID}, 0)
			require.ErrorIs(t, err, ErrorInvalidRequestID, "dst.AddTagToRequest() should return an error for invalid request ID")
		})
	})

	t.Run("03 Remove tag from request", func(t *testing.T) {
		dst := &Ctd{}

		t.Run("01 Invalid token", func(t *testing.T) {
			dst.Init(url, "invalid token")
			err := dst.RemoveTagFromRequest(ctx, tagID, int64(requestID))
			require.ErrorIs(t, err, ErrorInvalidToken, "dst.RemoveTagFromRequest() should return an error for invalid token")
		})

		t.Run("02 Valid token", func(t *testing.T) {
			dst.Init(url, token)
			dst.AddTagToRequest(ctx, []int64{tagID}, int64(requestID))
			err := dst.RemoveTagFromRequest(ctx, tagID, int64(requestID))
			require.NoError(t, err, "dst.RemoveTagFromRequest() should not return an error for valid token and data")
		})

		t.Run("03 Remove same tag again", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.RemoveTagFromRequest(ctx, tagID, int64(requestID))
			require.ErrorIs(t, err, ErrorInvalidRequestID, "dst.RemoveTagFromRequest() should return an error when removing a tag that is not assigned")
		})

		t.Run("04 Invalid tag ID", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.RemoveTagFromRequest(ctx, 0, int64(requestID))
			require.ErrorIs(t, err, ErrorInvalidTagID, "dst.RemoveTagFromRequest() should return an error for invalid tag ID")
		})

		t.Run("05 Invalid request ID", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.RemoveTagFromRequest(ctx, tagID, 0)
			require.ErrorIs(t, err, ErrorInvalidRequestID, "dst.RemoveTagFromRequest() should return an error for invalid request ID")
		})
	})

	clientID := env.GetEnvInt("API_CLIENT_ID", 0)
	require.NotEqual(t, 0, clientID, "API_CLIENT_ID must be set in .env file or .settings")

	t.Run("04 Assign tags to client", func(t *testing.T) {
		dst := &Ctd{}
		t.Run("Invalid token", func(t *testing.T) {
			dst.Init(url, "invalid token")
			err := dst.AddTagToClient(ctx, []int64{tagID}, int64(clientID))
			require.ErrorIs(t, err, ErrorInvalidToken, "dst.AddTagToClient() should return an error for invalid token")
		})

		t.Run("Valid token", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToClient(ctx, []int64{tagID}, int64(clientID))
			require.NoError(t, err, "dst.AddTagToClient() should not return an error for valid token and data")
		})

		t.Run("Set same tags", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToClient(ctx, []int64{tagID}, int64(clientID))
			require.NoError(t, err, "dst.AddTagToClient() should not return an error when setting the same tags again")
		})

		t.Run("Empty tag ids", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToClient(ctx, []int64{}, int64(clientID))
			require.ErrorIs(t, err, ErrorInvalidParameters, "dst.AddTagToClient() should return an error for empty tag IDs")
		})

		t.Run("Invalid tag IDs", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToClient(ctx, []int64{0, -1}, int64(clientID))
			require.NoError(t, err, "dst.AddTagToClient() should not return an error for valid invalid tag ids")
		})

		t.Run("Invalid client ID", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.AddTagToClient(ctx, []int64{tagID}, 0)
			require.ErrorIs(t, err, ErrorInvalidClientID, "dst.AddTagToClient() should return an error for invalid client ID")
		})
	})

	t.Run("05 Remove tag from client", func(t *testing.T) {
		dst := &Ctd{}

		t.Run("01 Invalid token", func(t *testing.T) {
			dst.Init(url, "invalid token")
			err := dst.RemoveTagFromClient(ctx, tagID, int64(clientID))
			require.ErrorIs(t, err, ErrorInvalidToken, "dst.RemoveTagFromClient() should return an error for invalid token")
		})

		t.Run("02 Valid token", func(t *testing.T) {
			dst.Init(url, token)
			dst.AddTagToClient(ctx, []int64{tagID}, int64(clientID))
			err := dst.RemoveTagFromClient(ctx, tagID, int64(clientID))
			require.NoError(t, err, "dst.RemoveTagFromClient() should not return an error for valid token and data")
		})

		t.Run("03 Remove same tag again", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.RemoveTagFromClient(ctx, tagID, int64(clientID))
			require.ErrorIs(t, err, ErrorInvalidClientID, "dst.RemoveTagFromClient() should return an error when removing a tag that is not assigned")
		})

		t.Run("04 Invalid tag ID", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.RemoveTagFromClient(ctx, 0, int64(clientID))
			require.ErrorIs(t, err, ErrorInvalidTagID, "dst.RemoveTagFromClient() should return an error for invalid tag ID")
		})

		t.Run("05 Invalid client ID", func(t *testing.T) {
			dst.Init(url, token)
			err := dst.RemoveTagFromClient(ctx, tagID, 0)
			require.ErrorIs(t, err, ErrorInvalidClientID, "dst.RemoveTagFromClient() should return an error for invalid client ID")
		})
	})
}
