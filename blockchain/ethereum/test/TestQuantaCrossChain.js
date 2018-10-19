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
const SimpleToken = artifacts.require("SimpleToken");

const Helpers = require('./helpers.js');
const TrufAssert = require('truffle-assertions');
const Web3Utils = require('web3-utils');

const catchAssert = require("./exceptions.js").catchInvalidOpcode;  // asserts throw invalid op code?
const catchRevert = require("./exceptions.js").catchRevert;

async function grantEther(contract, wei) {
  await contract.sendTransaction({from:web3.eth.coinbase,value:wei});
}


async function fetchCurrentTxId(contract) {
  var txId = Number(await contract.txIdLast());
  assert(txId >= 0);
  return txId;
}


async function quickPayment(contract, signer, erc20Addr, toAddr, amount) {
  var txId = await fetchCurrentTxId(contract);
  var nextTxId = txId + 1;

  var vrs = Helpers.makeVRS(
    signer, nextTxId,  erc20Addr, toAddr, amount, debug=false);
  let result = await contract.paymentTx(
    nextTxId, erc20Addr, toAddr, amount,
    [vrs[0]],
    [vrs[1]],
    [vrs[2]]);

    return result;
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
  var contract;

  it("should deploy our contract", async () => {
    contract = await QuantaCrossChain.deployed();
  });

  it("should revert paymentTx with no sigs", async () => {
    await catchRevert(
      contract.paymentTx(1, 0, accounts[2], 1, [], [], []),
    );
  });

  it("should revert paymentTx with no sigs and zeros", async () => {
    await catchRevert(
      contract.paymentTx(0, 0, 0, 0, [], [], []),
    );
  });

  it("should have txId==0", async () => {
    var txId = await fetchCurrentTxId(contract);
    assert.equal(txId, 0);
  });
})

contract('QuantaCrossChain one signer', async (accounts) => {
  var contract;

  it("should assign the one and only signer", async () => {
    contract = await QuantaCrossChain.deployed();
    await contract.assignInitialSigners([accounts[0]]);
  });

  it("should revert paymentTx with no sigs", async () => {
    await catchRevert(
      contract.paymentTx(0, 0, 0, 0, [], [], []),
    );
  });

  it("should fail paymentTx with bad sig", async () => {
    var txId = await fetchCurrentTxId(contract);
    assert(txId == 0);
    var nextTxId = txId + 1;
    let result = await contract.paymentTx(nextTxId, 0, 0, 0, [0], [0], [0]);

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
    var txId = await fetchCurrentTxId(contract);
    var nextTxId = txId + 1;
    var amount = web3.toWei(1, "ether");
    var toAddr = accounts[2];

    await grantEther(contract, amount);

    var vrs = Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);
    let result = await contract.paymentTx(
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
    var txId = await fetchCurrentTxId(contract);
    var nextTxId = txId + 1;
    var amount = 1;
    var toAddr = accounts[2];

    await grantEther(contract, amount);

    var vrs = Helpers.makeVRS(
      accounts[0], nextTxId,  0, toAddr, amount);
    let result = await contract.paymentTx(
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
    await catchRevert(
      contract.paymentTx(0, 0, 0, 0, [], [], []),
    );
  });

  it("should allow paymentTx with correct sig [erc20]", async () => {
    let erc20 = await SimpleToken.new();

    var txId = await fetchCurrentTxId(contract);
    var nextTxId = txId + 1;
    var amount = 10;
    var toAddr = accounts[2];

    await erc20.transfer(contract.address, amount);  // give our contract enough tokens

    var vrs = Helpers.makeVRS(
      accounts[0], nextTxId,  erc20.address, toAddr, amount, debug=false);
    let result = await contract.paymentTx(
      nextTxId, erc20.address, toAddr, amount,
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
        (ev.erc20Addr == erc20.address) &&
        (ev.toAddr == toAddr)
      );
    });
  });
})


