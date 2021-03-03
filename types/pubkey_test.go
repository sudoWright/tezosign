package types

import (
	"testing"
)

func Test_PubKeyAddress(t *testing.T) {
	type args struct {
		pubKey  PubKey
		address string
	}

	testCases := []struct {
		name      string
		args      args
		expResult string
		wantErr   bool
	}{
		{
			name: "Success case",
			args: args{
				pubKey:  "edpkuEZ8FpnCWY2mNUpNYaF4GH3zZuCYKoNjZJPvtoKEki2ZfbFPbS",
				address: "tz1boE6s8tS3pcxetHuaAPZWzHicMa39jSfj",
			},
			expResult: "",
			wantErr:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			address, gotErr := test.args.pubKey.Address()
			if address.String() != test.args.address || test.wantErr != (gotErr != nil) {
				t.Errorf("wantErr: %t | err: %v", test.wantErr, gotErr)
			}
		})
	}
}
