package ctd

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCtdApi_StatisticsRating(t *testing.T) {
	ctx := context.Background()
	url, api := getCredentials(t)

	t.Run("1 correct api", func(t *testing.T) {
		ctd := Ctd{Url: url, Token: api}
		data, total, err := ctd.StatisticsRating(ctx, time.Now().Add(-24*time.Hour), 0, 10)
		require.NoError(t, err, "ctd.StatisticsRating(ctx, time.Now(), 0, 10)")
		require.NotNil(t, data, "ctd.StatisticsRating(ctx, time.Now(), 0, 10)")
		require.GreaterOrEqual(t, total, 0, "total should be greater than or equal to zero")
	})

	t.Run("2 incorrect api", func(t *testing.T) {
		ctd := Ctd{Url: url, Token: "incorrect"}
		_, total, err := ctd.StatisticsRating(ctx, time.Now().Add(-24*time.Hour), 0, 10)
		require.ErrorIs(t, err, ErrorInvalidToken, "ctd.StatisticsRating(ctx, time.Now(), 0, 10)")
		require.Zero(t, total, "total should be zero on error")
	})
}
