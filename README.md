# Crosschain Architecture


## Configuring Lumen

lumen set config:network "custom;http://testnet-02.quantachain.io:8000;QUANTA Test Network ; September 2018"

## Testnet

ISSUER:
[QCISRUJ73RQBHB3C4LA6X537LPGSFZF3YUZ6MOPUOUJR5A63I5TLJML4](http://testnet-02.quantachain.io:8000/accounts/QCISRUJ73RQBHB3C4LA6X537LPGSFZF3YUZ6MOPUOUJR5A63I5TLJML4)

SIGNERS:
1. QAIS24MZVMNZFOGDFDXFSPIZTSNS46H7Y2JTUYXNYMDBK6ZEIBEX5JDN
2. QBIS326UFEHDA36IZTLBSKBF245DJ37JBMF3FEC45AWVIRM36KDB2LFQ
3. QCN2DWLVXNAZW6ALR6KXJWGQB4J2J5TBJVPYLQMIU2TDCXIOBID5WRU5


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
