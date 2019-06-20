// GoKraken project GoKraken.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

type KrakenApi struct {
	key    string
	secret string
	client *http.Client
}

// orderbookresponse hard-coded
type OrderBookResponse struct {
	Result map[string]CurrencyOrderbook
	Errors map[string][]string
}

type CurrencyOrderbook struct {
	Asks [][]interface{} `json:"asks"`
	Bids [][]interface{} `json:"bids"`
}

type AssetPairResponse struct {
	Result map[string]interface{}
	Errors map[string][]string
}

func (a *KrakenApi) get_time() map[string]interface{} {

	var timeResp map[string]interface{}
	// create a new request
	req, err := http.NewRequest("POST", "https://api.kraken.com/0/public/Time", nil)
	checkErr(err)
	req.Header.Set("Api-key", a.key)
	req.Header.Set("Api-Sign", a.secret)
	resp, err := a.client.Do(req)
	checkErr(err)
	json.NewDecoder(resp.Body).Decode(&timeResp)
	return timeResp
}

func (a *KrakenApi) get_asset_pairs() AssetPairResponse {

	// get list of asset-pairs and info about them (fees, etc)

	req, err := http.NewRequest("POST", "https://api.kraken.com/0/public/AssetPairs", nil)
	checkErr(err)
	req.Header.Set("Api-key", a.key)
	req.Header.Set("Api-Sign", a.secret)
	resp, err := a.client.Do(req)
	defer resp.Body.Close()
	checkErr(err)
	assetPairResp := AssetPairResponse{}
	err = json.NewDecoder(resp.Body).Decode(&assetPairResp)
	checkErr(err)
	return assetPairResp

}

func (a *KrakenApi) get_order_book(assetPair string) OrderBookResponse {
	// given string "assetPair" of currency pair (example, BCHXBT)
	// retrieve full orderbook
	// returns array for asks and bids <price>, <volume>, <timestamp>

	//var data map[string]interface{}
	params := url.Values{}
	params.Add("pair", assetPair)
	req, err := http.NewRequest("POST",
		"https://api.kraken.com/0/public/Depth",
		strings.NewReader(params.Encode()))
	// NewReader reads bytes in as a string, params being a mapping-
	// and .Encode() turning it to a bytes array

	req.Header.Set("Api-key", a.key)
	req.Header.Set("Api-Sign", a.secret)
	resp, err := a.client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	orderResp := OrderBookResponse{}
	err = json.NewDecoder(resp.Body).Decode(&orderResp)
	checkErr(err)
	return orderResp
}

func main() {

	// POC demo's - all of this will end up in a KrakenAPI Package

	// create our Kraken Object
	api := KrakenApi{
		key:    "YOUR KEY HERE",
		secret: "YOUR SECRET KEY HERE",
		client: &http.Client{},
	}

	krakenTime := api.get_time()
	fmt.Println("Time according to Kraken:", krakenTime)
	testKeys := getKeys(krakenTime)
	for i := range testKeys {
		if krakenTime[testKeys[i]] != nil {
			fmt.Println(testKeys[i], ": ", krakenTime[testKeys[i]])
		}
	}

	// test getting valid asset pairs from Kraken
	assertPairResp := api.get_asset_pairs()
	assetPairsKeys := getKeys(assertPairResp.Result)
	fmt.Println("Valid Asset Pairs are:\n", assetPairsKeys)

	// test getting full orderbook for a currency pair
	orderBookTest := api.get_order_book("BCHXBT")
	fmt.Println("Price:", orderBookTest.Result["BCHXBT"].Asks[0][0])
	fmt.Println("Volume:", orderBookTest.Result["BCHXBT"].Asks[0][1])

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
