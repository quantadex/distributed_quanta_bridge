/* Run this test in `truffle console`
*  1. run ganache with fixed block time:
*     ganache-cli -m "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat" -u 0 -u 1 --gasLimit 0x2FEFD800000 -a 100 --defaultBalanceEther 10000
*
*  2. Run truffle console in another terminal:
*     truffle console
*
*  3. Run this test:
*     test test/TestQuantaCrossChain.js
*
*  The test will take about 59 seconds (468931ms for buying tickets)
* */

const QuantaCrossChain = artifacts.require("QuantaCrossChain");

const Helpers = require('./helpers.js');
const TrufAssert = require('truffle-assertions');
const Web3Utils = require('web3-utils');

const catchRevert = require("./exceptions.js").catchRevert;


async function grantEther(contract, wei) {
  await contract.sendTransaction({from:web3.eth.coinbase,value:wei});
}


async function fetchCurrentTxId(contract) {
  var txId = Number(await contract.txIdLast());
  assert(txId >= 0);
  return txId;
}


contract('test signer', async (accounts) => {
  it("it should generate the correct signature message", async () => {
    let unsignedMsg = Helpers.toQuantaPaymentSigMsg(
      1,
      "0xf17f52151ebef6c7334fad080c5704d77216b732",
      "0xc5fdf4076b8f3a5357c5e395ab970b5b54098fef",
      1,
      debug=false, preamble=false);

    // https://github.com/ethereum/web3.js/blob/v1.0.0-beta.35/packages/web3-eth-accounts/src/index.js#L246
    var message = Web3Utils.hexToBytes(unsignedMsg);
    var messageBuffer = Buffer.from(message);
    var preamble = "\x19Ethereum Signed Message:\n" + message.length;
    var preambleBuffer = Buffer.from(preamble);
    var ethMessage = Buffer.concat([preambleBuffer, messageBuffer]);

    let hexMsg = "0x" + ethMessage.toString("hex");

    assert.equal(
      hexMsg,
      "0x19457468657265756d205369676e6564204d6573736167653a0a38300000000000000001f17f52151ebef6c7334fad080c5704d77216b732c5fdf4076b8f3a5357c5e395ab970b5b54098fef0000000000000000000000000000000000000000000000000000000000000001");
  });

  it("it should generate the correct preamble message", async () => {
    let msg = Helpers.toQuantaPaymentSigMsg(1, "0xf17f52151ebef6c7334fad080c5704d77216b732", "0xc5fdf4076b8f3a5357c5e395ab970b5b54098fef", 1, debug=false, preamble=true);
    assert.equal(
      msg,
      "0x19457468657265756d205369676e6564204d6573736167653a0a38300000000000000001f17f52151ebef6c7334fad080c5704d77216b732c5fdf4076b8f3a5357c5e395ab970b5b54098fef0000000000000000000000000000000000000000000000000000000000000001");
  });
})


contract('QuantaCrossChain no signers', async (accounts) => {
  it("should revert paymentTx with no sigs", async () => {
    let inst = await QuantaCrossChain.deployed();

    await catchRevert(
      inst.paymentTx(1, 0, accounts[2], 1, [], [], []),
    );
  });

  it("should revert paymentTx with no sigs and zeros", async () => {
    let inst = await QuantaCrossChain.deployed();

    await catchRevert(
      inst.paymentTx(0, 0, 0, 0, [], [], []),
    );
  });

  it("should have txId==0", async () => {
    let inst = await QuantaCrossChain.deployed();

    var txId = await fetchCurrentTxId(inst);
    assert.equal(txId, 0);
  });
})


