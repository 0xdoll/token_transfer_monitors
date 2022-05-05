package main

import (
	"solana/lib"
)

func main() {
	// Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB // Solana USDT address
	um_confirmed := lib.TokenTransferMonitor{COMMITMENT: "confirmed", MonitorTokenAddress: "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"}
	um_confirmed.Monitor()
}
