const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    QuantaCrossChain.deployed().then(function(instance) {
        console.log("Watching event...")
        instance.TransactionResult().watch((err, result) => {
            console.log(err, result)
        })
    })
}