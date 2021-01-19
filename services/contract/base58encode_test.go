package contract

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func Test_EncodeBase58ToPrimBytes(t *testing.T) {
	type args struct {
		base58 string
	}

	testCases := []struct {
		name      string
		args      args
		expResult string
		wantErr   bool
	}{
		{
			name: "KT address",
			args: args{
				base58: "KT1LAuGLiaCF9A72qZtFvVhyzzNFg86fwFnV",
			},
			expResult: "017f1df41f643db8039663fd5eb3b025e07efbaf3d00",
			wantErr:   false,
		},
		{
			name: "tz1 address",
			args: args{
				base58: "tz1dBT7PKeSDbPK1No7KNhTvrr3XoLe8vKLH",
			},
			expResult: "0000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d0",
			wantErr:   false,
		},
		{
			name: "tz2 address",
			args: args{
				base58: "tz29nEixktH9p9XTFX7p8hATUyeLxXEz96KR",
			},
			expResult: "0001101368afffeb1dc3c089facbbe23f5c30b787ce9",
			wantErr:   false,
		},
		{
			name: "tz3 address",
			args: args{
				base58: "tz3Mo3gHekQhCmykfnC58ecqJLXrjMKzkF2Q",
			},
			expResult: "0002101368afffeb1dc3c089facbbe23f5c30b787ce9",
			wantErr:   false,
		},
		{
			name: "ED25519PublicKey",
			args: args{
				base58: "edpkuNVuqdPhCsrYqkq21qW2hYTSZWMjQQjfyogoPZ2AfqCmonziNh",
			},
			expResult: "005ffdd5422addf020a689a1660e1e8c5a0247ed5bfd7ea4f4194b1a2d9f8129cb",
			wantErr:   false,
		},
		{
			name: "P256PublicKey",
			args: args{
				base58: "p2pk64iwFyjuvy1SYwkMXeM5GwYGdqQZPwwBViGvhkqM7nGyEwgjpM7",
			},
			expResult: "020213ebf302f60ddcc2168c3d5b2e1f9a9bfef1325682610e1578eecd0ea0846d74",
			wantErr:   false,
		},
		{
			name: "Secp256k1PublicKey",
			args: args{
				base58: "sppk7d8CHGV9SCVDi9ciUVAyGTSLExWRSBAJN4vcFpqWEYbWf9ZNr8D",
			},
			expResult: "0103f713b3d4447a11d5de2c190a67a1164f85b1b265a02331e2b24aee6afbacf286",
			wantErr:   false,
		},
		{
			name: "edSig",
			args: args{
				base58: "edsigtwo6iJyKdGMKKFxSqVT6KvhHuJK1whHdZo4rDF5rRhxpYHiZpnpBHtLRs3BEHyfFW3C8cSCQ7Zu55Kr339cN6M8PbeiMEz",
			},
			expResult: "b75be147bbbee4c2cb4b50942453d4c7866da234142537ea70fc3859e4db9e27b731e99c5371ab1d77d6683bcff6a480449011bf52481f98096e322975238c0d",
			wantErr:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := encodeBase58ToPrimBytes(test.args.base58)
			if test.wantErr && gotErr == nil {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}

			if !test.wantErr && !reflect.DeepEqual(hex.EncodeToString(got), test.expResult) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}
		})
	}
}
