const QuantaCrossChain = artifacts.require("QuantaCrossChain");

var myaccounts = [
    "0xc010F4FF759605e792db672C2aF70fdCedB8D974",
    "0xc020f640E9cCc0f3fc8D4e429aa8828BEB814853",
    "0xc0305048615C1331F836298DC7e9Cc5079Cd6814"
]

module.exports = function(callback) {
    QuantaCrossChain.deployed().then(function(instance) {
        return instance.assignInitialSigners(myaccounts).then((err) => {
            console.log(err)
            callback(err)
        })
    }).catch((e) => {
        callback(e)
    })
}