contract('QuantaCrossChain one signer', async (accounts) => {
  it("should revert paymentTx with no sigs", async () => {
    let inst = await QuantaCrossChain.deployed();
    inst.voteAddSigner(accounts[0]);

    await catchRevert(
      inst.paymentTx(0, 0, 0, 0, [], [], []),
    );
  });

  it("should fail paymentTx with bad sig", async () => {
    let inst = await QuantaCrossChain.deployed();
    inst.voteAddSigner(accounts[0]);

    var txId = await fetchCurrentTxId(inst);
    assert(txId == 0);
    var nextTxId = txId + 1;
    let result = await inst.paymentTx(nextTxId, 0, 0, 0, [0], [0], [0]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        !ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == false) &&
        (ev.txId == txId) &&  // should not advance
        (ev.amount == 0) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == 0)
      );
    });
  });

  it("should allow paymentTx with correct sig [native eth]", async () => {
    let inst = await QuantaCrossChain.deployed();
    inst.voteAddSigner(accounts[0]);

    var txId = await fetchCurrentTxId(inst);
    var nextTxId = txId + 1;
    var amount = web3.toWei(1, "ether");
    var toAddr = accounts[2];

    await grantEther(inst, amount);

    var vrs = await Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);
    let result = await inst.paymentTx(
      nextTxId, 0, toAddr, amount,
      [vrs[0]],
      [vrs[1]],
      [vrs[2]]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == true) &&
        (ev.txId == nextTxId) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });
  });

  it("should allow paymentTx with extra sig [native eth]", async () => {
    let inst = await QuantaCrossChain.deployed();
    inst.voteAddSigner(accounts[0]);

    var txId = await fetchCurrentTxId(inst);
    var nextTxId = txId + 1;
    var amount = 1;
    var toAddr = accounts[2];

    await grantEther(inst, amount);

    var vrs = await Helpers.makeVRS(
      accounts[0], nextTxId,  0, toAddr, amount);
    let result = await inst.paymentTx(
      nextTxId, null, toAddr, amount,
      [vrs[0], vrs[0]],
      [vrs[1], vrs[1]],
      [vrs[2], vrs[2]]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 2) &&
        (ev.verified[0] == true) &&
        (ev.verified[1] == false) &&
        (ev.txId == nextTxId) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });
  });

  it("should revert paymentTx with no sigs", async () => {
    let inst = await QuantaCrossChain.deployed();
    inst.voteAddSigner(accounts[0]);

    await catchRevert(
      inst.paymentTx(0, 0, 0, 0, [], [], []),
    );
  });

  /*
  it("should allow paymentTx with correct sig [erc20]", async () => {
    let inst = await QuantaCrossChain.deployed();

    // var erc20 = await SimpleERC20Token.new({gas: 1000000000, gasPrice: 1});
    let erc20 = await SimpleERC20Token.deployed();

    console.log("erc20: " + erc20);
    inst.voteAddSigner(accounts[0]);

    var txId = await fetchCurrentTxId(inst);
    var nextTxId = txId + 1;
    var amount = web3.toWei(1, "ether");
    var toAddr = accounts[2];

    await grantEther(inst, amount);

    var vrs = await Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);
    let result = await inst.paymentTx(
      nextTxId, erc20Addr, toAddr, amount,
      [vrs[0]],
      [vrs[1]],
      [vrs[2]]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == true) &&
        (ev.txId == nextTxId) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == erc20Addr) &&
        (ev.toAddr == toAddr)
      );
    });
  });
  */
})


contract('QuantaCrossChain two signers', async (accounts) => {
  it("should not allow paymentTx with only one sig [native eth]", async () => {
    let inst = await QuantaCrossChain.deployed();
    inst.voteAddSigner(accounts[0]);
    inst.voteAddSigner(accounts[1]);

    var txId = await fetchCurrentTxId(inst);
    assert(txId == 0);
    var nextTxId = txId + 1;

    var amount = 1;
    var toAddr = accounts[2];

    var vrs = await Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);
    let result = await inst.paymentTx(
      nextTxId, 0, toAddr, amount,
      [vrs[0]],
      [vrs[1]],
      [vrs[2]]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        !ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == true) &&
        (ev.txId == txId) &&  // should not advance
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });
  });
})
