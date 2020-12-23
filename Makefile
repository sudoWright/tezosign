generate_rpc_client:
	swagger generate client -t services/rpc_client -f ./services/rpc_client/client.yml -A tezosrpc
