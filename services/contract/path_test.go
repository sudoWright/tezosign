package contract

import (
	"reflect"
	"testing"
	"tezosign/models"

	"blockwatch.cc/tzindex/micheline"
)

func Test_BuildMichelsonPath(t *testing.T) {
	type args struct {
		actionType   models.ActionType
		actionParams string
	}

	testCases := []struct {
		name      string
		args      args
		expResult string
		wantErr   bool
	}{
		{
			name: "Storage",
			args: args{
				actionType:   models.StorageUpdate,
				actionParams: `{"int":"1"}`,
			},
			expResult: `{"args":[{"int":"1"}],"prim":"Right"}`,
			wantErr:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			param := &micheline.Prim{}
			_ = param.UnmarshalJSON([]byte(test.args.actionParams))
			got, gotErr := buildMichelsonPath(test.args.actionType, param)
			if test.wantErr != (gotErr != nil) {
				t.Errorf("wantErr: %t | err: %v", test.wantErr, gotErr)
			}

			path, _ := got.MarshalJSON()

			if !test.wantErr && !reflect.DeepEqual(string(path), test.expResult) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, path, test.expResult, gotErr)
			}
		})
	}
}

func Test_GetMichelsonParamsByPath(t *testing.T) {
	type args struct {
		actionType   models.ActionType
		actionParams string
	}

	testCases := []struct {
		name      string
		args      args
		expResult string
		wantErr   bool
	}{
		{
			name: "Storage",
			args: args{
				actionType:   models.StorageUpdate,
				actionParams: `{"args":[{"int":"1"}],"prim":"Right"}`,
			},
			expResult: `{"int":"1"}`,
			wantErr:   false,
		},
		{
			name: "Reject",
			args: args{
				actionType:   models.CustomPayload,
				actionParams: `{"prim":"Left","args":[{"prim":"Right","args":[[{"prim":"DROP"},{"prim":"NIL","args":[{"prim":"operation"}]}]]}]}`,
			},
			expResult: `[{"prim":"DROP"},{"args":[{"prim":"operation"}],"prim":"NIL"}]`,
			wantErr:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			param := &micheline.Prim{}
			_ = param.UnmarshalJSON([]byte(test.args.actionParams))
			got, gotErr := getMichelsonParamsByActionType(test.args.actionType, param)
			if test.wantErr != (gotErr != nil) {
				t.Errorf("wantErr: %t | err: %v", test.wantErr, gotErr)
			}

			path, _ := got.MarshalJSON()

			if !test.wantErr && !reflect.DeepEqual(string(path), test.expResult) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, path, test.expResult, gotErr)
			}
		})
	}
}
