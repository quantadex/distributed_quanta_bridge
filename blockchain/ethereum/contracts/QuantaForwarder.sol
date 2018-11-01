pragma solidity ^0.4.24;

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

  function flush() public {
    emit LogFlushed(msg.sender, address(this).balance);
    destinationAddress.transfer(address(this).balance);
  }
}