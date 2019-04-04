# Crosschain Architecture

## Setup

```
cd cli/ethereum && go build
./ethereum -config ../../node/ropsten/node1.yml  # run for each node

cd node && go build
./node -config ropsten/node1.yml # run for each node
```

```
docker run --name bitcoind -v "$PWD/bitcoin-data:/data" nicolasdorier/docker-bitcoin:0.17.0 bitcoind -testnet -deprecatedrpc=signrawtransaction -txindex -deprecatedrpc=accounts
```

## Deploy Testnet

```
ssh -i ~/.ssh/testnet-oregon.pem  ec2-user@ec2-54-188-223-216.us-west-2.compute.amazonaws.com
ssh -i ~/.ssh/testnet-oregon.pem  ec2-user@ec2-34-221-59-194.us-west-2.compute.amazonaws.com
ssh -i ~/.ssh/testnet-oregon.pem  ec2-user@ec2-34-219-198-107.us-west-2.compute.amazonaws.com

ssh ec2-user@ec2-54-188-223-216.us-west-2.compute.amazonaws.com
ssh ec2-user@ec2-34-221-59-194.us-west-2.compute.amazonaws.com
ssh ec2-user@ec2-34-219-198-107.us-west-2.compute.amazonaws.com
```

```
IP Addresses  Public         Internal
Crosschain1: 54.188.223.216 192.168.137.186
Crosschain2: 34.221.59.194 192.168.174.110
Crosschain3: 34.219.198.107 192.168.171.58
```
## Stopping node

```
docker-compose  stop crosschain1

## Pull latest
$(aws ecr get-login --no-include-email --region us-east-1)
docker pull 691216021071.dkr.ecr.us-east-1.amazonaws.com/quanta-bridge:latest

## Restart all dockers
docker-compose  up --force-recreate --build -d crosschain_eth
docker-compose  up --force-recreate --build -d crosschain_btc
docker-compose  up --force-recreate --build -d crosschain1
```

## Get logs

docker-compose  logs -f crosschain1
docker-compose  logs -f crosschain_eth
docker-compose  logs -f crosschain_btc


ETH Forwarding Contract Deployed:
https://ropsten.etherscan.io/tx/0x2255866778b5fc0e1d416b93d71d9a42c2fb66b9493cf20a40085f3aacaf8dc2

Contract: [0x21b2d8f88f5d60f28739ebeefc09220687ff00f1](https://ropsten.etherscan.io/address/0x21b2d8f88f5d60f28739ebeefc09220687ff00f1)
Forwards: [0xBD770336fF47A3B61D4f54cc0Fb541Ea7baAE92d](https://ropsten.etherscan.io/address/0xBD770336fF47A3B61D4f54cc0Fb541Ea7baAE92d)
QUANTA Addr: [QDIX3EOMEWN7OLZ3BEIN5DE7MCVSAP6547FFM3FFITQSTFXWUK4XA2NB](http://testnet-02.quantachain.io:8000/accounts/QDIX3EOMEWN7OLZ3BEIN5DE7MCVSAP6547FFM3FFITQSTFXWUK4XA2NB) (test2)


**Test deposits:**

0.1234 ETH to 0x21b2d8f88f5d60f28739ebeefc09220687ff00f1
TX: https://ropsten.etherscan.io/tx/0x8784265b4f2aa6c6460cc8bf3be46fb4d269a4d00c099bd205a4ae1ad1b96138


**Test Withdrawal**

lumen pay 0.10 TETH --from test2 --to issuer --memohash MHhiYTc1NzNDMGU4MDVlZjcxQUNCN2YxYzRhNTVFN2IwYWY0MTZFOTZB

Sends back to : 0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A

## Open questions

It's possible that eth timestamp is a head of QUANTA timestamp, vice versa. So we should throttle based on universal time.

# Truffle Contract Interaction Tips

* Start the truffle console

    $ cd blockchain/ethereum
    $ npm install truffle-hdwallet-provider
    $ npm install dotenv
    $ export MNENOMIC='empower furnace...'
    $ export INFURA_API_KEY='0e17d......'
    $ truffle console --network ropsten

* get the contract by address  

    truffle(ropsten)> const contract = QuantaCrossChain.at("0xbd770336ff47a3b61d4f54cc0fb541ea7baae92d")
    truffle(ropsten)> const util = require('ethereumjs-util')

* convert a private ethereum key to the public ethereum address    

    truffle(ropsten)> util.bufferToHex(util.publicToAddress(util.privateToPublic("0xc87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3")))
    '0x627306090abab3a6e1400e9345bc60c78a8bef57'

* check to see if the public ethereum address is a signer for the contract

    truffle(ropsten)> contract.isSigner("0x0xba7573c0e805ef71acb7f1c4a55e7b0af416e96a")
    true

* get the current ethereum block number    

    web3.eth.getBlockNumber((err, res) => {console.log(res)})
    4345932

## Links
