## install quanta-cli

curl -L -o /usr/local/bin/quanta-cli https://github.com/quantadex/quanta-cli/releases/download/0.5.1/quanta-cli.linux.amd64
chmod +x /usr/local/bin/quanta-cli

npm i -g ganache-cli@6.2.3

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
> cd blockchain/ethereum
> geth --datadir ./geth/data init ./geth/genesis.json

In one terminal:
> make start-geth
> truffle compile
> truffle migrate --network test
> truffle exec scripts/print_accounts.js  --network test
> truffle exec scripts/init_signers.js  --network test

Note: make sure your EthereumBlockStart config is 0, setup your configurations
​      in a privatenet

Your private ethereum:

(0) c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3
(1) ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f
(2) 0dbbe8e4ae425a6d2687f1a7e3ba17bc98c673636790f1b8ad91193c05875ef1

Use the private key, and setup your (4) key in your metamask:

0xc88b703fb08cbea894b6aeff5a544fb92e78a18e19814cd85da83b71f772aa6c

Create forwarding:
> truffle exec scripts/new_forward_address.js 0x9fbda871d559710256a2502a2517b794b482db40 QC2PFVZGWAZ2MXZPIHQLO2LJ5F5PZCVRN756FEFVHLGMWANQK3SOTKWW --network test
forwarding:  0xfb88de099e13c3ed21f80a7a1e49f8caecf10df6

## setup nodes

issuer is the xc_signer public address
node 1 nodekey is the private address of node1

Trust Addr: 0x3e758fcaea2788bc1d7f6cd2366d64b872db67c6

## get those top memonic
truffle exec scripts/print_accounts.js  --network ropsten


## create database
```
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
```
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


# Testing tips:

Always know what you should expect to see.



**Running the nodes:**

1. We should expect the nodes to process deposits for each block.  

```
Got blocks [1976 1977 1978 1979 1980 1981 1982 1983 1984 1985 1986 1987 1988 1989 1990 1991 1992 1993 1994]
***** Start # of blocks=0 man.N=3,man.Q=3 *** 
http://testnet-02.quantachain.io:8000//accounts/QBU734ICFLPITA6IYTZCUJEQWLIOV2R5DHVCAXGJ4AB4XMBTMM4MSB7P/payments?order=asc&limit=100&cursor=0
2018/11/28 11:54:19 I [5000] QuantaToCoin refunds []
2018/11/28 11:54:19 I [5000] Next cursor is = 0, numRefunds=0
2018/11/28 11:54:22 D [5000] Coin2Quanta: No new block last=1994 top=1994
```



2. Creating a new forwarding address

```
   truffle exec scripts/new_forward_address.js 0x9fbda871d559710256a2502a2517b794b482db40 QC2PFVZGWAZ2MXZPIHQLO2LJ5F5PZCVRN756FEFVHLGMWANQK3SOTKWW --network test
   
   
   2018/11/28 11:36:14 I [5000] New Forwarder Address ETH->QUANTA address, 0xAa588d3737B611baFD7bD713445b314BD453a5C8 -> QC2PFVZGWAZ2MXZPIHQLO2LJ5F5PZCVRN756FEFVHLGMWANQK3SOTKWW
   
```

Take note of the block # when it was executed.  so you know which block number to look for in the logs


3. Forwarding address should show up on each of the node:

```
   2018/11/28 11:36:14 I [5000] New Forwarder Address ETH->QUANTA address, 0xAa588d3737B611baFD7bD713445b314BD453a5C8 -> QC2PFVZGWAZ2MXZPIHQLO2LJ5F5PZCVRN756FEFVHLGMWANQK3SOTKWW
```



4. Make a deposit, take a note of amount, and block number



5. Withdrawal

Send Refund:

http://tomeko.net/online_tools/hex_to_base64.php
0x0d1d4e623d10f9fba5db95830f7d3839406c6af2 -> DR1OYj0Q+ful25WDD304OUBsavI=
quanta-cli pay 0.15 KETH --from xc_demo --to xc_signer --memotext "DR1OYj0Q+ful25WDD304OUBsavI="


