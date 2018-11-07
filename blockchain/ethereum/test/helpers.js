const BN = require('bn.js')
const ByteBuffer = require("bytebuffer");
const Utils = require('ethereumjs-util')
const Web3Utils = require('web3-utils');
const assertjs = require('assert');

const BIG_ENDIAN = 'be';
const ADDRESS_0X0 = '0x0000000000000000000000000000000000000000'
const ADDRESS_0X0_PART = '0000000000000000000000000000000000000000'


function _toQuantaPaymentSigMsg(txId, erc20Addr, toAddr, amount, debug=false, preamble=false) {
  if (debug) console.log(`{ "txId": ${txId}, "erc20Addr": ${erc20Addr}, "toAddr": ${toAddr}, "amount": ${amount}}`);

  var erc20AddrPart = ADDRESS_0X0_PART;
  if (typeof erc20Addr == 'string') {
    assertjs(Utils.isHexPrefixed(erc20Addr));
    erc20AddrPart = Utils.stripHexPrefix(erc20Addr);
  }

  assertjs(typeof txId == 'number');
  assertjs(typeof toAddr == 'string');
  assertjs(Utils.isHexPrefixed(toAddr));
  assertjs((typeof amount == 'number') || (typeof amount == 'string'))
  // amount must be a number or a string, i.e. for big numbers like "1000000000000000000"

  var totalBytes = 0;
  var elems = [
    new BN(txId).toArrayLike(Buffer, BIG_ENDIAN, 8),
    new BN(erc20AddrPart, 16).toArrayLike(Buffer, BIG_ENDIAN, 20),
    new BN(Utils.stripHexPrefix(toAddr), 16).toArrayLike(Buffer, BIG_ENDIAN, 20)
  ];

  var bnPart;
  if (typeof amount == 'string') {
    bnPart = new BN(amount, 10);  // assume a Base10 number
  } else if (typeof amount == 'number') {
    bnPart = new BN(amount);
  } else {
    assertjs(false);
  }
  elems.push(bnPart.toArrayLike(Buffer, BIG_ENDIAN, 32));

  for (i=0; i< elems.length; i++) {
    var elem = elems[i];
    totalBytes += elem.length;

    if (debug) {
      console.log("[Buffer #" + (i+1) + "/" + elems.length + " bytes="
                  + elem.length + " (" + (elem.length*8) + " bits)] "
                  + elem.toString('hex'));
    }
  }

  var buf = Buffer.concat(elems, totalBytes);
  var msg = buf.toString('hex');

  assertjs(buf.length == totalBytes);

  if (debug) {
    console.log("totalBytes=" + totalBytes + " (" + (totalBytes*8) + " bits)");
    console.log("buf[" + buf.length + " bytes / " + msg.length + " hex] " + msg);
  }

  if (preamble) {
    var preambleBuf = Buffer.from('\x19Ethereum Signed Message:\n' + buf.length);
    totalBytes += preambleBuf.length;

    buf = Buffer.concat([preambleBuf, buf], totalBytes);
    msg = buf.toString('hex');

    if (debug) {
      console.log("totalBytesWithPreamble=" + totalBytes + " (" + (totalBytes*8) + " bits)");
      console.log("bufWithPreamble[" + buf.length + " bytes / " + msg.length + " hex] " + msg);
    }
  }

  return "0x" + msg;
}


/** signes the appropriate msg and returns the embedded v r s */
module.exports.makeVRS = function(account, txId, erc20Addr, toAddr, amount, debug=false) {
  let msg = _toQuantaPaymentSigMsg(txId, erc20Addr, toAddr, amount, debug=debug);
  let sig1 = web3.eth.sign(account, msg);
  var sig = sig1.slice(2);
  var r = `0x${sig.slice(0, 64)}`;
  var s = `0x${sig.slice(64, 128)}`;
  var v = Web3Utils.toDecimal(sig.slice(128, 130)) + 27;

  // https://github.com/dcodeIO/bytebuffer.js
  var bb = new ByteBuffer(ByteBuffer.DEFAULT_CAPACITY, ByteBuffer.BIG_ENDIAN);

  bb.append(r, "binary").flip();
  var br = bb.toString("binary");

  bb.clear();
  bb.append(s, "binary").flip();
  var bs = bb.toString("binary");

  var bv = v;  // don't need fancy ByteBuffer for integers

  return [bv, br, bs];
}


/**
 * Generates the Quanta Payment message to be signed from the input parameters.
 *
 * @param preamble if true, then the message will be prepended with the
 *                 '\x19Ethereum Signed Message:\n{messageLength}' string
 *                 (used for debugging and diagnostics)
 *
 * @return string - the hex encoded message
 */
module.exports.toQuantaPaymentSigMsg = function(txId, erc20Addr, toAddr, amount, debug=false, preamble=false) {
  return _toQuantaPaymentSigMsg(txId, erc20Addr, toAddr, amount, debug, preamble);
}
