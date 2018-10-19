pushd ../../../blockchain/ethereum
truffle compile
truffle-export-abi
popd
$GOPATH/bin/abigen --pkg contracts --type TrustContract --out trust_contract.go --abi ../../../blockchain/ethereum/build/ABI.json --solc solcjs
