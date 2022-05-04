# real-time Token Transfer Monitor on Blockchains

This repo contains codes to monitor token Transfer events on blockchains.
Users can add intereted addresses to `redis`, then start the program to monitor their transfer events (or balance change).

## Suppoted chain networks:

1. Ethereum: ERC20 Tokens
2. Binance Smart Chain: ERC20 Tokens
3. Tron
4. Solana

## Basic idea

### Ethereum, BSC, Tron

Query latest mint block every certain seconds, then filter target token transfer events/logs.

### Solana

Use the officail websocket api to subscribe events omitted by official token program `TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA`, the use `memcmp` to filter the traget solana token transfer events.