Check trust:

```
geth attach ./geth/data/geth.ipc
web3.fromWei(eth.getBalance('0xa4392264a2d8c998901d10c154c91725b1bf0158'),'ether')

eth.sendTransaction({from:eth.accounts[0], to:"0xa4392264a2d8c998901d10c154c91725b1bf0158", value: web3.toWei(10, "ether")})

```

Check refund:

```
> eth.accounts[4]
"0x0d1d4e623d10f9fba5db95830f7d3839406c6af2"
web3.fromWei(eth.getBalance(eth.accounts[4]),'ether')
```


Leader node:

```
2018/11/30 23:48:59 I [5000] Start new round bda75e3a4c3cfd5822c512dd6f0e9408c680e8b2b5188b98e62a61be458029c2 ETH to=0x0d1d4e623D10F9FBA5Db95830F7d3839406C6AF2 amount=15000000000000000
2018/11/30 23:48:59 I [cosi] Start new round leader=true nodes=3 threshold=3
2018/11/30 23:48:59 I [cosi] Got commitment 2/3
2018/11/30 23:48:59 I [cosi] Got commitment 3/3
2018/11/30 23:48:59 I [cosi] Got total of 3 commitments, moving forward
2018/11/30 23:48:59 I [5000] Sign msg 6f53681a1d6e2c1a9d3200efe771f4e16c2974c46adf2475c57766539500134d425e001d9c244471ba0a91699c57fac7f2a5c69fa62d72c0be8276da3a33219e00
2018/11/30 23:48:59 I [cosi] Got signature 2/3
2018/11/30 23:48:59 I [cosi] Got signature 3/3
2018/11/30 23:48:59 I [5000] Great! Cosi successfully signed refund
Sending from 0x627306090abaB3A6e1400e9345bC60c78a8BEf57
Submit to contract=0x9FBDa871d559710256a2502A2517b794B482Db40 erc20=0x0000000000000000000000000000000000000000 to=0x0d1d4e623D10F9FBA5Db95830F7d3839406C6AF2 amount=15000000000000000
signatures (3) [6f53681a1d6e2c1a9d3200efe771f4e16c2974c46adf2475c57766539500134d425e001d9c244471ba0a91699c57fac7f2a5c69fa62d72c0be8276da3a33219e00 b27a23339259ff2df83e3f33f69094778bd6fe508602aceb926e069414fdc4397fbabda0e33fde87f3e32ad00d52b95db01800d46e9b28f9b4c5fc1551b2880400 7c451aa8541380161b07e7210dd9dc5ec221fac6b840da8dbcda7edc72adcb214ed91d4728e05e7149e7528e5f5eb187fae98aca1a13e326d82108e2ef170acc01]
prepare to send to contract
2018/11/30 23:48:59 I [5000] Submitted withdrawal in tx=0x7e22ddcdd66b23247de1c0c4c7c6db3bc6ca7da7f31c7f2a9e41e5a3f74fb1ee
```


### WHY PRIVATE NET KEEPS FAILING?

We need to bootstrap the precompiled contracts on the testnet by
sending 1 wei to each contract in 1,2,3,4

https://ethereum.stackexchange.com/questions/15479/list-of-pre-compiled-contracts
https://ethereum.stackexchange.com/questions/440/whats-a-precompiled-contract-and-how-are-they-different-from-native-opcodes
https://medium.com/@rbkhmrcr/precompiles-solidity-e5d29bd428c4

```

eth.sendTransaction({from:eth.accounts[0], to:"0x0000000000000000000000000000000000000001", value: 1});
eth.sendTransaction({from:eth.accounts[0], to:"0x0000000000000000000000000000000000000002", value: 1});
eth.sendTransaction({from:eth.accounts[0], to:"0x0000000000000000000000000000000000000003", value: 1});
eth.sendTransaction({from:eth.accounts[0], to:"0x0000000000000000000000000000000000000004", value: 1});


```