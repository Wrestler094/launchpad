// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title Presale
 * @dev Contract for managing token presales
 */
contract Presale is Ownable, ReentrancyGuard {
    IERC20 public token;
    
    uint256 public rate; // tokens per wei
    uint256 public softCap;
    uint256 public hardCap;
    uint256 public deadline;
    uint256 public raised;
    uint256 public tokensSold;
    
    bool public presaleActive;
    bool public presaleFinalized;
    
    mapping(address => uint256) public contributions;
    mapping(address => uint256) public tokensPurchased;
    
    event TokensPurchased(address indexed buyer, uint256 amount, uint256 tokens);
    event PresaleFinalized(bool successful);
    event FundsWithdrawn(address indexed owner, uint256 amount);
    
    modifier presaleIsActive() {
        require(presaleActive, "Presale is not active");
        require(block.timestamp <= deadline, "Presale has ended");
        require(raised < hardCap, "Hard cap reached");
        _;
    }
    
    constructor(
        address _token,
        uint256 _rate,
        uint256 _softCap,
        uint256 _hardCap,
        uint256 _deadline,
        address _owner
    ) Ownable(_owner) {
        require(_token != address(0), "Invalid token address");
        require(_rate > 0, "Rate must be positive");
        require(_softCap < _hardCap, "Soft cap must be less than hard cap");
        require(_deadline > block.timestamp, "Deadline must be in the future");
        
        token = IERC20(_token);
        rate = _rate;
        softCap = _softCap;
        hardCap = _hardCap;
        deadline = _deadline;
        presaleActive = true;
    }
    
    /**
     * @dev Purchase tokens with ETH
     */
    function buyTokens() external payable presaleIsActive nonReentrant {
        require(msg.value > 0, "Must send ETH");
        require(raised + msg.value <= hardCap, "Would exceed hard cap");
        
        uint256 tokens = msg.value * rate;
        require(token.balanceOf(address(this)) >= tokens, "Not enough tokens in contract");
        
        contributions[msg.sender] += msg.value;
        tokensPurchased[msg.sender] += tokens;
        raised += msg.value;
        tokensSold += tokens;
        
        require(token.transfer(msg.sender, tokens), "Token transfer failed");
        
        emit TokensPurchased(msg.sender, msg.value, tokens);
    }
    
    /**
     * @dev Finalize the presale
     */
    function finalizePresale() external onlyOwner {
        require(!presaleFinalized, "Presale already finalized");
        require(block.timestamp > deadline || raised >= hardCap, "Presale still active");
        
        presaleActive = false;
        presaleFinalized = true;
        
        bool successful = raised >= softCap;
        emit PresaleFinalized(successful);
        
        if (successful) {
            // Transfer raised funds to owner
            payable(owner()).transfer(raised);
            emit FundsWithdrawn(owner(), raised);
        }
    }
    
    /**
     * @dev Withdraw funds (only if presale was successful)
     */
    function withdrawFunds() external onlyOwner {
        require(presaleFinalized, "Presale not finalized");
        require(raised >= softCap, "Soft cap not reached");
        
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        payable(owner()).transfer(balance);
        emit FundsWithdrawn(owner(), balance);
    }
    
    /**
     * @dev Get refund if presale failed
     */
    function getRefund() external nonReentrant {
        require(presaleFinalized, "Presale not finalized");
        require(raised < softCap, "Presale was successful");
        require(contributions[msg.sender] > 0, "No contribution found");
        
        uint256 contribution = contributions[msg.sender];
        contributions[msg.sender] = 0;
        
        payable(msg.sender).transfer(contribution);
    }
    
    /**
     * @dev Emergency withdraw remaining tokens
     */
    function withdrawRemainingTokens() external onlyOwner {
        require(presaleFinalized, "Presale not finalized");
        
        uint256 remainingTokens = token.balanceOf(address(this));
        if (remainingTokens > 0) {
            require(token.transfer(owner(), remainingTokens), "Token transfer failed");
        }
    }
    
    /**
     * @dev Get presale info
     */
    function getPresaleInfo() external view returns (
        uint256 _rate,
        uint256 _softCap,
        uint256 _hardCap,
        uint256 _deadline,
        uint256 _raised,
        uint256 _tokensSold,
        bool _active,
        bool _finalized
    ) {
        return (rate, softCap, hardCap, deadline, raised, tokensSold, presaleActive, presaleFinalized);
    }
}