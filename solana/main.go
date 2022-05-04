package main

import (
	"solana/lib"
)

func main() {
	// TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA: Solana USDT address
	um_confirmed := lib.TokenTransferMonitor{COMMITMENT: "confirmed", MonitorTokenAddress: "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"}
	um_confirmed.Monitor()
}
