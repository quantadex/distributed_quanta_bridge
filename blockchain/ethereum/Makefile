ganache:
	ganache-cli -m "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat" -a 8 --gasLimit 0x2FEFD800000 --defaultBalanceEther 10000 -p 7545
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

migrate-ropsten:
	truffle migrate --network ropsten

console-ropsten:
	truffle console --network ropsten

contract-info-ropsten:
	truffle exec --network=ropsten scripts/print_contract_address.js
