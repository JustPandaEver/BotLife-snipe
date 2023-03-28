# Cryptocurrency Transaction Handler

A collection of cryptocurrency trading bots to be used on the Binance Smart Chain.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Contributing](#contributing)
- [License](#license)

## Features

### Antirug Bot

- Prevent the trader from being scammed

### Go Presale Bot

- Purchase initial coin drops before anyone else

### Mempool Bot

- Deprecated JavaScript implementation of a presale bot

### Presale Bot

- Another deprecated JavaScript implementation of a presale bot

### Telegram bot

- A bot that can purchase coins showing high interest in a channel on the messaging app Telegram

### Scraper

- A Go web scraper to gather information for the presale bots

## Technologies Used

Cryptocurrency Transaction Handler is built using the following technologies:

- Go: A programming language designed for building fast, efficient, and scalable software
- GoEthereum: A client implementation of the Ethereum protocol

## Contributing

Contributions to Cryptocurrency Transaction Handler are welcome! If you have any bug reports, feature requests, or other suggestions, please open an issue or submit a pull request.

---

We will soon have the most powerful bot on the **Binance Smart Chain**. Created in JavaScript.

## ToDo

- [x] Fix sendtx.js
- [x] Optimize
- [x] Comment buy.js
- [x] Comment sendTx.js
- [x] Rename the scripts
- [x] Create config file
- [x] Link contract.js with bot.js
    - [x] Add block delay option to bot.js
    - [x] Separate sell to a different script
    - [x] Add argument for minimum token output
    - [x] Add percent to sell argument
    - [x] Determine when to sell
    - [x] Add verbose
- [x] Add option for minimum amount out
- [x] Add selector for buy method
- [x] Add selector for fees
- [x] Find gas of liquidity add and match it
- [x] Buy with balance in contract
- [x] Test the bot
- [ ] Link to bloXroute
- [ ] Develop frontend
- [ ] Create script to fill wallets
- [ ] Rewrite in GoLang
- [ ] Develop limit buy/sell
- [ ] Develop bot to buy tokens shared in call channels
- [ ] Check Solidity contract
- [ ] Develop contract to launch
    - [ ] Automatically enable/disable auto LP based on ratio
    - [ ] Max wallet
    - [ ] Max transfer
    - [ ] Marketing wallet

# PreSale Bot ToDo

- [x] Add approve
- [x] Support DX
- [x] Support Uniswap
- [ ] Add claim
- [ ] Add sell
- [ ] Add pinkWhitelist
