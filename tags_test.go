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

	tagID := env.GetEnvInt("API_TAG_ID", 0)
	require.NotEqual(t, 0, tagID, "API_TAG_ID must be set in .env file or .settings")

	tests := []struct {
		name   string
		token  string
		tsg_id int
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

	t.Run("GetTagsList", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)

		offset := 0
		limit := 10
		total := 0
		var (
			got *[]TagItem
			err error
		)

		t.Run("GetTagsList first page", func(t *testing.T) {
			got, total, err = dst.GetTags(ctx, offset, limit)
			require.NoError(t, err, "dst.GetTags() should not return an error")
			require.NotNil(t, got, "dst.GetTags() should return data")
			require.Greater(t, len(*got), 0, "dst.GetTags() should return non-empty tag list")
			require.Greater(t, total, 0, "dst.GetTags() should return non-zero total count")
		})

		offset += limit
		t.Run("GetTagsList second page", func(t *testing.T) {
			got2, total2, err := dst.GetTags(ctx, offset, limit)
			require.NoError(t, err, "dst.GetTags() should not return an error on second call")
			require.NotNil(t, got2, "dst.GetTags() should return data on second call")
			require.Greater(t, len(*got2), 0, "dst.GetTags() should return non-empty tag list on second call")
			require.Greater(t, total2, 0, "dst.GetTags() should return non-zero total count on second call")
			require.Equal(t, total, total2, "dst.GetTags() should return the same total count on both calls")
			require.NotEqual(t, got, got2, "dst.GetTags() should return different data on different calls with different offsets")
		})

		t.Run("GetTagsList with invalid offset", func(t *testing.T) {
			got, _, err = dst.GetTags(ctx, 100000, limit)
			require.NoError(t, err, "dst.GetTags() should return an error for invalid offset")
			require.Equal(t, len(*got), 0, "dst.GetTags() should return nil data for invalid offset")
		})

		t.Run("GetAllTags", func(t *testing.T) {
			allTags, err := dst.GetAllTags(ctx)
			require.NoError(t, err, "dst.GetAllTags() should not return an error")
			require.NotNil(t, allTags, "dst.GetAllTags() should return data")
			require.Greater(t, len(*allTags), 0, "dst.GetAllTags() should return non-empty tag list")
		})
	})
}
