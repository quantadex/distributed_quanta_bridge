const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    QuantaCrossChain.deployed().then(function(instance) {
        return web3.eth.getAccounts(function(err,accounts) {
            console.log(instance, accounts)
            return instance.assignInitialSigners(accounts.slice(0,3))
        });
    })
}