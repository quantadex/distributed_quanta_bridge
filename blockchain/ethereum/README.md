# Cross chain Ethereum Smart Contract

# Prerequisites

## System

* `node >= 8.11.3`
* `npm >= 5.6`

## Ethereum Development Tools

    $ npm install -g truffle ganache-cli

# Check versions

    $ npm --version
    # 6.4.1

    $ node --version
    # v8.12.0

    $ truffle version
    Truffle v4.1.14 (core: 4.1.14)
    Solidity v0.4.24 (solc-js)

    $ ganache-cli --version
    Ganache CLI v6.1.8 (ganache-core: 2.2.1)

# Start a local development Ethereum node

Run this command to start your own private node:

    $ ganache-cli -p 7545

    # or set some options

    # need at least 8 accounts for the tests
    $ ganache-cli -m "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat" --gasLimit 0x2FEFD800000 -a 8 --defaultBalanceEther 10 -p 7545
    # or
    $ make ganache

# Compile the contracts

## Using the truffle console

    # or `make console`
    $ truffle console

    truffle(development)> test

    truffle(development)> compile
    truffle(development)> migrate

## Using single truffle commands

Compile the contracts.

    # or `make compile`
    $ truffle compile

# Testing

Run the following command to test the contracts.

    # or `make test`
    $ truffle test

# Deploy to the development network

    # or `make migrate`
    $ truffle migrate --network development

    Using network 'development'.

# Debugger

In three separate terminals:

    make ganache
    make debugger
    make console

If you see an exception in the ganache logs, take a note of the transaction id.

For example:

    Transaction: 0xae4b8f2e781cb7260853fb3a7a807f3b71e8b6ae4510c0d72ba203fd4ad70241
    Gas usage: 14600000
    Block Number: 46
    Block Time: Wed Oct 17 2018 12:22:18 GMT-0700 (PDT)
    Runtime Error: invalid opcode

In the console, debug the transaction id like so:

    truffle(development)> debug 0xae4b8f2e781cb7260853fb3a7a807f3b71e8b6ae4510c0d72ba203fd4ad70241
    ...
    Commands:
    (enter) last command entered (step next)
    (o) step over, (i) step into, (u) step out, (n) step next
    (;) step instruction, (p) print instruction, (h) print this help, (q) quit
    (b) toggle breakpoint, (c) continue until breakpoint
    (+) add watch expression (`+:<expr>`), (-) remove watch expression (-:<expr>)
    (?) list existing watch expressions
    (v) print variables and values, (:) evaluate expression - see `v`

Then use the commands to step through the code.
It would also help to add revert, assert and require statements in the code to test for conditions you were expecting.
    
# Troubleshooting

* Problem: Could not connect to your Ethereum client.

      Could not connect to your Ethereum client. Please check that your Ethereum client:
          - is running
          - is accepting RPC connections (i.e., "--rpc" option is used in geth)
          - is accessible over the network
          - is properly configured in your Truffle configuration file (truffle.js)

  Solution: Make sure `ganache-cli` is running.

# References

* [Truffle Framework Configuration](http://truffleframework.com/docs/advanced/configuration>)
