// GoKraken project GoKraken.go
package main

import (
	"fmt"
	"net/http"

	"github.com/KrakenTypes"
)

const (
	API_URL        = "https://api.kraken.com"
	API_VERSION    = "0"
	API_USER_AGENT = "TestUser"
)

// List of valid public methods
var publicMethods = []string{
	"Time",
	"Assets",
	"AssetPairs",
	"Ticker",
	"OHLC",
	"Depth",
	"Trades",
	"Spread",
}

// List of valid private methods
var privateMethods = []string{
	"Balance",
	"TradeBalance",
	"OpenOrders",
	"ClosedOrders",
	"QueryOrders",
	"TradesHistory",
	"QueryTrades",
	"OpenPositions",
	"Ledgers",
	"QueryLedgers",
	"TradeVolume",
	"AddOrder",
	"CancelOrder",
	"DepositMethods",
	"DepositAddresses",
	"DepositStatus",
	"WithdrawInfo",
	"Withdraw",
	"WithdrawStatus",
	"WithdrawCancel",
}

// These represent the minimum order sizes for the respective coins
// Should be monitored through here: https://support.kraken.com/hc/en-us/articles/205893708-What-is-the-minimum-order-size-
const (
	MinimumREP  = 0.3
	MinimumXBT  = 0.002
	MinimumBCH  = 0.002
	MinimumDASH = 0.03
	MinimumDOGE = 3000.0
	MinimumEOS  = 3.0
	MinimumETH  = 0.02
	MinimumETC  = 0.3
	MinimumGNO  = 0.03
	MinimumICN  = 2.0
	MinimumLTC  = 0.1
	MinimumMLN  = 0.1
	MinimumXMR  = 0.1
	MinimumXRP  = 30.0
	MinimumXLM  = 300.0
	MinimumZEC  = 0.02
	MinimumUSDT = 5.0
)

func main() {

	// POC demo's - all of this will end up in a KrakenAPI Package

	// create our Kraken Object
	api := KrakenTypes.KrakenApi{
		Key:    "FAKE_KEY",
		Secret: "FAKE_SECRET",
		Client: &http.Client{},
	}

	krakenTime := api.Get_time()
	fmt.Println("Time according to Kraken:", krakenTime)
	testKeys := getKeys(krakenTime)
	for i := range testKeys {
		if krakenTime[testKeys[i]] != nil {
			fmt.Println(testKeys[i], ": ", krakenTime[testKeys[i]])
		}
	}

	// test getting valid asset pairs from Kraken
	assertPairResp := api.Get_asset_pairs()
	assetPairsKeys := getKeys(assertPairResp.Result)
	fmt.Println("Valid Asset Pairs are:\n", assetPairsKeys)

	// test getting full orderbook for a currency pair
	orderBookTest := api.Get_order_book("BCHXBT")
	fmt.Println("Debug - orderbooktest is", orderBookTest)
	fmt.Println("Price:", orderBookTest.Result["BCHXBT"].Asks[0][0])
	fmt.Println("Volume:", orderBookTest.Result["BCHXBT"].Asks[0][1])

	// test getting balance
	balance := api.Get_balance()
	fmt.Println("balance is", balance.Result)

	// buy BCH with XBT
	api.Submit_Order("BCHXBT", "buy", "0.05")

	// sell BCH for XBT
	api.Submit_Order("BCHXBT", "sell", "0.05")
}

// several calls, serveral errosr to wrap+check
func checkErr(err error) {

	if err != nil {
		fmt.Println("An actual error occurred")
		panic(err)
	}

}

// during development we may be unsure what keys we're dealilng with
// this function allows us to quickly debug printout slices of keys
func getKeys(mapping map[string]interface{}) []string {

	keys := make([]string, len(mapping))
	for k := range mapping {
		keys = append(keys, k)
	}
	return keys
}

