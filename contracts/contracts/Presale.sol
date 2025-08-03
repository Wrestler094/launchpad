// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Presale is ReentrancyGuard, Ownable {
    IERC20 public token;
    uint256 public rate; // tokens per ETH
    uint256 public hardCap; // maximum ETH to raise
    uint256 public softCap; // minimum ETH to raise
    uint256 public deadline;
    uint256 public totalRaised;
    bool public presaleEnded;
    bool public goalReached;
    
    mapping(address => uint256) public contributions;
    address[] public contributors;
    
    event TokensPurchased(address indexed buyer, uint256 amount, uint256 tokens);
    event PresaleFinalized(bool goalReached, uint256 totalRaised);
    event RefundClaimed(address indexed contributor, uint256 amount);
    
    constructor(
        address _token,
        uint256 _rate,
        uint256 _hardCap,
        uint256 _softCap,
        uint256 _durationInDays,
        address _owner
    ) Ownable(_owner) {
        require(_token != address(0), "Invalid token address");
        require(_rate > 0, "Rate must be greater than 0");
        require(_hardCap > _softCap, "Hard cap must be greater than soft cap");
        require(_durationInDays > 0, "Duration must be greater than 0");
        
        token = IERC20(_token);
        rate = _rate;
        hardCap = _hardCap;
        softCap = _softCap;
        deadline = block.timestamp + (_durationInDays * 1 days);
    }
    
    modifier presaleActive() {
        require(!presaleEnded && block.timestamp < deadline, "Presale not active");
        _;
    }
    
    function buyTokens() external payable nonReentrant presaleActive {
        require(msg.value > 0, "Must send ETH");
        require(totalRaised + msg.value <= hardCap, "Would exceed hard cap");
        
        uint256 tokenAmount = msg.value * rate;
        require(token.balanceOf(address(this)) >= tokenAmount, "Insufficient tokens");
        
        if (contributions[msg.sender] == 0) {
            contributors.push(msg.sender);
        }
        
        contributions[msg.sender] += msg.value;
        totalRaised += msg.value;
        
        token.transfer(msg.sender, tokenAmount);
        
        emit TokensPurchased(msg.sender, msg.value, tokenAmount);
        
        if (totalRaised >= hardCap) {
            _finalizePresale();
        }
    }
    
    function finalizePresale() external onlyOwner {
        require(block.timestamp >= deadline || totalRaised >= hardCap, "Cannot finalize yet");
        _finalizePresale();
    }
    
    function _finalizePresale() internal {
        presaleEnded = true;
        goalReached = totalRaised >= softCap;
        
        if (goalReached) {
            payable(owner()).transfer(address(this).balance);
        }
        
        emit PresaleFinalized(goalReached, totalRaised);
    }
    
    function claimRefund() external nonReentrant {
        require(presaleEnded, "Presale not ended");
        require(!goalReached, "Goal was reached, no refunds");
        require(contributions[msg.sender] > 0, "No contribution to refund");
        
        uint256 refundAmount = contributions[msg.sender];
        contributions[msg.sender] = 0;
        
        payable(msg.sender).transfer(refundAmount);
        
        emit RefundClaimed(msg.sender, refundAmount);
    }
    
    function withdrawUnsoldTokens() external onlyOwner {
        require(presaleEnded, "Presale not ended");
        uint256 unsoldTokens = token.balanceOf(address(this));
        if (unsoldTokens > 0) {
            token.transfer(owner(), unsoldTokens);
        }
    }
    
    function getContributorsCount() external view returns (uint256) {
        return contributors.length;
    }
    
    function isActive() external view returns (bool) {
        return !presaleEnded && block.timestamp < deadline && totalRaised < hardCap;
    }
}