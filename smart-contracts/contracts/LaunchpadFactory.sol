// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./MyToken.sol";
import "./Presale.sol";

/**
 * @title LaunchpadFactory
 * @dev Factory contract for deploying tokens and presales
 */
contract LaunchpadFactory {
    struct TokenInfo {
        address tokenAddress;
        string name;
        string symbol;
        uint256 totalSupply;
        address creator;
        uint256 createdAt;
    }
    
    struct PresaleInfo {
        address presaleAddress;
        address tokenAddress;
        uint256 rate;
        uint256 softCap;
        uint256 hardCap;
        uint256 deadline;
        address creator;
        uint256 createdAt;
    }
    
    TokenInfo[] public tokens;
    PresaleInfo[] public presales;
    
    mapping(address => uint256[]) public userTokens;
    mapping(address => uint256[]) public userPresales;
    mapping(address => uint256) public tokenToIndex;
    mapping(address => uint256) public presaleToIndex;
    
    event TokenCreated(
        address indexed tokenAddress,
        string name,
        string symbol,
        uint256 totalSupply,
        address indexed creator
    );
    
    event PresaleCreated(
        address indexed presaleAddress,
        address indexed tokenAddress,
        uint256 rate,
        uint256 softCap,
        uint256 hardCap,
        uint256 deadline,
        address indexed creator
    );
    
    /**
     * @dev Create a new ERC20 token
     */
    function createToken(
        string memory name,
        string memory symbol,
        uint256 totalSupply
    ) external returns (address) {
        MyToken newToken = new MyToken(name, symbol, totalSupply, msg.sender);
        address tokenAddress = address(newToken);
        
        TokenInfo memory tokenInfo = TokenInfo({
            tokenAddress: tokenAddress,
            name: name,
            symbol: symbol,
            totalSupply: totalSupply,
            creator: msg.sender,
            createdAt: block.timestamp
        });
        
        tokens.push(tokenInfo);
        uint256 tokenIndex = tokens.length - 1;
        userTokens[msg.sender].push(tokenIndex);
        tokenToIndex[tokenAddress] = tokenIndex;
        
        emit TokenCreated(tokenAddress, name, symbol, totalSupply, msg.sender);
        
        return tokenAddress;
    }
    
    /**
     * @dev Create a new presale
     */
    function createPresale(
        address tokenAddress,
        uint256 rate,
        uint256 softCap,
        uint256 hardCap,
        uint256 deadline
    ) external returns (address) {
        require(tokenAddress != address(0), "Invalid token address");
        
        Presale newPresale = new Presale(
            tokenAddress,
            rate,
            softCap,
            hardCap,
            deadline,
            msg.sender
        );
        
        address presaleAddress = address(newPresale);
        
        PresaleInfo memory presaleInfo = PresaleInfo({
            presaleAddress: presaleAddress,
            tokenAddress: tokenAddress,
            rate: rate,
            softCap: softCap,
            hardCap: hardCap,
            deadline: deadline,
            creator: msg.sender,
            createdAt: block.timestamp
        });
        
        presales.push(presaleInfo);
        uint256 presaleIndex = presales.length - 1;
        userPresales[msg.sender].push(presaleIndex);
        presaleToIndex[presaleAddress] = presaleIndex;
        
        emit PresaleCreated(presaleAddress, tokenAddress, rate, softCap, hardCap, deadline, msg.sender);
        
        return presaleAddress;
    }
    
    /**
     * @dev Get total number of tokens created
     */
    function getTokenCount() external view returns (uint256) {
        return tokens.length;
    }
    
    /**
     * @dev Get total number of presales created
     */
    function getPresaleCount() external view returns (uint256) {
        return presales.length;
    }
    
    /**
     * @dev Get tokens created by a user
     */
    function getUserTokens(address user) external view returns (uint256[] memory) {
        return userTokens[user];
    }
    
    /**
     * @dev Get presales created by a user
     */
    function getUserPresales(address user) external view returns (uint256[] memory) {
        return userPresales[user];
    }
    
    /**
     * @dev Get token info by index
     */
    function getTokenInfo(uint256 index) external view returns (TokenInfo memory) {
        require(index < tokens.length, "Token index out of bounds");
        return tokens[index];
    }
    
    /**
     * @dev Get presale info by index
     */
    function getPresaleInfo(uint256 index) external view returns (PresaleInfo memory) {
        require(index < presales.length, "Presale index out of bounds");
        return presales[index];
    }
}