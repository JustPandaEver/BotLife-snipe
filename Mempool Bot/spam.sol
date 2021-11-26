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
}

contract buyTokens{

    IUniswapV2Router02 public immutable uniswapV2Router;

    address _owner;
    address toAddy = 0x79708C882483A93C9F6f1D385a5301e1E08B0438;
    address WBNB = 0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c;

    mapping(address => bool) public allowed;

    address[] allowedAddy =  [0x79708C882483A93C9F6f1D385a5301e1E08B0438];

    event WhitelistedAddressAdded(address addr);

    modifier onlyOwner(){
        require(msg.sender == _owner,"Only Contract owner can call this function");
        _;
    }

    modifier onlySniper(){
        require(allowed[msg.sender], "Gtfo virgin");
        _;
    }

    constructor(){
        _owner = msg.sender;
        IUniswapV2Router02 _uniswapV2Router = IUniswapV2Router02(0x10ED43C718714eb63d5aA57B78B54704E256024E);
        uniswapV2Router = _uniswapV2Router;

        for (uint i; i< allowedAddy.length;i++){
          allowed[allowedAddy[i]] = true;
        }
    }

    function memBuy(address tokenAddress) public onlySniper {
        address[] memory path = new address[](2);
        path[0] = WBNB;
        path[1] = tokenAddress;

        uniswapV2Router.swapExactETHForTokens{value: address(this).balance}(1, path, toAddy, block.timestamp + 21);
    }

    function buy(address tokenAddress) public onlySniper {
        address[] memory path = new address[](2);
        path[0] = WBNB;
        path[1] = tokenAddress;
        uniswapV2Router.getAmountsOut(address(this).balance, path);

        uniswapV2Router.swapExactETHForTokens{value: address(this).balance}(1, path, toAddy, block.timestamp + 21);
    }

    function buy2(address tokenAddress, uint amtOut, int x) public onlySniper {
        address[] memory path = new address[](2);
        path[0] = WBNB;
        path[1] = tokenAddress;
        uniswapV2Router.getAmountsIn(amtOut, path);

        for (int i = 0; i < x; i++) {
          uniswapV2Router.swapETHForExactTokens{value: address(this).balance}(amtOut, path, toAddy, block.timestamp + 21);
        }

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
