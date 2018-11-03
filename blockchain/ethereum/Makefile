ganache:
	ganache-cli -m "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat" -a 8 --gasLimit 0x2FEFD800000 --defaultBalanceEther 10000 -p 7545
console:
	truffle console
debugger:
	truffle develop --log
compile:
	truffle compile

.PHONY:
test: .PHONY
	truffle test
migrate:
	truffle migrate
abi:
	truffle-export-abi

mine-geth:
	geth attach geth/data/geth.ipc --exec "miner.start()"
start-geth:
	geth --datadir ./geth/data --networkid 15 --ipcpath ./geth/data/geth.ipc --rpc  --unlock 0  --mine --password ./geth/emptypassword.txt  --rpcapi "db,eth,net,web3,miner" --nodiscover  --port "30304" --rpcaddr "0.0.0.0" --rpccorsdomain "*"