contract('QuantaCrossChain two signers', async (accounts) => {
  var contract;

  it("should assign the two initial signers", async () => {
    contract = await QuantaCrossChain.deployed();
    await contract.assignInitialSigners([accounts[0], accounts[1]]);
  });

  it("should fail paymentTx fast with only one sig [native eth]", async () => {
    var txId = await fetchCurrentTxId(contract);
    assert(txId == 0);
    var nextTxId = txId + 1;

    var amount = 1;
    var toAddr = accounts[2];

    var vrs = Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);
    let result = await contract.paymentTx(
      nextTxId, 0, toAddr, amount,
      [vrs[0]],
      [vrs[1]],
      [vrs[2]]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        !ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == false) &&  // should have skipped
        (ev.txId == txId) &&  // should not advance
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });

    txId = await fetchCurrentTxId(contract);
    assert(txId == 0);
  });

  it("should allow paymentTx with both sigs", async () => {
    var txId = await fetchCurrentTxId(contract);
    var nextTxId = txId + 1;

    var amount = 6;
    var toAddr = accounts[2];

    var vrs0 = Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);
    var vrs1 = Helpers.makeVRS(
      accounts[1], nextTxId,  null, toAddr, amount);

    await grantEther(contract, amount);

    let result = await contract.paymentTx(
      nextTxId, 0, toAddr, amount,
      [vrs0[0], vrs1[0]],
      [vrs0[1], vrs1[1]],
      [vrs0[2], vrs1[2]]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 2) &&
        (ev.verified[0] == true) &&
        (ev.verified[1] == true) &&
        (ev.txId == nextTxId) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });

    txId = await fetchCurrentTxId(contract);
    assert(txId == nextTxId);
  });

  // no need to test this scenario since sig validation is separate from currency transfer
  // it("should allow paymentTx with both sigs [erc20]"

  it("should not allow paymentTx with dupe sig [native eth]", async () => {
    var txId = await fetchCurrentTxId(contract);
    var nextTxId = txId + 1;

    var amount = 7;
    var toAddr = accounts[2];

    var vrs0 = Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);

    await grantEther(contract, amount);

    let result = await contract.paymentTx(
      nextTxId, 0, toAddr, amount,
      [vrs0[0], vrs0[0]],  // try to trick the contract by giving dupe sigs
      [vrs0[1], vrs0[1]],
      [vrs0[2], vrs0[2]]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        !ev.success &&
        (ev.verified.length == 2) &&
        (ev.verified[0] == true) &&
        (ev.verified[1] == false) &&  // second dupe sig should not pass
        (ev.txId == txId) &&  // should not advance
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });

    txId = await fetchCurrentTxId(contract);
    assert(txId == txId);  // should be the old txId
  });

  it("should allow paymentTx with both sigs and a dupe sig", async () => {
    var txId = await fetchCurrentTxId(contract);
    var nextTxId = txId + 1;

    var amount = 6;
    var toAddr = accounts[2];

    var vrs0 = Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);
    var vrs1 = Helpers.makeVRS(
      accounts[1], nextTxId,  null, toAddr, amount);

    await grantEther(contract, amount);

    let result = await contract.paymentTx(
      nextTxId, 0, toAddr, amount,
      [vrs0[0], vrs1[0], vrs0[0]],
      [vrs0[1], vrs1[1], vrs0[1]],
      [vrs0[2], vrs1[2], vrs0[2]]);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 3) &&
        (ev.verified[0] == true) &&
        (ev.verified[1] == true) &&
        (ev.verified[2] == false) &&
        (ev.txId == nextTxId) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });

    txId = await fetchCurrentTxId(contract);
    assert(txId == nextTxId);
  });

  it("should fail paymentTx fast with bad sig [native eth]", async () => {
    var txId = await fetchCurrentTxId(contract);
    var nextTxId = txId + 1;

    var amount = 8;
    var toAddr = accounts[2];

    var vrs0 = Helpers.makeVRS(
      accounts[0], nextTxId,  null, toAddr, amount);
    var vrsX = Helpers.makeVRS(
      accounts[2], nextTxId,  null, toAddr, amount);  // accounts[2] is invalid

    await grantEther(contract, amount);

    let result = await contract.paymentTx(
      nextTxId, 0, toAddr, amount,
      [vrsX[0], vrs0[0]],
      [vrsX[1], vrs0[1]],
      [vrsX[2], vrs0[2]]);

    // should fail AND not bother to validate the second (valid) signature
    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        !ev.success &&
        (ev.verified.length == 2) &&
        (ev.verified[0] == false) &&
        (ev.verified[1] == false) &&
        (ev.txId == txId) &&  // should not advance
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });

    txId = await fetchCurrentTxId(contract);
    assert(txId == txId);  // should be the old txId
  });
});


