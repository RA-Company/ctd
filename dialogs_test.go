package ctd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtd_Dialogs(t *testing.T) {
	ctx := context.Background()
	//faker := gofakeit.New(0)

	url, token := getCredentials(t)

	t.Run("01 GetDialogs Incorrect token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, "incorrect_token")
		params := &GetDialogsParams{
			Limit:  10,
			Offset: 0,
			State:  "any",
			Order:  "desc",
		}
		got, err := dst.GetDialogs(ctx, params)
		require.ErrorIs(t, err, ErrorInvalidToken, "dst.APIGetDialogs() error")
		require.Nil(t, got, "dst.APIGetDialogs() should return nil data on error")
	})

	var dialog_id = int64(0)
	t.Run("02 GetDialogs Correct token", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		params := &GetDialogsParams{
			Limit:  10,
			Offset: 0,
			State:  "any",
			Order:  "desc",
		}
		got, err := dst.GetDialogs(ctx, params)
		require.NoError(t, err, "dst.APIGetDialogs() error")
		require.NotNil(t, got, "dst.APIGetDialogs() should return data")
		require.Greater(t, len(*got), 0, "dst.APIGetDialogs() should return some dialogs")
		dialog_id = (*got)[0].ID
	})

	t.Run("03 GetDialog Incorrect ID", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		got, err := dst.GetDialog(ctx, 0)
		require.ErrorIs(t, err, ErrorInvalidID, "dst.APIGetDialog() error")
		require.Nil(t, got, "dst.APIGetDialog() should return nil data on error")
	})

	t.Run("04 GetDialog Correct ID", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		got, err := dst.GetDialog(ctx, dialog_id)
		require.NoError(t, err, "dst.GetDialog() error")
		require.NotNil(t, got, "dst.GetDialog() should return data")
		require.Equal(t, dialog_id, got.ID, "dst.GetDialog() should return correct dialog ID")
	})
}
