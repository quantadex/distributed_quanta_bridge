const QuantaForwarder = artifacts.require("QuantaForwarder");

module.exports = function(callback) {
    erc20 = process.argv[4].trim()
    forw = process.argv[5].trim()

    return QuantaForwarder.at(forw).flushTokens(erc20).then((e) => {
        console.log(e)
    });
}