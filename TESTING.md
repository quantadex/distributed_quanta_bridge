## install quanta-cli

curl -L -o /usr/local/bin/quanta-cli https://github.com/quantadex/quanta-cli/releases/download/0.5.1/quanta-cli.linux.amd64
chmod +x /usr/local/bin/quanta-cli

https://github.com/quantadex/quanta_book/wiki/Accessing-testnet

## Set up the keys
quanta-cli set config:network test

quanta-cli account new xc_signer
QBU734ICFLPITA6IYTZCUJEQWLIOV2R5DHVCAXGJ4AB4XMBTMM4MSB7P ZBW6WP3ZRMUFGC4GDTWJOXJTCYUJPN37GJXYMFPZZFRCLBAGXJ3YRWQZ

quanta-cli account new xc_node1
QCU3VV7R7NCTXH6Z6XXJAHGRZYWQSBXN2TG4FXGHHVYOLX5MUV4JGQ7B ZDFTL63XHTY7XMLQHKW22TGVGBHNKTZQTWK7N3HZCPLDCTNN63T2BJFB

quanta-cli account new xc_node2
QAZV6UM4QWVJ236RCAVIGAD3CMVD3O5CN2VOGI4PFIGONTR6QWSU3SMH ZBVZ3QWC46GRO4YT5BXWPJ6UAFGXB4FCLP2G5AC64AC5C4KZMQ6H4PT5

quanta-cli account new xc_node3
QALTHTQZ3F6XC7QVQDDAUQ44VOTXKMFNT3IMG46MW6II4LMCGGUL6NAY ZDUTFAMBO3VTR7ADO33APQSNFDIDBJSCDB5SX7LFKG55JFX5SSWZDKSF

quanta-cli friendbot xc_signer
quanta-cli friendbot xc_node1
quanta-cli friendbot xc_node2
quanta-cli friendbot xc_node3


quanta-cli  signer add QCU3VV7R7NCTXH6Z6XXJAHGRZYWQSBXN2TG4FXGHHVYOLX5MUV4JGQ7B 1 --to xc_signer

quanta-cli  signer add QAZV6UM4QWVJ236RCAVIGAD3CMVD3O5CN2VOGI4PFIGONTR6QWSU3SMH 1 --to xc_signer

quanta-cli  signer add QALTHTQZ3F6XC7QVQDDAUQ44VOTXKMFNT3IMG46MW6II4LMCGGUL6NAY 1 --to xc_signer

## set threshold medium/high to 3 signers
quanta-cli signer thresholds xc_signer 1 3 3

## Setup infura and get a key for ropsten

http://infura.io

## install metamask

Setup a memomic pass phrase:
https://support.dex.top/hc/en-us/articles/360004125614-How-to-Create-Mnemonic-Phrase-with-MetaMask-

Save your phrases
and

Note: Memonic key is the root key that generate other keys


## Deploy in ethereum in ropsten

https://faucet.metamask.io/
Click request 1 Ether from faucet (with metamask loggedin)

## Deploy the trust contract

In ethereum_distributed

> cd blockchain/ethereum
> export MNENOMIC=xxx
> export INFURA_API_KEY=xxx
> truffle compile
> truffle migrate --network ropsten

Note the the Trust contract.

## Deploy in privatenet

### deploy geth
geth --datadir ./geth/data init ./geth/genesis.json
make start-geth


> cd blockchain/ethereum
> export MNENOMIC=xxx
> truffle compile
> truffle migrate --network ropsten --reset


## setup nodes

issuer is the xc_signer public address
node 1 nodekey is the private address of node1

Trust Addr: 0x3e758fcaea2788bc1d7f6cd2366d64b872db67c6

## get those top memonic
truffle exec scripts/print_accounts.js  --network ropsten


## create database

quoc@MacBook-Pro-8:~/Projects/go/src/github.com/quantadex/distributed_quanta_bridge/node$ psql
psql (9.6.3)
Type "help" for help.
quoc=# create database crosschain_1;
ACREATE DATABASE
quoc=# create database crosschain_2;
CREATE DATABASE
quoc=# create database crosschain_3;
CREATE DATABASE
quoc=# \q

## setup the signers for trust contract

Change the following:
Create 3 configurations, node1.yml, node2.yml, node3.yml


