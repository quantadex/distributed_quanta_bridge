# Quanta Chain Contract Deployment

## Ropsten TESTNET Deployment 2018-10-24 14:17

Deployed the first Quanta chain contract to Ropsten!

Deployment Time: about 2 minutes
Deployment Cost: 0.3573677 ETH

Contract Address: [0xfe5e4bc6efa05cef33c7eb4adaa933cfa1288417d355509b6f37bee108fd145c](https://ropsten.etherscan.io/tx/0xfe5e4bc6efa05cef33c7eb4adaa933cfa1288417d355509b6f37bee108fd145c)

### Notes

* expensive first time deployment due to needing to deploythe library contract files as well
* the contract itself, however, still costs 0.23 ETH

### Console Log

    $ git rev-parse --short HEAD
    9126249
    $ export MNENOMIC='empower furnace...'
    $ export INFURA_API_KEY='0e17d4...'
    $ truffle migrate --network ropsten

    Using network 'ropsten'.

    Running migration: 1_initial_migration.js
      Deploying Migrations...
      ... 0x71c64215af80c35115f912218a7dda627b6f8540e6e7d488240600258a3a90f2
      Migrations: 0x139ac76fdeeb43c583da0ca9811301b80d06f051
    Saving successful migration to network...
      ... 0x01ca2e8c3ea8c27abfa471ff1cdccdb96e17510dcff384398472850ea3eade4e
    Saving artifacts...
    Running migration: 2_QuantaCrossChain.js
      Deploying libbytes...
      ... 0x63c6058c11e48f7ef5dc0f9ba354d847b8cf88b3424823e7d6035ab78c0f4c12
      libbytes: 0x2e80ca25f890367075a9b4f341dd5f15d98a1a70
      Deploying ECTools...
      ... 0x7f5a35098c3ef698b86fafa03ce59f906f7c7798ef2d2c2883883296a5818b27
      ECTools: 0xc6a7e958ff2c0d6ae0ef518101c5e523a6eba45f
      Linking ECTools to QuantaCrossChain
      Deploying QuantaCrossChain...
      ... 0xfe5e4bc6efa05cef33c7eb4adaa933cfa1288417d355509b6f37bee108fd145c
      QuantaCrossChain: 0xbd770336ff47a3b61d4f54cc0fb541ea7baae92d
    Saving successful migration to network...
      ... 0x83e7c5727c651542627eb51940ea18a6d31185b9be9109e5cc282a87fe8da5c4
    Saving artifacts...

### Transaction Details

https://ropsten.etherscan.io/address/0xe2ac9076b2e846864ef038d01ef2322771aa6df4

    Txhash	Blockno	UnixTimestamp	DateTime	From	To	ContractAddress	Value_IN(ETH)	Value_OUT(ETH)	CurrentValue @ $0/Eth	TxnFee(ETH)	TxnFee(USD)	Historical $Price/Eth	Status	ErrCode
    0x71c64215af80c35115f912218a7dda627b6f8540e6e7d488240600258a3a90f2	4294601	1540415824	10/24/2018 9:17:04 PM	0xe2ac9076b2e846864ef038d01ef2322771aa6df4		0x139ac76fdeeb43c583da0ca9811301b80d06f051	0	0	0	0.0224195	0			
    0x01ca2e8c3ea8c27abfa471ff1cdccdb96e17510dcff384398472850ea3eade4e	4294602	1540415841	10/24/2018 9:17:21 PM	0xe2ac9076b2e846864ef038d01ef2322771aa6df4	0x139ac76fdeeb43c583da0ca9811301b80d06f051		0	0	0	0.0041908	0			
    0x63c6058c11e48f7ef5dc0f9ba354d847b8cf88b3424823e7d6035ab78c0f4c12	4294604	1540415881	10/24/2018 9:18:01 PM	0xe2ac9076b2e846864ef038d01ef2322771aa6df4		0x2e80ca25f890367075a9b4f341dd5f15d98a1a70	0	0	0	0.0074748	0			
    0x7f5a35098c3ef698b86fafa03ce59f906f7c7798ef2d2c2883883296a5818b27	4294606	1540415965	10/24/2018 9:19:25 PM	0xe2ac9076b2e846864ef038d01ef2322771aa6df4		0xc6a7e958ff2c0d6ae0ef518101c5e523a6eba45f	0	0	0	0.0917492	0			
    0xfe5e4bc6efa05cef33c7eb4adaa933cfa1288417d355509b6f37bee108fd145c	4294608	1540415972	10/24/2018 9:19:32 PM	0xe2ac9076b2e846864ef038d01ef2322771aa6df4		0xbd770336ff47a3b61d4f54cc0fb541ea7baae92d	0	0	0	0.2288426	0			
    0x83e7c5727c651542627eb51940ea18a6d31185b9be9109e5cc282a87fe8da5c4	4294610	1540415994	10/24/2018 9:19:54 PM	0xe2ac9076b2e846864ef038d01ef2322771aa6df4	0x139ac76fdeeb43c583da0ca9811301b80d06f051		0	0	0	0.0026908	0			



