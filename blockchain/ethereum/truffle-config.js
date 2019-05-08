const HDWalletProvider = require("truffle-hdwallet-provider");

// export MNENOMIC="empower furnace ..."
// export INFURA_API_KEY="xxx"
require('dotenv').config()  // Store environment-specific variable from '.env' to process.env
const PrivateKeyProvider = require("truffle-privatekey-provider");

module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!

  // network_id: identifier for network based on ethereum blockchain. Find out more at https://github.com/ethereumbook/ethereumbook/issues/110
  // gas: gas limit
  // gasPrice: gas price in gwei

  networks: {
      development: {
          host: "127.0.0.1",
          port: 7545,
          network_id: "*", // Match any network id
          gas: 6721975,
          gasPrice: 21
      },
      test: {
          host: "127.0.0.1",
          port: 7545,
          network_id: "*", // Match any network id
          gas: 3000000,
          gasPrice: 30,
      },
      // testnets
      ropsten: {
          provider: function () {
              return new HDWalletProvider(process.env.MNENOMIC,
                  "https://ropsten.infura.io/v3/" + process.env.INFURA_API_KEY)
          },
          network_id: 3,
          gas: 4700000,
          // 50 wei is way too small, comment this out and let the system suggest the gas price
          // may be a little expensive, but transaction will go through.

          //gasPrice: 50
      },
      mainnet: {
          provider: function () {
              return new PrivateKeyProvider(process.env.PRIVATE_KEY,
                  "https://mainnet.infura.io/v3/" + process.env.INFURA_API_KEY)
          },
          network_id: 1,
          gasPrice: 4000000000,
          gas: 1000000
      },
  },
  solc: {
    optimizer: {
      enabled: true,
      runs: 500
    }
  }
};
