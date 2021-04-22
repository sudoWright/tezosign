package services

import (
	"encoding/hex"
	"testing"
	"tezosign/types"
)

func Test_verifysign(t *testing.T) {
	type args struct {
		payload   string
		pubKey    types.PubKey
		signature types.Signature
	}

	testCases := []struct {
		name      string
		args      args
		expResult string
		wantErr   bool
	}{
		{
			name: "Ledger raw",
			args: args{
				payload:   "616e7920737472696e6720746861742077696c6c206265207369676e6564",
				pubKey:    "edpkv13wgJVsEQGiQmw6M2gt9SCu55ajuZDiS9Xyxq375tBUtv8Fjh",
				signature: "edsigu7CP5oCEFGGJC1ixnUK85NkZ449coGeFpitcvC267r81jwoinkX5DKSrfemwJS97bmejSnkm5Nrfdmrsqpri4VCQrEs3Bz",
			},
			expResult: "",
			wantErr:   false,
		},
		{
			name: "Ed25519 signature",
			args: args{
				payload:   "05070707070a00000004a83650210a00000016019ce13845659ff2582555ec08dc322007f6493e800007070000050505050505050507070a00000016000032bb7d0084f79711f757d66b791d5290f88eb28000a80f",
				pubKey:    "edpkvGDDYVjo8sz2dJD9mD4ufTVgvzRzZLj4PYApw9JFPLtQ9uDwQ8",
				signature: "sigotZvkArG58FivwJFeVetWvJk2GjWWCBHHwEKqmCvebhRcGmfKQKsTEv9gC779ZiZCNwSYmKf7MvDPBgwJH5m1SEZEXX6a",
			},
			expResult: "",
			wantErr:   false,
		},
		{
			name: "Ledger signature",
			args: args{
				payload:   "05070707070a00000004a83650210a00000016019ce13845659ff2582555ec08dc322007f6493e800007070000050505050505050507070a00000016000032bb7d0084f79711f757d66b791d5290f88eb28000a80f",
				pubKey:    "edpkuEZ8FpnCWY2mNUpNYaF4GH3zZuCYKoNjZJPvtoKEki2ZfbFPbS",
				signature: "edsigtYjH3GbyEQodeauFpWYq8Dfmpk5iyL4J2EFFFiZ3GMU5RY5R1Hj2hSZ1ZUqGQvueHqB98hk4dyhARWTLLFree4opMoU3Hb",
			},
			expResult: "",
			wantErr:   false,
		},
		//TODO fix tests
		//{
		//	name: "Secp256k1 signature",
		//	args: args{
		//		payload:   "05070707070a00000004a83650210a00000016019ce13845659ff2582555ec08dc322007f6493e800007070000050505050505050507070a00000016000032bb7d0084f79711f757d66b791d5290f88eb28000a80f",
		//		pubKey:    "sppk7bfFYv6qG9NwDM1k9x7RVJCQexGkU15WtVqSWMCzJxpwaCbtCWV",
		//		signature: "spsig1E1YVJEE9HpA9M5trj4Vaszx2a2FNn8yqenbsd7BP65si738yayrQVGfkuN5oy8VWBoH4iK6B4bLfbbiLxjSPz5WhGgUPZ",
		//	},
		//	expResult: "",
		//	wantErr:   false,
		//},
		//{
		//	name: "P256 signature",
		//	args: args{
		//		payload:   "05070707070a00000004a83650210a00000016019ce13845659ff2582555ec08dc322007f6493e800007070000050505050505050507070a00000016000032bb7d0084f79711f757d66b791d5290f88eb28000a80f",
		//		pubKey:    "p2pk67k5frPpxhB417bhm1n3wqH3sYKerBASTYyKXTRwkeCXBUvaaSf",
		//		signature: "p2sigk6NNw846iQ85yPuQxG9n1P2Hyumvka7zPLMxpGR6g8kT7qAWo2WrKby6uTXiRCqQbGoYnkMQAPonLeZ1CGvwWzYKUxmX7",
		//	},
		//	expResult: "",
		//	wantErr:   false,
		//},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			crPub, _ := test.args.pubKey.CryptoPublicKey()
			bt, _ := hex.DecodeString(test.args.payload)

			gotErr := verifySign(bt, test.args.signature, crPub)
			if test.wantErr != (gotErr != nil) {
				t.Errorf("wantErr: %t | err: %v", test.wantErr, gotErr)
			}
		})
	}
}
