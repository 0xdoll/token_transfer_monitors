package lib

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/imroc/req/v3"
)

var (
	rdb              *redis.Client
	reqClient        *req.Client
	ctx              context.Context
	SOLANA_ADDRS_KEY string = "SOLANA_USDT_ADDR"
	REDIS_URI               = os.Getenv("REDIS_URI")
)

func init() {
	if len(REDIS_URI) == 0 {
		REDIS_URI = "localhost:6379"
	}
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     REDIS_URI,
		Password: "",
		DB:       0,
	})
	reqClient = req.C().SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36")
}

type TokenTransferMonitor struct {
	Conn                *websocket.Conn
	SubScriptionId      int64
	COMMITMENT          string
	MonitorTokenAddress string
}

func (um *TokenTransferMonitor) Subscribe() {
	if um.COMMITMENT == "" {
		um.COMMITMENT = "confirmed"
	}
	var err error
	for i := 0; i < 10; i++ {
		um.Conn, _, err = websocket.DefaultDialer.Dial("ws://api.mainnet-beta.solana.com", nil)
		if err != nil {
			fmt.Println(um.COMMITMENT, " ", i, " times dial err: ", err)
			continue
		}
		break
	}
	err = um.Conn.WriteMessage(1, []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"programSubscribe","params":["TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",{"commitment":"%s","encoding":"jsonParsed","filters":[{"memcmp":{"offset":0,"bytes":"%s"}}]}]}`, um.COMMITMENT, um.MonitorTokenAddress)))

	if err != nil {
		fmt.Println("\n", um.COMMITMENT, " write err:", err)
	}
}

func (um *TokenTransferMonitor) Monitor() {
	if um.Conn == nil {
		um.Subscribe()
	}
	defer func(Conn *websocket.Conn, subScriptionId int64) {
		err := Conn.WriteMessage(1, []byte(fmt.Sprintf(`{"jsonrpc":"2.0", "id":1, "method":"programUnsubscribe", "params":[%d]}`, subScriptionId)))
		err = Conn.Close()
		if err != nil {
			fmt.Println(um.COMMITMENT, " um.Conn.Close error: ", err.Error())
		}
		fmt.Println(um.COMMITMENT, " um.Conn.Close Successfully: ")
	}(um.Conn, um.SubScriptionId)
	for {
		psr := ProgramSubResp{}
		err := um.Conn.ReadJSON(&psr)
		if err != nil {
			fmt.Println("\n", um.COMMITMENT, " read err:", err, " now re-subscribe...")
			um.Subscribe()
			continue
		}
		if psr.Result != nil {
			switch result := psr.Result.(type) {
			case float64:
				um.SubScriptionId = int64(result)
				fmt.Println(um.COMMITMENT, " Renew SubScriptionId: ", um.SubScriptionId)
				continue
			}
		}
		ownerAccount := psr.Params.Result.Value.Account.Data.Parsed.Info.Owner
		if isIn, _ := rdb.SIsMember(ctx, SOLANA_ADDRS_KEY, ownerAccount).Result(); isIn {
			fmt.Println("\n", um.COMMITMENT, "(", time.Now().Format("01-02 15:04:05.000"), ")", psr)
		}
	}
}

// https://public-api.solscan.io/account/splTransfers?account=AAVBgjLnAAHc4RCrjU3A93TwvAyYaA151YF5jTZkSmHG&fromTime=1647707503&toTime=1647709445&offset=0&limit=50
