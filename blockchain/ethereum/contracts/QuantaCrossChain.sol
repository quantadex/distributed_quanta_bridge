pragma solidity ^0.4.24;

import { libbytes } from "../libraries/libbytes.sol";
import { ECTools } from "../libraries/ECTools.sol";


// TODO: make this contract ownable

contract QuantaCrossChain {
  /** keys are the signers' address, value is a 1 if active */
  mapping(address=>uint8) private signers;

  /** the number of total signers added */
  uint256 private totalSigners = 0;

  /** the last successfully used txId */
  uint64 public txIdLast = 0;

  event TransactionResult(bool success,
                          uint64 txId,
                          address erc20Addr,
                          address toAddr,
                          uint256 amount,
                          bool[] verified);


  function paymentTx(uint64 txId,
                     address erc20Addr,
                     address toAddr,
                     uint256 amount,
                     uint8[] v, bytes32[] r, bytes32[] s) public {
    uint n = v.length;

    require(n != 0);
    require(n == r.length);
    require(n == s.length);
    require(txId == txIdLast+1);

    bool[] memory verified = new bool[](n);
    bytes memory sigMsg = toQuantaPaymentSignatureMessage(txId, erc20Addr, toAddr, amount);

    bytes32 signed = ECTools.toEthereumSignedMessage(string(sigMsg));

    address addr;

    // TODO: check for duplicates
    // TODO: many to one?
    for(uint i=0; i<n; i++) {
      addr = ECTools.recoverSignerVRS(signed, v[i], r[i], s[i]);
      if (signers[addr] == 1) {
        verified[i] = true;
        n--;
      }
    }

    if ((v.length-n) == totalSigners) {
      // TODO: check sufficient balances

      if (erc20Addr == 0) {
        // https://solidity.readthedocs.io/en/v0.4.24/units-and-global-variables.html#address-related
        toAddr.transfer(amount);
      } else {
        // https://theethereum.wiki/w/index.php/ERC20_Token_Standard#The_ERC20_Token_Standard_Interface

        // TODO: test with FixedSupply ERC20 sample contract
        // TODO: replace erc20 copy from zeppelin v1.12.0
        // https://github.com/OpenZeppelin/openzeppelin-solidity/blob/v1.12.0/contracts/token/ERC20/ERC20.sol
        // https://medium.com/taipei-ethereum-meetup/smart-contract-unit-testing-use-erc20-token-contract-as-an-example-d150c2700834

        // TODO: use https://github.com/OpenZeppelin/openzeppelin-solidity/blob/v1.12.0/contracts/token/ERC20/SafeERC20.sol
        revert("not supported, erc20 transfers");
        // ERC20(erc20Addr).transfer(toAddr, amount);
      }
    }

     // advance the txId
     if ((v.length-n) == totalSigners) {
       txIdLast++;
     }

     emit TransactionResult((v.length-n) == totalSigners, txIdLast, erc20Addr, toAddr, amount, verified);  // , v, r, s, sigMsg);
  }

  function voteAddSigner(address signer) public {
    require(signer != 0x0);

    if (signers[signer] == 0) {
      totalSigners++;
      signers[signer] = 1;
    }
  }

  function voteRemoveSigner(address signer) public {
    require(signer != 0x0);

    if (signers[signer] == 1) {
      delete signers[signer];
      totalSigners--;
    }
  }

  /**
   * @return hex encoded signature message
   */
  function toQuantaPaymentSignatureMessage(uint64 txId,
                                           address erc20Addr,
                                           address toAddr,
                                           uint256 amount) internal pure returns (bytes) {
     // preSignedMsg := preamble + header + body
     // preamble := '\x19Ethereum Signed Message:\n'
     // header := length(body)
     //
     //     bytes /     bits | param
     // -------------------------------------------------------
     //  26 bytes / 208 bits | '\x19Ethereum Signed Message:\n'
     //   4 bytes /  32 bits | msg body length
     //   8 bytes /  64 bits | txId
     //  20 bytes / 160 bits | erc20Addr
     //  20 bytes / 160 bits | toAddr
     //  32 bytes / 256 bits | amount
     // -------------------------------------------------------
     // 110 bytes / 848 bits | total

     bytes8 b_txId = bytes8(txId);  // libbytes.uint64ToBytes(txId);
     bytes memory b_erc20Addr = libbytes.addressToBytes(erc20Addr);
     bytes memory b_toAddr = libbytes.addressToBytes(toAddr);
     bytes32 b_amount = bytes32(amount);

     string memory s = new string(b_txId.length + b_erc20Addr.length + b_toAddr.length + b_amount.length);
     bytes memory b_msg = bytes(s);
     uint k = 0;
     for (uint i = 0; i < b_txId.length; i++) b_msg[k++] = b_txId[i];
     for (i = 0; i < b_erc20Addr.length; i++) b_msg[k++] = b_erc20Addr[i];
     for (i = 0; i < b_toAddr.length; i++) b_msg[k++] = b_toAddr[i];
     for (i = 0; i < b_amount.length; i++) b_msg[k++] = b_amount[i];

     return b_msg;
  }

  // https://ethereum.stackexchange.com/questions/42995/how-to-send-ether-to-a-contract-in-truffle-test/43011
  // @notice Logs the address of the sender and amounts paid to the contract
  event Fund(address indexed _from, uint _value);

  // @notice Will receive any eth sent to the contract
  function () external payable {
    emit Fund(msg.sender, msg.value);
  }
}
