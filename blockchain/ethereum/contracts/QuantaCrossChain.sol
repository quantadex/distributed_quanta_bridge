pragma solidity ^0.4.24;

import { Ownable } from "./zeppelin/ownership/Ownable.sol";
import { ERC20 } from "./zeppelin/token/ERC20.sol";
import { libbytes } from "../libraries/libbytes.sol";
import { ECTools } from "../libraries/ECTools.sol";


/**
 * Quanta Cross Chain contract.
 *
 * Initial deployment will instantiate an unusable contract.
 * To make the contract usable, the owner must call assignInitialSigners() once.
 */
contract QuantaCrossChain is Ownable {
  /** keys are the signers' address, value is a 1 if active */
  mapping(address=>uint8) private signers;

  /** the number of total signers ratified */
  uint256 private totalSigners = 0;

  /** the last successfully used txId */
  uint64 public txIdLast = 0;

  event TransactionResult(bool success,
                          uint64 txId,
                          address erc20Addr,
                          address toAddr,
                          uint256 amount,
                          bool[] verified);

  /**
   * Assigns the initial list of signers.
   * Must be called once after contract instantiation.
   * All other methods will fail or revert until this is called.
   * May only be called once by the owner. Further calls are ignored.
  */
  function assignInitialSigners(address[] initialSigners) external onlyOwner {
    assert(totalSigners == 0);  // penalize (using assertions) if tried twice
    require(initialSigners.length > 0);

    // TODO: should we temporarily set totalSigners to MAX_UINT256?

    for(uint i=0; i<initialSigners.length; i++) {
      signers[initialSigners[i]] = 1;
    }

    totalSigners = i;
  }

  function getTotalSigners() public view onlyOwner returns (uint256) {
    return totalSigners;
  }

  function paymentTx(uint64 txId,
                     address erc20Addr,
                     address toAddr,
                     uint256 amount,
                     uint8[] v, bytes32[] r, bytes32[] s) external {  // FIXME: use external instead of public since it uses last gas
    uint n = v.length;

    require(n != 0);
    require(n == r.length);
    require(n == s.length);
    require(txId == txIdLast+1);

    bytes memory sigMsg = toQuantaPaymentSignatureMessage(txId, erc20Addr, toAddr, amount);

    bool[] memory verified = new bool[](n);
    bool success = validateSignatures(sigMsg, n, v, r, s, verified);

    if (success) {
      // if insufficient balance, methods will just revert

      if (erc20Addr == 0) {
        // https://solidity.readthedocs.io/en/v0.4.24/units-and-global-variables.html#address-related
        toAddr.transfer(amount);
      } else {
        // https://theethereum.wiki/w/index.php/ERC20_Token_Standard#The_ERC20_Token_Standard_Interface
        ERC20(erc20Addr).transfer(toAddr, amount);
      }

      txIdLast++;
    }

    emit TransactionResult(success, txIdLast, erc20Addr, toAddr, amount, verified);
  }


  function validateSignatures(bytes sigMsg,
                              uint numSigs,
                              uint8[] v,
                              bytes32[] r,
                              bytes32[] s,
                              bool[] outVerified) internal view returns (bool) {
    bytes32 signed = ECTools.toEthereumSignedMessage(string(sigMsg));
    address[] memory validated = new address[](numSigs);
    uint numValidated = 0;

    for(uint i=0; i<numSigs; i++) {
      address addr = ECTools.recoverSignerVRS(signed, v[i], r[i], s[i]);
      if (signers[addr] == 1) {
        // wen anticipate the number of signers to be small, so just do a linear scan
        for(uint j=0; j<numValidated; j++) {
          if (validated[j] == addr) {
            addr = 0x0;
            break;
          }
        }
        if (addr != 0x0) {
          outVerified[i] = true;
          validated[numValidated] = addr;
          numValidated++;
        }
      }
    }

    return numValidated == totalSigners;
  }


  function voteAddSigner(address signer) external {
    require(signer != 0x0);

    // TODO: put the proposed proposed signer in a mapping, and only
    // move them to the signers list when all current signers have
    // ratified them
    if (signers[signer] == 0) {
      totalSigners++;
      signers[signer] = 1;
    }
  }

  function voteRemoveSigner(address signer) external {
    require(signer != 0x0);

    // TODO: put the proposed proposed signer in a mapping, and only
    // move them to the signers list when all current signers have
    // agreed to remove them

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
