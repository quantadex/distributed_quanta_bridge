const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    QuantaCrossChain.deployed().then( async function(instance) {
        var txId = await instance.txIdLast()
        var balance = await web3.eth.getBalance(instance.address)
        console.log(instance.address, Number(txId), balance.toNumber())
    })
}