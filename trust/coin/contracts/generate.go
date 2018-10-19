//go:generate sh generate_abi.sh
//not:generate /usr/local/Cellar/ethereum/1.8.14/bin/abigen --pkg contracts --type TrustContract --out trust_contract.go --abi ../../../blockchain/ethereum/build/ABI.json

package contracts

