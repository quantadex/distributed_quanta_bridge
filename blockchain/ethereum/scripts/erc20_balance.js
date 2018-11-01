const SimpleToken = artifacts.require("SimpleToken");

module.exports = function(callback) {
    erc20Address = process.argv[4].trim();
    destAddress = process.argv[5].trim();

    return SimpleToken.at(erc20Address).balanceOf(destAddress).then((e) => {
        console.log(e.toNumber())
    });
}