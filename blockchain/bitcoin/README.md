# Bitcoin

## For linux

https://bitcoin.org/en/full-node#ubuntu-1604

sudo apt-add-repository ppa:bitcoin/bitcoin
sudo apt-get install bitcoin-qt bitcoind

mkdir ~/.bitcoin
cp data/bitcoin.conf ~/.bitcoin


## run server
mkdir data
bitcoind

## run client

### Dump private keys

bitcoin-cli dumpwallet wallet.txt


### Generate a block

$ bitcoin-cli generate 1
[
  "19e7e275e4fee0dfb929883a17dac00b46d1b51ea4032a179583dd865e4b853b"
]

### Create multisig

bitcoin-cli addmultisigaddress 2 '["2N6VzBTQGFukxwmRcKkQTJR3rK18RNZmBtL","2N6VzBTQGFukxwmRcKkQTJR3rK18RNZmBtL"]'

### Get raw transaction

bitcoin-cli getrawtransaction 8c8aec696c84498574da7104f3ca4d8019147134f7ad0e962e2aa1d18b840080 1

### Send payment

bitcoin-cli sendtoaddress "2N6VzBTQGFukxwmRcKkQTJR3rK18RNZmBtL" 0.0001


## Useful Links

https://github.com/anders94/bitcoin-2-of-3-multisig

https://bitcoin-rpc.github.io/en/doc/0.17.99/rpc/rawtransactions/decoderawtransaction/

https://samsclass.info/141/proj/pBitc1.htm

https://en.bitcoin.it/wiki/Running_Bitcoin#Bitcoin.conf_Configuration_File