const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    QuantaCrossChain.deployed().then(function(instance) {
        console.log("Send payment")
        return instance.paymentTx(
            1, 0, "0xffff", 150,[],[],[]).then((err) => {

            console.log(err)
            callback()
        });
    })
}