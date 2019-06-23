package KrakenTypes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type KrakenApi struct {
	Key    string
	Secret string
	Client *http.Client
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

func (a *KrakenApi) Get_time() map[string]interface{} {

	var timeResp map[string]interface{}
	// create a new request
	req, err := http.NewRequest("POST", "https://api.kraken.com/0/public/Time", nil)
	checkErr(err)
	req.Header.Set("Api-key", a.Key)
	req.Header.Set("Api-Sign", a.Secret)
	resp, err := a.Client.Do(req)
	checkErr(err)
	json.NewDecoder(resp.Body).Decode(&timeResp)
	return timeResp
}

func (a *KrakenApi) Get_asset_pairs() AssetPairResponse {

	// get list of asset-pairs and info about them (fees, etc)

	req, err := http.NewRequest("POST", "https://api.kraken.com/0/public/AssetPairs", nil)
	checkErr(err)
	req.Header.Set("Api-key", a.Key)
	req.Header.Set("Api-Sign", a.Secret)
	resp, err := a.Client.Do(req)
	defer resp.Body.Close()
	checkErr(err)
	assetPairResp := AssetPairResponse{}
	err = json.NewDecoder(resp.Body).Decode(&assetPairResp)
	checkErr(err)
	return assetPairResp

}

func (a *KrakenApi) Get_balance() {
	// retrieve current balances for user
	req, err := http.NewRequest("GET", "https://api.kraken.com/0/private/Balance", nil)
	checkErr(err)
	req.Header.Set("Api-key", a.Key)
	req.Header.Set("Api-Sign", a.Secret)
	resp, err := a.Client.Do(req)
	defer resp.Body.Close()
	checkErr(err)
	balanceResp := AssetPairResponse{}
	err = json.NewDecoder(resp.Body).Decode(&balanceResp)
	fmt.Println(balanceResp)

}

func (a *KrakenApi) Get_order_book(assetPair string) OrderBookResponse {
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

	req.Header.Set("Api-key", a.Key)
	req.Header.Set("Api-Sign", a.Secret)
	resp, err := a.Client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	orderResp := OrderBookResponse{}
	err = json.NewDecoder(resp.Body).Decode(&orderResp)
	checkErr(err)
	return orderResp
}

// several calls, serveral errosr to wrap+check
func checkErr(err error) {

	if err != nil {
		fmt.Println("An actual error occurred")
		panic(err)
	}

}
