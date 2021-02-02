package contract

import (
	"msig/models"
	"msig/types"
	"reflect"
	"testing"
)

//TODO Add negative cases
func Test_BuildContractSignPayload(t *testing.T) {
	type args struct {
		networkID       string
		counter         int64
		operationParams models.ContractOperationRequest
	}

	testCases := []struct {
		name      string
		args      args
		expResult types.Payload
		wantErr   bool
	}{
		{
			name: "transaction",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				operationParams: models.ContractOperationRequest{
					ContractID: "KT1LAuGLiaCF9A72qZtFvVhyzzNFg86fwFnV",
					Type:       models.Transfer,
					To:         "tz1dBT7PKeSDbPK1No7KNhTvrr3XoLe8vKLH",
					Amount:     1010,
				},
			},

			//Left Left Left Left
			expResult: "05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d0007070000050505050505050507070a000000160000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d000b20f",
			wantErr:   false,
		},
		{
			name: "delegation",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				operationParams: models.ContractOperationRequest{
					ContractID: "KT1LAuGLiaCF9A72qZtFvVhyzzNFg86fwFnV",
					Type:       models.Delegation,
					To:         "tz3Mo3gHekQhCmykfnC58ecqJLXrjMKzkF2Q",
				},
			},
			//Left Left Left Right
			expResult: "05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d0007070000050505050505050805090a0000001502101368afffeb1dc3c089facbbe23f5c30b787ce9",
			wantErr:   false,
		},
		{
			name: "delegation tz2",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				operationParams: models.ContractOperationRequest{
					ContractID: "KT1LAuGLiaCF9A72qZtFvVhyzzNFg86fwFnV",
					Type:       models.Delegation,
					To:         "tz29nEixktH9p9XTFX7p8hATUyeLxXEz96KR",
				},
			},

			expResult: "05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d0007070000050505050505050805090a0000001501101368afffeb1dc3c089facbbe23f5c30b787ce9",
			wantErr:   false,
		},
		{
			name: "none delegation",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				operationParams: models.ContractOperationRequest{
					ContractID: "KT1LAuGLiaCF9A72qZtFvVhyzzNFg86fwFnV",
					Type:       models.Delegation,
					To:         "",
				},
			},

			expResult: "05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d000707000005050505050505080306",
			wantErr:   false,
		},
		{
			name: "fa transfer with default From",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				operationParams: models.ContractOperationRequest{
					ContractID: "KT1LAuGLiaCF9A72qZtFvVhyzzNFg86fwFnV",
					AssetID:    "KT1NtGnEjacAkBph7k9HWVrN38PoYjcXTxdY",
					Type:       models.FATransfer,
					To:         "tz1dBT7PKeSDbPK1No7KNhTvrr3XoLe8vKLH",
					Amount:     110,
				},
			},

			expResult: "05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d000707000005050505050807070a00000016019ce13845659ff2582555ec08dc322007f6493e800007070a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d0007070a000000160000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d000ae01",
			wantErr:   false,
		},
		{
			name: "fa transfer with custom From",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				operationParams: models.ContractOperationRequest{
					ContractID: "KT1LAuGLiaCF9A72qZtFvVhyzzNFg86fwFnV",
					AssetID:    "KT1NtGnEjacAkBph7k9HWVrN38PoYjcXTxdY",
					Type:       models.FATransfer,
					To:         "tz1dBT7PKeSDbPK1No7KNhTvrr3XoLe8vKLH",
					From:       "tz29nEixktH9p9XTFX7p8hATUyeLxXEz96KR",
					Amount:     110,
				},
			},

			expResult: "05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d000707000005050505050807070a00000016019ce13845659ff2582555ec08dc322007f6493e800007070a000000160001101368afffeb1dc3c089facbbe23f5c30b787ce907070a000000160000c06b6aa5308a9a89a628ebb8234d5055bf9ba1d000ae01",
			wantErr:   false,
		},
		{
			name: "storage update",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				operationParams: models.ContractOperationRequest{
					ContractID: "KT1LAuGLiaCF9A72qZtFvVhyzzNFg86fwFnV",
					Type:       models.StorageUpdate,
					Threshold:  1,
					Keys:       []types.PubKey{"edpkuNVuqdPhCsrYqkq21qW2hYTSZWMjQQjfyogoPZ2AfqCmonziNh", "p2pk64iwFyjuvy1SYwkMXeM5GwYGdqQZPwwBViGvhkqM7nGyEwgjpM7", "sppk7d8CHGV9SCVDi9ciUVAyGTSLExWRSBAJN4vcFpqWEYbWf9ZNr8D"},
				},
			},

			expResult: "05070707070a000000049caecab90a00000016017f1df41f643db8039663fd5eb3b025e07efbaf3d000707000005080707000102000000740a00000021005ffdd5422addf020a689a1660e1e8c5a0247ed5bfd7ea4f4194b1a2d9f8129cb0a00000022020213ebf302f60ddcc2168c3d5b2e1f9a9bfef1325682610e1578eecd0ea0846d740a000000220103f713b3d4447a11d5de2c190a67a1164f85b1b265a02331e2b24aee6afbacf286",
			wantErr:   false,
		},
		{
			name: "custom payload",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				counter:   1,
				operationParams: models.ContractOperationRequest{
					ContractID:    "KT1WKnsxYnYpTfgCZDuJ9mmv7f6c38Aea9wF",
					Type:          models.CustomPayload,
					CustomPayload: `[{"prim":"PUSH","args":[{"prim": "string"},{"string":"Was inserted"}]},{"prim":"FAILWITH"}]`,
				},
			},

			expResult: "05070707070a000000049caecab90a0000001601ee7d9fd9c644b230e516ffb8ef2bcb9629bb1a37000707000105050508020000001707430368010000000c57617320696e7365727465640327",
			wantErr:   false,
		},
		{
			name: "custom payload hex",
			args: args{
				networkID: "NetXjD3HPJJjmcd",
				counter:   1,
				operationParams: models.ContractOperationRequest{
					ContractID:    "KT1WKnsxYnYpTfgCZDuJ9mmv7f6c38Aea9wF",
					Type:          models.CustomPayload,
					CustomPayload: "0x05020000001707430368010000000c57617320696e7365727465640327",
				},
			},

			expResult: "05070707070a000000049caecab90a0000001601ee7d9fd9c644b230e516ffb8ef2bcb9629bb1a37000707000105050508020000001707430368010000000c57617320696e7365727465640327",
			wantErr:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := BuildContractSignPayload(test.args.networkID, test.args.counter, test.args.operationParams)
			if test.wantErr != (gotErr != nil) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}
			if !test.wantErr && !reflect.DeepEqual(got, test.expResult) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}
		})
	}
}
