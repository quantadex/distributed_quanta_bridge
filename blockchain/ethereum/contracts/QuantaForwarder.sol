pragma solidity ^0.4.24;
import { ERC20Basic } from "./zeppelin/token/ERC20Basic.sol";

contract QuantaForwarder {

  address public destinationAddress;
  string public quantaAddress;

  event LogForwarded(address indexed sender, uint amount);
  event LogFlushed(address indexed sender, uint amount);
  event LogCreated(address trust, string quanta);

  constructor(address trust, string quantaAddr) public {
    destinationAddress = trust;
    quantaAddress = quantaAddr;
    emit LogCreated(trust,quantaAddr);
  }

  function() payable public {
    emit LogForwarded(msg.sender, msg.value);
    destinationAddress.transfer(msg.value);
  }

  /**
   * Execute a token transfer of the full balance from the forwarder token to the parent address
   * @param tokenContractAddress the address of the erc20 token contract
   */
  function flushTokens(address tokenContractAddress) public {
    ERC20Basic instance = ERC20Basic(tokenContractAddress);
    address forwarderAddress = address(this);
    uint256 forwarderBalance = instance.balanceOf(forwarderAddress);
    if (forwarderBalance == 0) {
      return;
    }
    if (!instance.transfer(destinationAddress, forwarderBalance)) {
      revert();
    }
  }

  function flush() public {
    emit LogFlushed(msg.sender, address(this).balance);
    destinationAddress.transfer(address(this).balance);
  }
}