const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    return web3.eth.getAccounts(function(err,accounts) {
        console.log(accounts.slice(0,3))
    });
}