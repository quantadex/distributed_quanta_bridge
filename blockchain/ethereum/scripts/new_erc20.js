const SimpleToken = artifacts.require("SimpleToken");

module.exports = function(cb) {
    return SimpleToken.new().then((erc20) => {
      console.log(erc20.address)
      cb()
    })
}