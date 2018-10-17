const QUANTA_CROSS_CHAIN = artifacts.require("./QuantaCrossChain.sol");
const LIB_BYTES = artifacts.require("../libraries/libbytes.sol");
const LIB_ECTOOLS = artifacts.require("../libraries/ECTools.sol");


module.exports = function(deployer) {
	deployer.deploy(LIB_BYTES);
	deployer.link(LIB_BYTES, QUANTA_CROSS_CHAIN);

	deployer.deploy(LIB_ECTOOLS);
	deployer.link(LIB_ECTOOLS, QUANTA_CROSS_CHAIN);

	deployer.deploy(QUANTA_CROSS_CHAIN);
};
