const SimpleToken = artifacts.require("SimpleToken");

module.exports = function(callback) {
    erc20Address = process.argv[4].trim();
    destAddress = process.argv[5].trim();
    amountInWei = process.argv[6].trim();

    return SimpleToken.at(erc20Address).transfer(destAddress, parseInt(amountInWei)).then((e) => {
        console.log(e)
    });
}