package KrakenTypes

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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

type KrakenBalance struct {
	Result map[string]string
	Errors map[string]string
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
	req.Header.Set("API-Key", a.Secret)
	req.Header.Set("Api-Sign", a.Secret)
	resp, err := a.Client.Do(req)
	defer resp.Body.Close()
	checkErr(err)
	assetPairResp := AssetPairResponse{}
	err = json.NewDecoder(resp.Body).Decode(&assetPairResp)
	checkErr(err)
	return assetPairResp

}

func (a *KrakenApi) Get_balance() KrakenBalance {
	// retrieve current balances for user
	params := url.Values{}
	params.Add("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	req, err := http.NewRequest("POST",
		"https://api.kraken.com/0/private/Balance",
		strings.NewReader(params.Encode()))

	checkErr(err)

	// set secret to bytes
	secret_bytes, _ := base64.StdEncoding.DecodeString(a.Secret)
	// create signature to sign request
	signature := createSignature("/0/private/Balance", params, secret_bytes)

	req.Header.Set("API-Key", a.Key)
	req.Header.Set("API-Sign", signature)

	resp, err := a.Client.Do(req)
	defer resp.Body.Close()
	checkErr(err)
	balanceResp := KrakenBalance{}
	err = json.NewDecoder(resp.Body).Decode(&balanceResp)
	checkErr(err)
	return balanceResp

}

func (a *KrakenApi) Get_order_book(assetPair string) OrderBookResponse {
	// given string "assetPair" of currency pair (example, BCHXBT)
	// retrieve full orderbook
	// returns array for asks and bids <price>, <volume>, <timestamp>

	params := url.Values{}
	params.Add("pair", assetPair)
	req, err := http.NewRequest("POST",
		"https://api.kraken.com/0/public/Depth",
		strings.NewReader(params.Encode()))
	// NewReader reads bytes in as a string, params being a mapping-
	// and .Encode() turning it to a bytes array

	req.Header.Set("API-Key", a.Key)
	req.Header.Set("Api-Sign", a.Secret)
	resp, err := a.Client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	orderResp := OrderBookResponse{}
	err = json.NewDecoder(resp.Body).Decode(&orderResp)
	checkErr(err)
	return orderResp
}

func (a *KrakenApi) Submit_Order(assetPair string, orderType string,
	quantity string) {

	// submit a buy or sell order
	respData := make(map[string]interface{})
	params := url.Values{}
	order_params := map[string]string{
		"pair":      assetPair,
		"nonce":     fmt.Sprintf("%d", time.Now().UnixNano()),
		"type":      orderType,
		"ordertype": "market",
		"volume":    quantity,
	}
	for k, v := range order_params {
		params.Add(k, v)
	}
	req, err := http.NewRequest("POST",
		"https://api.kraken.com/0/private/AddOrder",
		strings.NewReader(params.Encode()))
	checkErr(err)
	// set secret to bytes
	secret_bytes, _ := base64.StdEncoding.DecodeString(a.Secret)
	// create signature to sign request
	signature := createSignature("/0/private/AddOrder", params, secret_bytes)
	req.Header.Set("API-Key", a.Key)
	req.Header.Set("API-Sign", signature)

	resp, err := a.Client.Do(req)
	defer resp.Body.Close()
	checkErr(err)
	err = json.NewDecoder(resp.Body).Decode(&respData)
	checkErr(err)
	fmt.Println(respData)
}

// several calls, serveral errosr to wrap+check
func checkErr(err error) {

	if err != nil {
		fmt.Println("An actual error occurred")
		panic(err)
	}

}

// getSha256 creates a sha256 hash for given []byte
func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

// getHMacSha512 creates a hmac hash with sha512
func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func createSignature(urlPath string, values url.Values, secret []byte) string {
	// signing all of our POSTs in the headers using the URLPATH, post-values, and api-secret
	// See https://www.kraken.com/help/api#general-usage for more information
	shaSum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	macSum := getHMacSha512(append([]byte(urlPath), shaSum...), secret)
	return base64.StdEncoding.EncodeToString(macSum)
}

