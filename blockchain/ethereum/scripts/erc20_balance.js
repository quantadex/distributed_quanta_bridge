const ERC20Basic = artifacts.require("zepplelin/token/ERC20Basic");

module.exports = function(callback) {
    erc20Address = process.argv[4].trim();
    destAddress = process.argv[5].trim();
    console.log("checking. ", erc20Address, destAddress);

    return ERC20Basic.at(erc20Address).balanceOf(destAddress).then((e) => {
        console.log(e);
        console.log(e.toNumber())
        callback();
    }).catch(e => {
        console.log("errr", e);

        callback();
    });

}