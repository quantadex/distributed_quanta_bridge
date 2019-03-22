module.exports = function(callback) {

    const acc = "0x3B391502ef005eea0aCb672CdCad10ADaEB51e1A"
    const mainAcc = "0xbA420EF5D725361d8fDc58Cb1e4fa62EDa9EC990"
    const privateKey = ""
    const gas = web3.toWei("1", "gwei")

    web3.eth.getBalance(acc, (err, balance) => {
        console.log("balance", balance)
        const amount = balance - gas
        var raw = {
            "from": acc,
            "to": mainAcc,
            "value": web3.toHex(amount),
            "gas": web3.toHex(gas),
            "chainId": 1
        };

        web3.eth.accounts.signTransaction(raw, (e,x)=>{}).then(signed=>{
            console.log(raw, signed)

            callback()
        })

    });
}