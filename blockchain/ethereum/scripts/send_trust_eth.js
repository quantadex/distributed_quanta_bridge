const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    QuantaCrossChain.deployed().then(function(instance) {
        return web3.eth.getAccounts(function(err,accounts) {
            return web3.eth.sendTransaction({from: accounts[0], to: instance.address, value: web3.toWei(100, 'ether')})
        })
    }).catch((ex) => {
        callback(ex)
    })
}