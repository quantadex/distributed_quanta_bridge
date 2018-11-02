const HDWalletProvider = require("truffle-hdwallet-provider");

// export MNENOMIC="empower furnace ..."
// export INFURA_API_KEY="xxx"
require('dotenv').config()  // Store environment-specific variable from '.env' to process.env

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
      gas: 14600000,
      gasPrice: 21
    },

    // testnets
    ropsten: {
      provider: function() {
        return new HDWalletProvider(process.env.MNENOMIC,
                                    "https://ropsten.infura.io/v3/" + process.env.INFURA_API_KEY)
      },
      network_id: 3,
      gas: 4700000,
      // gasPrice: 21
    }
  },

  solc: {
    optimizer: {
      enabled: true,
      runs: 500
    }
  }
};
