package contract

import (
	"reflect"
	"testing"
	"tezosign/types"
)

func Test_BuildFullTxPayload(t *testing.T) {
	type args struct {
		payload    types.Payload
		signatures []types.Signature
	}

	testCases := []struct {
		name      string
		args      args
		expResult string
		wantErr   bool
	}{
		{
			name: "Full tx FA transfer",
			args: args{
				payload:    "070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d000707000005050505050807070a00000016019ce13845659ff2582555ec08dc322007f6493e800007070a000000160001101368afffeb1dc3c089facbbe23f5c30b787ce907070a000000160000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d000ae01",
				signatures: []types.Signature{"edsigtwo6iJyKdGMKKFxSqVT6KvhHuJK1whHdZo4rDF5rRhxpYHiZpnpBHtLRs3BEHyfFW3C8cSCQ7Zu55Kr339cN6M8PbeiMEz"},
			},
			expResult: `{"args":[{"args":[{"int":"0"},{"args":[{"args":[{"args":[{"args":[{"bytes":"019ce13845659ff2582555ec08dc322007f6493e8000"},{"args":[{"bytes":"0001101368afffeb1dc3c089facbbe23f5c30b787ce9"},{"args":[{"bytes":"0000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d0"},{"int":"110"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Right"}],"prim":"Left"}],"prim":"Left"}],"prim":"Pair"},[{"args":[{"bytes":"b75be147bbbee4c2cb4b50942453d4c7866da234142537ea70fc3859e4db9e27b731e99c5371ab1d77d6683bcff6a480449011bf52481f98096e322975238c0d"}],"prim":"Some"}]],"prim":"Pair"}`,
			wantErr:   false,
		},
		{
			name: "Full tx FA transfer with watermark",
			args: args{
				payload:    "05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d000707000005050505050807070a00000016019ce13845659ff2582555ec08dc322007f6493e800007070a000000160001101368afffeb1dc3c089facbbe23f5c30b787ce907070a000000160000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d000ae01",
				signatures: []types.Signature{"edsigtwo6iJyKdGMKKFxSqVT6KvhHuJK1whHdZo4rDF5rRhxpYHiZpnpBHtLRs3BEHyfFW3C8cSCQ7Zu55Kr339cN6M8PbeiMEz"},
			},
			expResult: `{"args":[{"args":[{"int":"0"},{"args":[{"args":[{"args":[{"args":[{"bytes":"019ce13845659ff2582555ec08dc322007f6493e8000"},{"args":[{"bytes":"0001101368afffeb1dc3c089facbbe23f5c30b787ce9"},{"args":[{"bytes":"0000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d0"},{"int":"110"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Right"}],"prim":"Left"}],"prim":"Left"}],"prim":"Pair"},[{"args":[{"bytes":"b75be147bbbee4c2cb4b50942453d4c7866da234142537ea70fc3859e4db9e27b731e99c5371ab1d77d6683bcff6a480449011bf52481f98096e322975238c0d"}],"prim":"Some"}]],"prim":"Pair"}`,
			wantErr:   false,
		},
		{
			name: "Full tx FA transfer with watermark and prefix",
			args: args{
				payload:    "0x05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d000707000005050505050807070a00000016019ce13845659ff2582555ec08dc322007f6493e800007070a000000160001101368afffeb1dc3c089facbbe23f5c30b787ce907070a000000160000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d000ae01",
				signatures: []types.Signature{"edsigtwo6iJyKdGMKKFxSqVT6KvhHuJK1whHdZo4rDF5rRhxpYHiZpnpBHtLRs3BEHyfFW3C8cSCQ7Zu55Kr339cN6M8PbeiMEz"},
			},
			expResult: `{"args":[{"args":[{"int":"0"},{"args":[{"args":[{"args":[{"args":[{"bytes":"019ce13845659ff2582555ec08dc322007f6493e8000"},{"args":[{"bytes":"0001101368afffeb1dc3c089facbbe23f5c30b787ce9"},{"args":[{"bytes":"0000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d0"},{"int":"110"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Right"}],"prim":"Left"}],"prim":"Left"}],"prim":"Pair"},[{"args":[{"bytes":"b75be147bbbee4c2cb4b50942453d4c7866da234142537ea70fc3859e4db9e27b731e99c5371ab1d77d6683bcff6a480449011bf52481f98096e322975238c0d"}],"prim":"Some"}]],"prim":"Pair"}`,
			wantErr:   false,
		},
		{
			name: "Full tx FA transfer with prefix",
			args: args{
				payload:    "0x070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d000707000005050505050807070a00000016019ce13845659ff2582555ec08dc322007f6493e800007070a000000160001101368afffeb1dc3c089facbbe23f5c30b787ce907070a000000160000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d000ae01",
				signatures: []types.Signature{"edsigtwo6iJyKdGMKKFxSqVT6KvhHuJK1whHdZo4rDF5rRhxpYHiZpnpBHtLRs3BEHyfFW3C8cSCQ7Zu55Kr339cN6M8PbeiMEz"},
			},
			expResult: `{"args":[{"args":[{"int":"0"},{"args":[{"args":[{"args":[{"args":[{"bytes":"019ce13845659ff2582555ec08dc322007f6493e8000"},{"args":[{"bytes":"0001101368afffeb1dc3c089facbbe23f5c30b787ce9"},{"args":[{"bytes":"0000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d0"},{"int":"110"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Pair"}],"prim":"Right"}],"prim":"Left"}],"prim":"Left"}],"prim":"Pair"},[{"args":[{"bytes":"b75be147bbbee4c2cb4b50942453d4c7866da234142537ea70fc3859e4db9e27b731e99c5371ab1d77d6683bcff6a480449011bf52481f98096e322975238c0d"}],"prim":"Some"}]],"prim":"Pair"}`,
			wantErr:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, _, gotErr := BuildFullTxPayload(test.args.payload, test.args.signatures)
			if test.wantErr != (gotErr != nil) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}

			if !test.wantErr && !reflect.DeepEqual(string(got), test.expResult) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}
		})
	}
}
