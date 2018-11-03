pragma solidity ^0.4.24;

import './zeppelin/token/StandardToken.sol';


contract SimpleToken is StandardToken {
    string  public  constant name = "SimpleToken";
    string  public  constant symbol = "ST";
    uint8   public  constant decimals = 9;
    uint256 public INITIAL_SUPPLY = 10000000000 * (10 ** uint256(decimals));

    //uint256 public INITIAL_SUPPLY = 10000000000;  //  * (10 ** uint256(decimals));

    constructor() public {
        totalSupply_ = INITIAL_SUPPLY;
        balances[msg.sender] = INITIAL_SUPPLY;
        emit Transfer(address(0x0), msg.sender, INITIAL_SUPPLY);
    }
}
