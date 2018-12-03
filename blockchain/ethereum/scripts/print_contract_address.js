/**
 * Usage:
 *     truffle exec scripts/print_contract_address.js
 */

const QuantaCrossChain = artifacts.require("QuantaCrossChain");

module.exports = function(callback) {
    QuantaCrossChain.deployed().then( async function(instance) {
        var txId = await instance.txIdLast()
        var signers = await instance.numSigners()
        // var signersAddr = await instance.signers();

        var balance = await web3.eth.getBalance(instance.address)
        console.log(instance.address, Number(txId), balance.toNumber());
        console.log("Signers", signers.toNumber())
        //console.log(instance)
        callback();
    })
}
