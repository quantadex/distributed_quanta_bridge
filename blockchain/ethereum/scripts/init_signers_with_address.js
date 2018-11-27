const QuantaCrossChain = artifacts.require("QuantaCrossChain");

var myaccounts = [
    "0x0833030c730792fDD9b77Cc54F43d7921C356Bf1",
    "0xe0006458963c3773B051E767C5C63FEe24Cd7Ff9",
    "0xba7573C0e805ef71ACB7f1c4a55E7b0af416E96A"
]

module.exports = function(callback) {
    QuantaCrossChain.deployed().then(function(instance) {
        return web3.eth.getAccounts(function(err,accounts) {
            console.log(instance, accounts)
            instance.assignInitialSigners(myaccounts)
            callback()
        });
    })
}