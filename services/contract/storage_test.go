package contract

import (
	"msig/types"
	"reflect"
	"testing"
)

func Test_BuildContractStorage(t *testing.T) {
	type args struct {
		threshold uint
		pubKeys   []types.PubKey
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
				threshold: 1,
				pubKeys:   []types.PubKey{"edpkuNVuqdPhCsrYqkq21qW2hYTSZWMjQQjfyogoPZ2AfqCmonziNh", "p2pk64iwFyjuvy1SYwkMXeM5GwYGdqQZPwwBViGvhkqM7nGyEwgjpM7", "sppk7d8CHGV9SCVDi9ciUVAyGTSLExWRSBAJN4vcFpqWEYbWf9ZNr8D"},
			},
			expResult: `{"args":[{"int":"0"},{"args":[{"int":"1"},[{"bytes":"005ffdd5422addf020a689a1660e1e8c5a0247ed5bfd7ea4f4194b1a2d9f8129cb"},{"bytes":"020213ebf302f60ddcc2168c3d5b2e1f9a9bfef1325682610e1578eecd0ea0846d74"},{"bytes":"0103f713b3d4447a11d5de2c190a67a1164f85b1b265a02331e2b24aee6afbacf286"}]],"prim":"Pair"}],"prim":"Pair"}`,
			wantErr:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := BuildContractStorage(test.args.threshold, test.args.pubKeys)
			if test.wantErr != (gotErr != nil) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}

			if !test.wantErr && !reflect.DeepEqual(string(got), test.expResult) {
				t.Errorf("wantErr: %t | results %s == %s | err: %v", test.wantErr, got, test.expResult, gotErr)
			}
		})
	}
}
