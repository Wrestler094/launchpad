// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MyToken
 * @dev ERC20 token that can be deployed by users through the Launchpad platform
 */
contract MyToken is ERC20, Ownable {
    uint256 private _totalSupply;
    
    constructor(
        string memory name,
        string memory symbol,
        uint256 totalSupply,
        address owner
    ) ERC20(name, symbol) Ownable(owner) {
        _totalSupply = totalSupply * 10**decimals();
        _mint(owner, _totalSupply);
    }
    
    /**
     * @dev Returns the total supply of tokens
     */
    function totalSupply() public view override returns (uint256) {
        return _totalSupply;
    }
    
    /**
     * @dev Mint additional tokens (only owner)
     */
    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
        _totalSupply += amount;
    }
    
    /**
     * @dev Burn tokens (only owner)
     */
    function burn(uint256 amount) public onlyOwner {
        _burn(msg.sender, amount);
        _totalSupply -= amount;
    }
}