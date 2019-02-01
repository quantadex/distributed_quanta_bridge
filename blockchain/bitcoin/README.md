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

bitcoin-cli addmultisigaddress 2 '["0477c122f4acfbba8fb90ab353082bf782a3e26b9aaa45de66abb0464d9bf7204f5255059dd7e215849c8f29412b2a7eb95c849441dd49435d0da0a0802b5e1e9a","2NEDF3RBHQuUHQmghWzFf6b6eeEnC7KjAtR"]'

### Get raw transaction

bitcoin-cli getrawtransaction 8c8aec696c84498574da7104f3ca4d8019147134f7ad0e962e2aa1d18b840080 1

### Send payment

bitcoin-cli sendtoaddress "2N6VzBTQGFukxwmRcKkQTJR3rK18RNZmBtL" 0.0001

### Bitcoin management in Crosschain

We're going to create crosschain address with the core (3) signatures + 1 raw public key from the Graphene, which
gives us a deterministic bitcoin address given for a Quanta address.

#### Discovering Deposits

Let's scan every deposit coming in for the signatures of these 3 core addresses, decode the raw public key, then
create a deposit from  x ->  QUANTA address

#### Encode a refund

We'll scan in our unspent for any address containing the 3 core addresses, and consider them to be spendable for
the refund transaction.  We'll use as part as FIFO to assign for the spending address.

Each node needs to scan for these 3 addresses in the tx, and the amount, before signing it.


## Useful Links

https://github.com/anders94/bitcoin-2-of-3-multisig

https://bitcoin-rpc.github.io/en/doc/0.17.99/rpc/rawtransactions/decoderawtransaction/

https://samsclass.info/141/proj/pBitc1.htm

https://en.bitcoin.it/wiki/Running_Bitcoin#Bitcoin.conf_Configuration_File

https://en.bitcoin.it/wiki/List_of_address_prefixes