contract('QuantaCrossChain one signer [erc20 balance cases]', async (accounts) => {
  var contract;

  it("should assign the one and only signer", async () => {
    contract = await QuantaCrossChain.deployed();
    await contract.assignInitialSigners([accounts[0]]);
  });

  it("should revert payment when not enough erc20 balance", async () => {
    let erc20 = await SimpleToken.new();

    var amount = 10;
    var toAddr = accounts[2];
    var balance = await erc20.balanceOf(contract.address);
    assert.equal(balance, 0);

    await catchRevert(
      quickPayment(contract, accounts[0], erc20.address, toAddr, amount),
    );

    await erc20.transfer(contract.address, 8);  // give our contract not enough tokens
    await catchRevert(
      quickPayment(contract, accounts[0], erc20.address, toAddr, amount),
    );
    balance = await erc20.balanceOf(contract.address);
    assert.equal(balance, 8);
  });

  it("should allow payment when enough erc20 balance", async () => {
    let erc20 = await SimpleToken.new();

    var amount = 10;
    var toAddr = accounts[2];
    var balance = await erc20.balanceOf(contract.address);
    assert.equal(balance, 0);

    await erc20.transfer(contract.address, 8);  // give our contract not enough tokens
    await erc20.transfer(contract.address, 2);  // 8 + 2 = 10
    balance = await erc20.balanceOf(contract.address);
    assert.equal(balance, 10);

    let result = await quickPayment(contract, accounts[0], erc20.address, toAddr, amount);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == true) &&
        (ev.txId == 1) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == erc20.address) &&
        (ev.toAddr == toAddr)
      );
    });

    balance = await erc20.balanceOf(contract.address);
    assert.equal(balance, 0);  // check remainder
  });

  it("should allow payment when too much erc20 balance", async () => {
    let erc20 = await SimpleToken.new();

    var txId = Number(await contract.txIdLast());
    var amount = 10;
    var toAddr = accounts[2];

    await erc20.transfer(contract.address, amount+5);
    var balance = await erc20.balanceOf(contract.address);
    assert.equal(balance, amount+5);

    let result = await quickPayment(contract, accounts[0], erc20.address, toAddr, amount);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == true) &&
        (ev.txId == txId+1) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == erc20.address) &&
        (ev.toAddr == toAddr)
      );
    });

    balance = await erc20.balanceOf(contract.address);
    assert.equal(balance, 5);  // check remainder
  });
});


contract('QuantaCrossChain one signer [native eth balance cases]', async (accounts) => {
  var contract;

  it("should assign the one and only signer", async () => {
    contract = await QuantaCrossChain.deployed();
    await contract.assignInitialSigners([accounts[0]]);
  });

  it("should revert payment when not enough native eth balance", async () => {
    var amount = 10;
    var toAddr = accounts[2];
    var balance = await web3.eth.getBalance(contract.address);
    assert.equal(balance, 0);

    await catchRevert(
      quickPayment(contract, accounts[0], null, toAddr, amount),
    );

    await grantEther(contract, 8);  // give our contract not enough ether
    await catchRevert(
      quickPayment(contract, accounts[0], null, toAddr, amount),
    );
    balance = await web3.eth.getBalance(contract.address);
    assert.equal(balance, 8);
  });

  it("should allow payment when enough native eth balance", async () => {
    var amount = 10;
    var toAddr = accounts[2];
    var balance = await web3.eth.getBalance(contract.address);
    assert.equal(balance, 8);   // we should already have 8 from the previous test case

    await grantEther(contract, 2);  // 8 + 2 = 10
    balance = await web3.eth.getBalance(contract.address);
    assert.equal(balance, 10);

    let result = await quickPayment(contract, accounts[0], null, toAddr, amount);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == true) &&
        (ev.txId == 1) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });

    balance = await web3.eth.getBalance(contract.address);
    assert.equal(balance, 0);
  });

  it("should allow payment when too much native eth balance", async () => {
    var txId = Number(await contract.txIdLast());
    var amount = 10;
    var toAddr = accounts[2];

    await grantEther(contract, amount+5);
    var balance = await web3.eth.getBalance(contract.address);
    assert.equal(balance, amount+5);

    let result = await quickPayment(contract, accounts[0], null, toAddr, amount);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        ev.success &&
        (ev.verified.length == 1) &&
        (ev.verified[0] == true) &&
        (ev.txId == txId+1) &&
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });

    balance = await web3.eth.getBalance(contract.address);
    assert.equal(balance, 5);  // check remainder
  });
});


