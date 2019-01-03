

module.exports = function(callback) {
    return web3.eth.getAccounts(function(err,accounts) {
        const t1 = web3.eth.sendTransaction({from: accounts[0], to: "0x0000000000000000000000000000000000000001", value: 1})
        const t2 = web3.eth.sendTransaction({from: accounts[0], to: "0x0000000000000000000000000000000000000002", value: 1})
        const t3 = web3.eth.sendTransaction({from: accounts[0], to: "0x0000000000000000000000000000000000000003", value: 1})
        const t4 = web3.eth.sendTransaction({from: accounts[0], to: "0x0000000000000000000000000000000000000004", value: 1})
        console.log(t1, t2, t3, t4)
        callback()
    })
}