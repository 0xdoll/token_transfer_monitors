package lib

import "strconv"

type TokenTx struct {
	Id           string      `json:"_id"`
	Address      string      `json:"address"`
	Signature    []string    `json:"signature"`
	ChangeType   string      `json:"changeType"`
	ChangeAmount interface{} `json:"changeAmount"`
	Decimals     int         `json:"decimals"`
	PostBalance  interface{} `json:"postBalance"`
	PreBalance   interface{} `json:"preBalance"`
	TokenAddress string      `json:"tokenAddress"`
	Symbol       string      `json:"symbol"`
	BlockTime    int         `json:"blockTime"`
	Slot         int         `json:"slot"`
	Fee          int         `json:"fee"`
	Owner        string      `json:"owner"`
	Balance      struct {
		Amount   string `json:"amount"`
		Decimals int    `json:"decimals"`
	} `json:"balance"`
}

type TokenTxsResp struct {
	Total int       `json:"total"`
	Data  []TokenTx `json:"data"`
}

func ToInt64(v interface{}) int64 {
	var iv int64
	switch v := v.(type) {
	case float64:
		iv = int64(v)
	case string:
		iv, _ = strconv.ParseInt(v, 10, 64)
	}
	return iv
}

type ProgramSubResp struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Id      int         `json:"id"`

	Method string `json:"method"`
	Params struct {
		Result struct {
			Context struct {
				Slot int `json:"slot"`
			} `json:"context"`
			Value struct {
				Pubkey  string `json:"pubkey"`
				Account struct {
					Lamports int `json:"lamports"`
					Data     struct {
						Program string `json:"program"`
						Parsed  struct {
							Info struct {
								Delegate        string `json:"delegate"`
								DelegatedAmount struct {
									Amount         string  `json:"amount"`
									Decimals       int     `json:"decimals"`
									UiAmount       float64 `json:"uiAmount"`
									UiAmountString string  `json:"uiAmountString"`
								} `json:"delegatedAmount"`
								IsNative    bool   `json:"isNative"`
								Mint        string `json:"mint"`
								Owner       string `json:"owner"`
								State       string `json:"state"`
								TokenAmount struct {
									Amount         string  `json:"amount"`
									Decimals       int     `json:"decimals"`
									UiAmount       float64 `json:"uiAmount"`
									UiAmountString string  `json:"uiAmountString"`
								} `json:"tokenAmount"`
							} `json:"info"`
							Type string `json:"type"`
						} `json:"parsed"`
						Space int `json:"space"`
					} `json:"data"`
					Owner      string `json:"owner"`
					Executable bool   `json:"executable"`
					RentEpoch  int    `json:"rentEpoch"`
				} `json:"account"`
			} `json:"value"`
		} `json:"result"`
		Subscription int `json:"subscription"`
	} `json:"params"`
}
