#Contracts description

1. MSig contract
 
    File `contract.tz` or `contract.json`. 
 
    Contract is based on  
    `Source: https://github.com/murbard/smart-contracts/blob/master/multisig/michelson/generic.tz`
    
	Main core of signatures check was fully saved from the source. Same contract base was used in `tezos-client` generic msig contract.
	
	Our contract has been modified by adding static entrypoints for common SC actions to optimize contract call txs.
	
	To call MSig methods you should call `%main_parameter` with full path to internal action.
	All existing paths to internal actions can be found here: `https://github.com/atticlab/tezosign/blob/master/services/contract/path.go#L12`
	
	List of entrypoints:
	
    - Default. Entrypoint for account replenishments. Works without signatures check
    
    - Main. Entrypoint which requires users signatures to execute
    
	   General view of param `(pair (pair(counter (PATH TO INTERNAL ACTION AND PARAMS))) (list :sigs (option signature))`
	   
	    List of internal actions:
		
		- `:direct_action` Transfer mutez to address 
		- `:delegation` Set up delegate for msig contact 
		- `:transferFA` Call :transfer method of FA contract. `(pair address :FAContractAddress  (or ...)`
			- `:transferFA1.2` Params for call FA1.2 transfer
			- `list :transferFA2` Params for call FA2 transfer
		- `:vesting` Call of vesting contract  `(pair address :FAContractAddress  (or ...)`
			- `(option :setDelegate key_hash)` Call :setDelegate entrypoint of vesting contract
			- `(nat :vest)` Call :vest entrypoint of vesting contract
		- `(lambda unit (list operation))` Generic action with lambda function
		- `(pair (nat :threshold) (list :keys key))` Change of msig params: change payload signatures threshold and list of signers

2. Vesting contract
    
    File `vesting.tz` or `vesting.json`.
    
    Vesting contract was based on `Source: https://github.com/tqtezos/vesting-contract`
    
	Single change: default UNIT entrypoint to replenish the contract was added.
