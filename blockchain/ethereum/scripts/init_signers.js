const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    QuantaCrossChain.deployed().then(function(instance) {
        return web3.eth.getAccounts(function(err,accounts) {
            console.log("update signers", accounts.slice(0,3))
            return instance.assignInitialSigners(accounts.slice(0,3)).then((err) => {
                console.log(err)
                callback(err)
            })
        });
    }).catch((e) => {
        callback(e)
    })
}

