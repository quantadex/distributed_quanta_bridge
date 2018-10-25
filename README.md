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


## Open questions

It's possible that eth timestamp is a head of QUANTA timestamp, vice versa. So we should throttle based on universal time.

## Links