#!/bin/bash

cd ../../blockchain/ethereum/
geth --datadir ./geth/data/ init ./geth/genesis.json
make start-geth