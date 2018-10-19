const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    QuantaCrossChain.deployed().then(function(instance) {
        instance.txIdLast().then((txId) => {
            console.log(instance.address, Number(txId))
        })
    })
}