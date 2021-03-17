package contract

import (
	"reflect"
	"testing"
	"tezosign/types"
)

func Test_BuildVestingContractStorage(t *testing.T) {
	type args struct {
		vestingAddress types.Address
		delegateAdmin  types.Address
		timestamp      uint64
		secondsPerTick uint64
		tokensPerTick  uint64
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
				vestingAddress: "tz1NkT6YCFS3mDo6kfaMFKFrRiA7w2o5dkWp",
				delegateAdmin:  "tz1NkT6YCFS3mDo6kfaMFKFrRiA7w2o5dkWp",
				timestamp:      123,
				secondsPerTick: 10,
				tokensPerTick:  10,
			},
			expResult: `{"args":[{"args":[{"bytes":"0000221f3e8f57fccf16203bbd5f27590d365b190084"},{"bytes":"0000221f3e8f57fccf16203bbd5f27590d365b190084"}],"prim":"Pair"},{"args":[{"int":"0"},{"args":[{"int":"123"},{"args":[{"int":"10"},{"int":"10"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Pair"}`,
			wantErr:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := BuildVestingContractStorage(test.args.vestingAddress, test.args.delegateAdmin, test.args.timestamp, test.args.secondsPerTick, test.args.tokensPerTick)
			if test.wantErr != (gotErr != nil) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}

			if !test.wantErr && !reflect.DeepEqual(string(got), test.expResult) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}
		})
	}
}
