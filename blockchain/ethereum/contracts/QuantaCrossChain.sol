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
  // all the ratified signers in an unsorted array
  address[] public signers;
  uint public requiredVotes = 0;

  struct Poll {
    address addr;
    bool[] votes;
  }
  uint8 public MAX_CANDIDATES = 10;

  Poll[] additions;
  uint8 additionsHead;

  Poll[] removals;
  uint8 removalsHead;

  /** the last successfully used txId */
  uint64 public txIdLast = 0;

  event TransactionResult(bool success,
                          uint64 txId,
                          address erc20Addr,
                          address toAddr,
                          uint256 amount,
                          bool[] verified);

  // https://ethereum.stackexchange.com/questions/42995/how-to-send-ether-to-a-contract-in-truffle-test/43011
  // @notice Logs the address of the sender and amounts paid to the contract
  event Fund(address indexed _from, uint _value);

  /**
   * Assigns the initial list of signers.
   * Must be called once after contract instantiation.
   * All other methods will fail or revert until this is called.
   * May only be called once by the owner. Further calls are ignored.
  */
  function assignInitialSigners(address[] initialSigners) external onlyOwner {
    assert(signers.length == 0);  // penalize (using assertions) if tried twice
    require(initialSigners.length > 0);

    // TODO: should we temporarily set totalSigners to MAX_UINT256?
    // TODO: dupe check?
    for(uint i=0; i<initialSigners.length; i++) {
      signers.push(initialSigners[i]);
    }

    updateRequiredVotes();
  }


  function isSigner(address signer) external view returns(bool) {
    for (uint i=0; i < signers.length; i++) {
      if (signers[i] == signer) {
        return true;
      }
    }
    return false;
  }


  function numSigners() external view returns (uint) {
    return signers.length;
  }


  function getAddCandidateVotes(address candidate) external view returns(uint count) {
    for(uint8 i=0; i<additions.length; i++) {
      if (additions[i].addr == candidate) {
        bool[] storage votes = additions[i].votes;
        for(i=0; i<signers.length; i++) {
          if (votes[i]) {
            count++;
          }
        }
        break;
      }
    }
  }


  function numAddCandidates() external view returns (uint count) {
    return additions.length;
  }


  function getRemoveCandidateVotes(address candidate) external view returns(uint count) {
    for(uint8 i=0; i<removals.length; i++) {
      if (removals[i].addr == candidate) {
        bool[] storage votes = removals[i].votes;
        for(i=0; i<signers.length; i++) {
          if (votes[i]==true) {
            count++;
          }
        }
        break;
      }
    }
  }


  function numRemoveCandidates() external view returns(uint) {
    return removals.length;
  }


  function paymentTx(uint64 txId,
                     address erc20Addr,
                     address toAddr,
                     uint256 amount,
                     uint8[] v, bytes32[] r, bytes32[] s) external {
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


  function voteAddSigner(address candidate) external returns(uint votesNeeded) {
    require(candidate != 0x0);
    assert(signers.length > 0);
    require(candidate != msg.sender);

    // find existing candidate, if any
    uint8 idx = 0;

    for(uint8 i=0; i<additions.length; i++) {
      if (additions[i].addr == candidate) {
        idx = i;
        break;
      }
    }

    if (i == additions.length) {  // did not find it
      // need to add it
      if (additions.length == MAX_CANDIDATES) {
        // replace at the head

        // initialize
        Poll storage spoll = additions[additionsHead];
        spoll.addr = candidate;
        bool[] storage cvotes = spoll.votes;
        for (uint j=0; j < signers.length; j++) {
          cvotes[j] = false;
        }

        idx = additionsHead;

        // increment the head
        additionsHead++;
        if (additionsHead == MAX_CANDIDATES) {
          additionsHead = 0;
        }
      } else {
        // not filled yet, simple case
        Poll memory mpoll = Poll(candidate, new bool[](signers.length));
        additions.push(mpoll);
        idx = uint8(additions.length - 1);
      }
    }

    // integrated vote tallying
    votesNeeded = requiredVotes;
    bool[] storage votes = additions[idx].votes;
    for(i=0; i<signers.length; i++) {
      // sender must be one of the signers
      if (votes[i]) {
        votesNeeded--;
      } else if (signers[i]==msg.sender) {
        votes[i] = true;
        votesNeeded--;
      }

      if (votesNeeded == 0) {
        _promoteCandidateToSigner(candidate);
        break;
      }
    }
  }


  function voteRemoveSigner(address candidate) external returns(uint votesNeeded) {
    require(candidate != 0x0);
    assert(signers.length > 0);
    require(candidate != msg.sender);

    // find existing candidate, if any
    uint8 idx = 0;

    for(uint8 i=0; i<removals.length; i++) {
      if (removals[i].addr == candidate) {
        idx = i;
        break;
      }
    }

    if (i == removals.length) {  // did not find it
      // need to add it
      if (removals.length == MAX_CANDIDATES) {
        // replace at the head

        // overwrite and initialize
        Poll storage spoll = removals[removalsHead];
        spoll.addr = candidate;
        votes = spoll.votes;
        for (uint j=0; j < signers.length; j++) {
          votes[j] = false;
        }

        idx = removalsHead;

        // increment the head
        removalsHead++;
        if (removalsHead == MAX_CANDIDATES) {
          removalsHead = 0;
        }
      } else {
        // not filled yet, simple case
        Poll memory mpoll = Poll(candidate, new bool[](signers.length));
        removals.push(mpoll);
        idx = uint8(removals.length - 1);
      }
    }

    // integrated vote tallying
    votesNeeded = requiredVotes;
    bool[] storage votes = removals[idx].votes;
    for(i=0; i<signers.length; i++) {
      // sender must be one of the signers
      if (votes[i]) {
        votesNeeded--;
      } else if (signers[i]==msg.sender) {
        votes[i] = true;
        votesNeeded--;
      }

      if (votesNeeded == 0) {
        _removeSigner(candidate);
        break;
      }
    }
  }


  // @notice Will receive any eth sent to the contract
  function () external payable {
    emit Fund(msg.sender, msg.value);
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

  function validateSignatures(bytes sigMsg,
                              uint numSigs,
                              uint8[] v,
                              bytes32[] r,
                              bytes32[] s,
                              bool[] outVerified) internal view returns (bool) {
    // with 8 signatures, Gas = 107770

    uint abortThreshold = numSigs - signers.length + 1;
    bytes32 signed = ECTools.toEthereumSignedMessage(string(sigMsg));
    address[] memory validated = new address[](numSigs);
    uint numValidated = 0;
    address addr;

    uint j = 0;
    for(uint i=0; ((i<numSigs) && (numValidated<signers.length) && (abortThreshold>0)); i++) {
      addr = ECTools.recoverSignerVRS(signed, v[i], r[i], s[i]);

      // we anticipate the number of signers to be small, so just do a linear scan
      for(j=0; j<signers.length; j++) {
        if (signers[j] == addr) {
          break;
        }
      }

      if (j==signers.length) {
        abortThreshold--;
      } else {
        // dupe check: we anticipate the number of signers to be small, so just do a linear scan
        for(j=0; j<numValidated; j++) {
          if (validated[j] == addr) {
            addr = 0;
            break;
          }
        }

        if (addr == 0) {
          abortThreshold--;
        } else {
          outVerified[i] = true;
          validated[numValidated] = addr;
          numValidated++;
        }
      }
    }

    return numValidated == signers.length;
  }

  function updateRequiredVotes() internal {
    // solidity division always rounds down
    // simple majority
    requiredVotes = signers.length / 2 + 1;
  }


  function _removeSigner(address signer) internal {
    // swap the signer in each of the voting arrays
    uint endIdx = removals.length - 1;
    uint8 signerIdx;

    for(signerIdx=0; signerIdx < signers.length; signerIdx++) {
      if (signers[signerIdx] == signer) {
        break;
      }
    }

    for (uint i=0; i<=endIdx; i++) {
      // put our last candidate here
      removals[i].votes[signerIdx] = removals[i].votes[signers.length-1];
      removals[i].votes.length--;  // decrease the size

      if (removals[i].addr == signer) {
        if (i != endIdx) {
          // put the last one in this index, we will delete the tail eventually
          removals[i] = removals[endIdx];
        }
      }
    }

    removals.length--;  // delete the last one

    signers[signerIdx] = signers[signers.length - 1];
    signers.length--;

    // FIXME: the one who casts the winning vote, pays this gas!
    updateRequiredVotes();
  }

  function _promoteCandidateToSigner(address candidate) internal {
    // add our candidate to the end of the voting list for all other addCandidates

    uint8 lastIdx = uint8(additions.length - 1);

    for(uint8 i=0; i < lastIdx; i++) {
      if (additions[i].addr == candidate) { // our candidate
        // move the last index to here and remap, technically not a FIFO queue
        additions[i] = additions[lastIdx];
      }
      additions[i].votes.length++;  //make space for our new signer
    }

    additions.length--;       // delete our candidate
    signers.push(candidate);  // add our candidate

    // FIXME: the one who casts the winning vote, pays this gas!
    updateRequiredVotes();
  }
}
