package ctd

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtdApi_Operators(t *testing.T) {
	ctx := context.Background()

	url, token := getCredentials(t)

	type args struct {
		url   string
		token string
	}
	tests := []struct {
		name   string
		args   args
		isData bool
		error  error
	}{
		{
			name: "Incorrect token",
			args: args{
				url:   url,
				token: "incorrect_token",
			},
			isData: false,
			error:  ErrorInvalidToken,
		},
		{
			name: "Correct token",
			args: args{
				url:   url,
				token: token,
			},
			isData: true,
			error:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &Ctd{}
			dst.Init(tt.args.url, tt.args.token)
			got, total, err := dst.Operators(ctx, 0, 10)
			if tt.error != nil {
				require.ErrorIs(t, err, tt.error, "dst.Operators() error")
				require.Zero(t, total, "dst.Operators() should return zero total on error")
			} else {
				require.NoError(t, err, "dst.Operators() error")
				require.Greater(t, total, 0, "dst.Operators() should return some total")
			}
			if tt.isData {
				require.NotNil(t, got, "dst.Operators() should return data")
			} else {
				require.Nil(t, got, "dst.Operators() should return nil data on error")
			}

		})
	}

	t.Run("OperatorsStatuses", func(t *testing.T) {
		dst := &Ctd{}
		dst.Init(url, token)
		got, err := dst.APIOperatorStatuses(ctx)
		require.NoError(t, err, "dst.OperatorsStatuses() error")
		require.NotNil(t, got, "dst.OperatorsStatuses() should return data")
	})
}

func TestOperator_GetLegacyRole(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{
			name: "Admin role by AccessRightName",
			json: `{"id": 617,"email": "test@test.com","phone": null,"role": "admin","access_right_id": null,"online": 1,"offline_type": null,"avatar": null,"opened_dialogs": 0,"first_name": "John","last_name": "Doe","external_id": null,"last_visit": "2026-01-29T06:14:41 UTC","status_id": 0,"two_factor": false,"status": "admin","access_right_name": null}`,
			want: "admin",
		},
		{
			name: "Supervisor role by AccessRightName",
			json: `{"id":311416,"email":"test@test.com","phone":"","role":"operator","access_right_id":54521,"online":1,"offline_type":"busy","avatar":null,"opened_dialogs":0,"first_name":"John","last_name":"Doe","external_id":"","last_visit":"2024-11-13T10:55:52 UTC","status_id":1,"two_factor":false,"status":"enabled","access_right_name":"Супервайзер"}`,
			want: "supervisor",
		},
		{
			name: "Deleted role by AccessRightName",
			json: `{"id":313753,"email":"test@test.com","phone":"+7 701 763 0465","role":"deleted","access_right_id":null,"online":0,"offline_type":null,"avatar":null,"opened_dialogs":0,"first_name":"John","last_name":"Doe","external_id":"","last_visit":null,"status_id":0,"two_factor":false,"status":"deleted","access_right_name":null}`,
			want: "deleted",
		},
		{
			name: "Disabled role by AccessRightName",
			json: `{"id":313787,"email":"test@test.com","phone":"","role":"disabled","access_right_id":null,"online":1,"offline_type":null,"avatar":null,"opened_dialogs":0,"first_name":"John","last_name":"Doe","external_id":"","last_visit":"2025-03-05T12:07:02 UTC","status_id":0,"two_factor":true,"status":"disabled","access_right_name":null}`,
			want: "disabled",
		},
		{
			name: "Operator role by AccessRightName",
			json: `{"id":311416,"email":"test@test.com","phone":"","role":"operator","access_right_id":54521,"online":1,"offline_type":"busy","avatar":null,"opened_dialogs":0,"first_name":"John","last_name":"Doe","external_id":"","last_visit":"2024-11-13T10:55:52 UTC","status_id":1,"two_factor":false,"status":"enabled","access_right_name":"Operators"}`,
			want: "operator",
		},
		{
			name: "Supervisor role by Role field",
			json: `{"id":962,"email":"test@test.com","first_name":"John","last_name":"Doe","role":"supervisor","phone":"+7 700 390 1850","avatar":"https://test.com/companies/company_479/users/avatars/user5e4e802f76de7.jpg","last_visit":"2026-01-29T07:38:48 UTC","online":1,"offline_type":null,"external_id":null,"opened_dialogs":0,"status_id":0,"two_factor":true}`,
			want: "supervisor",
		},
		{
			name: "Admin role by Role field",
			json: `{"id":617,"email":"test@test.com","first_name":"John","last_name":"Doe","role":"admin","phone":"+7 700 390 1850","avatar":"https://test.com/companies/company_479/users/avatars/2023-4/user42b46c51f02587d6a4ae.png","last_visit":"2026-01-29T06:59:05 UTC","online":1,"offline_type":null,"external_id":"","opened_dialogs":1,"status_id":0,"two_factor":false}`,
			want: "admin",
		},
		{
			name: "Deleted role by Role field",
			json: `{"id":638,"email":"test@test.com","first_name":"John","last_name":"Doe","role":"deleted","phone":"","avatar":null,"last_visit":null,"online":0,"offline_type":null,"external_id":null,"opened_dialogs":0,"status_id":0,"two_factor":false}`,
			want: "deleted",
		},
		{
			name: "Operator role by Role field",
			json: `{"id":618,"email":"test@test.com","first_name":"John","last_name":"Doe","role":"operator","phone":"+7 777 222 0509","avatar":"https://test.com/companies/company_479/users/avatars/424272/usera66bc5820bb7ae46d587.jpg","last_visit":"2026-01-28T13:47:10 UTC","online":0,"offline_type":null,"external_id":"","opened_dialogs":0,"status_id":0,"two_factor":true}`,
			want: "operator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operator := Operator{}
			err := json.Unmarshal([]byte(tt.json), &operator)
			require.NoError(t, err, "json.Unmarshal() error")

			if got := operator.GetLegacyRole(); got != tt.want {
				t.Errorf("Operator.GetLegacyRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
