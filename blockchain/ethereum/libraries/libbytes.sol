pragma solidity ^0.4.24;

/** byte utilities */
library libbytes {
  // @credit https://ethereum.stackexchange.com/questions/884/how-to-convert-an-address-to-bytes-in-solidity
  function addressToBytes(address a) internal pure returns (bytes b) {
     assembly {
        let m := mload(0x40)
        mstore(add(m, 20), xor( 0x140000000000000000000000000000000000000000, a))
        mstore(0x40, add(m, 52))
        b := m
     }
  }
}
