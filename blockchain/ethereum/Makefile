ganache:
	ganache-cli -m "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat" -a 8 --gasLimit 0x2FEFD800000 --defaultBalanceEther 10000 -p 8545
migrate:
	truffle migrate
console:
	truffle console
debugger:
	truffle develop --log
compile:
	truffle compile
contract-info:
	truffle exec scripts/print_contract_address.js

.PHONY:
test: .PHONY
	truffle test
abi:
	truffle-export-abi

mine-geth:
	geth attach geth/data/geth.ipc --exec "miner.start()"
start-geth:
	geth --datadir ./geth/data --networkid 15 --ipcpath ./geth/data/geth.ipc --rpc  --unlock 0  --mine --password ./geth/emptypassword.txt  --rpcapi "db,eth,net,web3,miner,debug" --nodiscover  --port "30304" --rpcaddr "0.0.0.0" --rpccorsdomain "*" --rpcport 7545 --verbosity 3 --miner.gasprice 2000000000 --minerthreads=1 --etherbase=0x0000000000000000000000000000000000000000 --debug --vmdebug

migrate-ropsten:
	truffle migrate --network ropsten

console-ropsten:
	truffle console --network ropsten

contract-info-ropsten:
	truffle exec --network=ropsten scripts/print_contract_address.js