* change port from: 5000, 5001, 5002
* Issuer address, the public key xc_signer
* NodeKey: Private key for each of the node (node1, node2, node3)
* EthereumKeyStore: private key for ethereum (node1, node2, node3)
* EthereumTrustAddr: trust address from earlier
* Ropsten Url
* EthereumBlockStart: trust

```
ListenIp: 0.0.0.0
ListenPort: 5000
UsePrevKeys: true
KvDbName: ./data/node1
CoinName: ETH
IssuerAddress: QBU734ICFLPITA6IYTZCUJEQWLIOV2R5DHVCAXGJ4AB4XMBTMM4MSB7P
NodeKey: ZDFTL63XHTY7XMLQHKW22TGVGBHNKTZQTWK7N3HZCPLDCTNN63T2BJFB
HorizonUrl: http://testnet-02.quantachain.io:8000/
NetworkPassphrase: QUANTA Test Network ; September 2018
RegistrarIp: 0.0.0.0
RegistrarPort: 5100
EthereumNetworkId: 3
EthereumBlockStart: 4300900
EthereumRpc: https://ropsten.infura.io/v3/7b880b2fb55c454985d1c1540f47cbf6
EthereumKeyStore: 5E5C9D89B0393E673FDFB21DB2EB7D3368AB1D4D7F86293E7C8DD8822218A753
EthereumTrustAddr: 0x3e758fcaea2788bc1d7f6cd2366d64b872db67c6
HEALTH_INTERVAL: 5
DatabaseUrl: postgres://postgres:@localhost/crosschain_1
```
## Running the node

> quanta-cli account new xc_demo
QC2PFVZGWAZ2MXZPIHQLO2LJ5F5PZCVRN756FEFVHLGMWANQK3SOTKWW ZBH53BENZ5B3SFQ7YWIY4GMH3E73LAX3HYOANCOC6HAM7TICC24ZQPD2

truffle exec scripts/new_forward_address.js  <trust address> <QUANTA_account>
> truffle exec scripts/new_forward_address.js 0x3e758fcaea2788bc1d7f6cd2366d64b872db67c6 QC2PFVZGWAZ2MXZPIHQLO2LJ5F5PZCVRN756FEFVHLGMWANQK3SOTKWW  --network ropsten

Using network 'ropsten'.
Creating for  0x3e758fcaea2788bc1d7f6cd2366d64b872db67c6 QC2PFVZGWAZ2MXZPIHQLO2LJ5F5PZCVRN756FEFVHLGMWANQK3SOTKWW
0x1485427f94275b2f6d23dcaa5d07cb2059f1b481

This is my forwarding address: 0x1485427f94275b2f6d23dcaa5d07cb2059f1b481

Run the nodes:

```
./node -config privatenet/node1.yml -registry
postgres://postgres:@localhost/crosschain_1
2018/11/27 16:16:26 I [5000] Initialize ledger
2018/11/27 16:16:26 I [5000] Submitworker started
2018/11/27 16:16:26 I [5000] REST API started at :0...

./node -config privatenet/node2.yml
postgres://postgres:@localhost/crosschain_2
2018/11/27 16:16:26 I [5000] Initialize ledger
2018/11/27 16:16:26 I [5000] Submitworker started
2018/11/27 16:16:26 I [5000] REST API started at :0...

./node -config privatenet/node3.yml
postgres://postgres:@localhost/crosschain_3
2018/11/27 16:16:26 I [5000] Initialize ledger
2018/11/27 16:16:26 I [5000] Submitworker started
2018/11/27 16:16:26 I [5000] REST API started at :0...
```

Tips to reset:

Reboot from scratch:
- control + c out of each node
- delete database, and re-create again (not needed very often)
- remove the kv database, data/node*
- start registry, node2, node3

Normal reboot:
- control + c out of each node
- start registry, node2, node3

Basic test:
- Minimum block: Trust creation
- Create forward address (block)  - forwarding once
- Send  0.10 ETH to  0x1485427f94275b2f6d23dcaa5d07cb2059f1b481 with metamask (Note: block it is)
- Look at the node's logs
- you should see issued assets by QBU734ICFLPITA6IYTZCUJEQWLIOV2R5DHVCAXGJ4AB4XMBTMM4MSB7P


Withdraw:

setup asset:  (replace issuer)
> quanta-cli  asset set KETH QBU734ICFLPITA6IYTZCUJEQWLIOV2R5DHVCAXGJ4AB4XMBTMM4MSB7P --code ETH

See: https://quanta.gitbook.io/documentation/testnet/transfer-eth-erc-20



