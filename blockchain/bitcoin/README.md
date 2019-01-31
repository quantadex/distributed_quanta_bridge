# Bitcoin

## run server
mkdir data
bitcoind -datadir=data/


## run client

bitcoin-cli -datadir=data/ addmultisigaddress 2 '["2N6VzBTQGFukxwmRcKkQTJR3rK18RNZmBtL","2N6VzBTQGFukxwmRcKkQTJR3rK18RNZmBtL"]'
bitcoin-cli -datadir=blockchain/bitcoin/data/ addmultisigaddress 2 '["2N6VzBTQGFukxwmRcKkQTJR3rK18RNZmBtL","2N6VzBTQGFukxwmRcKkQTJR3rK18RNZmBtL"]'


bitcoin-cli -datadir=data/ getrawtransaction 8c8aec696c84498574da7104f3ca4d8019147134f7ad0e962e2aa1d18b840080 1

```
quoc@MacBook-Pro-8:bitcoin$ {
  "txid": "8c8aec696c84498574da7104f3ca4d8019147134f7ad0e962e2aa1d18b840080",
  "hash": "8c8aec696c84498574da7104f3ca4d8019147134f7ad0e962e2aa1d18b840080",
  "version": 2,
  "size": 187,
  "vsize": 187,
  "weight": 748,
  "locktime": 101,
  "vin": [
    {
      "txid": "e29d33356214b2f35b4a4f35caed22f4cf751b0240c7c112927a5ac4d753f471",
      "vout": 0,
      "scriptSig": {
        "asm": "304402203cdb38ef16dc04fa3743583fac5eec2116945447788057620217a21f426cbc2d02206d311d9dcad4caa7472ea0612cd50c32dee789777e44af5f991e37843d05269c[ALL]",
        "hex": "47304402203cdb38ef16dc04fa3743583fac5eec2116945447788057620217a21f426cbc2d02206d311d9dcad4caa7472ea0612cd50c32dee789777e44af5f991e37843d05269c01"
      },
      "sequence": 4294967294
    }
  ],
  "vout": [
    {
      "value": 1.00000000,
      "n": 0,
      "scriptPubKey": {
        "asm": "OP_HASH160 282a74269fb3cc22d2f1313124bfe91396379242 OP_EQUAL",
        "hex": "a914282a74269fb3cc22d2f1313124bfe9139637924287",
        "reqSigs": 1,
        "type": "scripthash",
        "addresses": [
          "2Mvubsh9MS5PitwA85ix26z3X7NfDZXPnw2"
        ]
      }
    },
    {
      "value": 48.99996240,
      "n": 1,
      "scriptPubKey": {
        "asm": "OP_HASH160 5ce444a2d906abe71adacadd71d6601b3305b282 OP_EQUAL",
        "hex": "a9145ce444a2d906abe71adacadd71d6601b3305b28287",
        "reqSigs": 1,
        "type": "scripthash",
        "addresses": [
          "2N1iPcH7iQ9mSRnck2ADSE7uhWJvaALU8iM"
        ]
      }
    }
  ],
  "hex": "020000000171f453d7c45a7a9212c1c740021b75cff422edca354f4a5bf3b2146235339de2000000004847304402203cdb38ef16dc04fa3743583fac5eec2116945447788057620217a21f426cbc2d02206d311d9dcad4caa7472ea0612cd50c32dee789777e44af5f991e37843d05269c01feffffff0200e1f5050000000017a914282a74269fb3cc22d2f1313124bfe9139637924287500210240100000017a9145ce444a2d906abe71adacadd71d6601b3305b2828765000000",
  "blockhash": "36d42f00d2af94bbf7d2bb688c8fef313a92173438106d928356bad2f56802b8",
  "confirmations": 6,
  "time": 1548816038,
  "blocktime": 1548816038
}
```