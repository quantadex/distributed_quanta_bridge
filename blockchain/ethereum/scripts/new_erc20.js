const SimpleToken = artifacts.require("SimpleToken");

module.exports = function(callback) {
    return SimpleToken.new().then((erc20) => {
      console.log(erc20.address)
    })
}