contract('QuantaCrossChain assign initial signers', async (accounts) => {
  var contract;

  it("should deploy the contract", async () => {
    contract = await QuantaCrossChain.deployed();

    let totalSigners = await contract.totalSigners();
    assert.equal(0, totalSigners);
  });

  it("should revert paymentTx with no sigs", async () => {
    await catchRevert(
      contract.paymentTx(0, 0, 0, 0, [], [], []),
    );
  });

  it("should revert with an empty assign signers list", async () => {
    contract = await QuantaCrossChain.deployed();
    await catchRevert(
      contract.assignInitialSigners([]),
    );

    let totalSigners = await contract.totalSigners();
    assert.equal(0, totalSigners);
  });

  it("should assign the first signer", async () => {
    await contract.assignInitialSigners([accounts[0]]);

    let totalSigners = await contract.totalSigners();
    assert.equal(1, totalSigners);
  });

  it("should assert failure on a second assign initial signers call", async () => {
    await catchAssert(
      contract.assignInitialSigners([]),
    )

    await catchAssert(
      contract.assignInitialSigners([accounts[0], accounts[1]]),
    )

    await catchAssert(
      contract.assignInitialSigners([]),
    )

    let totalSigners = await contract.totalSigners();
    assert.equal(1, totalSigners);
  });

  it("should allow payment", async () => {
    var amount = 10;
    var toAddr = accounts[2];

    await grantEther(contract, amount);
    let result = await quickPayment(contract, accounts[0], null, toAddr, amount);

    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return ev.success;
    });
  });
});


contract('QuantaCrossChain 7 signers', async (accounts) => {
  var contract;

  it("should deploy the contract", async () => {
    contract = await QuantaCrossChain.deployed();

    await contract.assignInitialSigners([accounts[0], accounts[1], accounts[2], accounts[3], accounts[4], accounts[5], accounts[6]]);

    let totalSigners = await contract.totalSigners();
    assert.equal(7, totalSigners);
  });

  it("should fail paymentTx fast with bad sig [native eth]", async () => {
    var txId = await fetchCurrentTxId(contract);
    var nextTxId = txId + 1;

    var amount = 8;
    var toAddr = accounts[2];

    await grantEther(contract, amount);

    var arrV = [];
    var arrR = [];
    var arrS = [];

    var vrs;
    vrs = Helpers.makeVRS(accounts[0], nextTxId,  null, toAddr, amount); arrV.push(vrs[0]); arrR.push(vrs[1]); arrS.push(vrs[2]);
    vrs = Helpers.makeVRS(accounts[1], nextTxId,  null, toAddr, amount); arrV.push(vrs[0]); arrR.push(vrs[1]); arrS.push(vrs[2]);

    // #7 is a trap
    vrs = Helpers.makeVRS(accounts[7], nextTxId,  null, toAddr, amount); arrV.push(vrs[0]); arrR.push(vrs[1]); arrS.push(vrs[2]);

    // skip accounts[2]
    vrs = Helpers.makeVRS(accounts[3], nextTxId,  null, toAddr, amount); arrV.push(vrs[0]); arrR.push(vrs[1]); arrS.push(vrs[2]);
    vrs = Helpers.makeVRS(accounts[4], nextTxId,  null, toAddr, amount); arrV.push(vrs[0]); arrR.push(vrs[1]); arrS.push(vrs[2]);
    vrs = Helpers.makeVRS(accounts[5], nextTxId,  null, toAddr, amount); arrV.push(vrs[0]); arrR.push(vrs[1]); arrS.push(vrs[2]);
    vrs = Helpers.makeVRS(accounts[6], nextTxId,  null, toAddr, amount); arrV.push(vrs[0]); arrR.push(vrs[1]); arrS.push(vrs[2]);

    let result = await contract.paymentTx(
      nextTxId, 0, toAddr, amount,
      arrV, arrR, arrS);

    // should fail AND not bother to validate the second (valid) signature
    TrufAssert.eventEmitted(result, 'TransactionResult', (ev) => {
      return (
        !ev.success &&
        (ev.verified.length == 7) &&
        (ev.verified[0] == true) &&  // it should of made it through the first two
        (ev.verified[1] == true) &&
        (ev.verified[2] == false) &&  // this should have been invalid
        (ev.verified[3] == false) &&  // the remainder should be skipped
        (ev.verified[4] == false) &&  // the remainder should be skipped
        (ev.verified[5] == false) &&  // the remainder should be skipped
        (ev.verified[6] == false) &&  // the remainder should be skipped
        (ev.txId == txId) &&  // should not advance
        (ev.amount == amount) &&
        (ev.erc20Addr == 0) &&
        (ev.toAddr == toAddr)
      );
    });

    txId = await fetchCurrentTxId(contract);
    assert(txId == txId);  // should be the old txId
  });
});
