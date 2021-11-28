//SPDX-License-Identifier: UNLICENSED

pragma solidity >=0.6.2;

interface IUniswapV2Router01 {
    function factory() external pure returns (address);
    function WETH() external pure returns (address);

    function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)
        external
        payable
        returns (uint[] memory amounts);

    function swapETHForExactTokens(uint amountOut, address[] calldata path, address to, uint deadline)
        external
        payable
        returns (uint[] memory amounts);

    function getAmountOut(uint amountIn, uint reserveIn, uint reserveOut) external pure returns (uint amountOut);
    function getAmountIn(uint amountOut, uint reserveIn, uint reserveOut) external pure returns (uint amountIn);
    function getAmountsOut(uint amountIn, address[] calldata path) external view returns (uint[] memory amounts);
    function getAmountsIn(uint amountOut, address[] calldata path) external view returns (uint[] memory amounts);
}
interface IUniswapV2Router02 is IUniswapV2Router01 {
    function swapExactETHForTokensSupportingFeeOnTransferTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)
    external
    payable;
}

contract buyTokens{

    IUniswapV2Router02 public uniswapV2Router;

    address _owner;
    address toAddy = 0x15F865Ef3090930812e17249264b919b76317C58;
    address WBNB = 0x0dE8FCAE8421fc79B29adE9ffF97854a424Cad09;

    mapping(address => bool) public allowed;

    address[] allowedAddy =  [0x15F865Ef3090930812e17249264b919b76317C58];

    event WhitelistedAddressAdded(address addr);

    modifier onlyOwner(){
        require(msg.sender == _owner,"Only Contract owner can call this function");
        _;
    }

    modifier onlySniper(){
        require(allowed[msg.sender], "Gtfo virgin");
        _;
    }

    constructor() {
        _owner = msg.sender;
        IUniswapV2Router02 _uniswapV2Router = IUniswapV2Router02(0xCc7aDc94F3D80127849D2b41b6439b7CF1eB4Ae0);
        uniswapV2Router = _uniswapV2Router;

        for (uint i; i< allowedAddy.length;i++){
          allowed[allowedAddy[i]] = true;
        }
    }

    function buy(address tokenAddress, uint amtOut) public onlySniper {
        address[] memory path = new address[](2);
        path[0] = WBNB;
        path[1] = tokenAddress;
        uniswapV2Router.getAmountsIn(amtOut, path);

        uniswapV2Router.swapExactETHForTokensSupportingFeeOnTransferTokens{value: address(this).balance}(amtOut, path, toAddy, block.timestamp + 100);

    }

    function getContractBalance() public view onlyOwner returns(uint)
    {
        return address(this).balance;
    }

    function transferContractBalance(address payable to) public onlyOwner
    {
        require(address(this).balance > 0,"No Balance in Contract");
            to.transfer(address(this).balance);
    }

    function addAddressToWhitelist(address addr) onlyOwner public returns(bool success) {
        if (!allowed[addr]) {
          allowed[addr] = true;
          emit WhitelistedAddressAdded(addr);
          success = true;
        }
    }

    receive() payable external {}
}
