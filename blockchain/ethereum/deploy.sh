rm -rf build
truffle compile
truffle migrate --network test --reset
truffle exec scripts/init_signers.js --network test
truffle exec scripts/send_trust_eth.js --network test
truffle exec scripts/print_contract_address.js --network test