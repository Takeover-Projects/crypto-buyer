package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetKeys(t *testing.T) {
	var testMap = make(map[string]interface{})
	testMap["testKey"] = "testValue"
	resp := getKeys(testMap)

	if resp[1] != "testKey" {
		t.Error("Incorrect Key returned", resp[1], "instead of", "testKey")
	}
}

func Get_Time_Server(req *http.Request) (*http.Response, error) {

	// mock for POST of Get_Time_Server used with httpmock

	fmt.Println("running me!")
	return httpmock.NewJsonResponse(200, map[string]interface{}{
		"result": map[string]string{
			"rfc1123": "Tue, 18 Jun 19 02:26:07 +0000 unixtime:1.560824767e+09",
		},
		"error": nil,
	})

}
  
func TestGetTime(t *testing.T) {

	// retrieving unix epoch timestamp from Kraken API
	// serves as a very easy example how to drive test-cases for objects using-
	// http clients

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact url match
	httpmock.RegisterResponder("POST", "https://api.kraken.com/0/public/Time",
		Get_Time_Server)

	api := KrakenApi{
		key:    "testkey1",
		secret: "testkey2",
		client: &http.Client{},
	}

	resp := api.get_time()
	if resp["result"].(map[string]interface{})["rfc1123"] != "Tue, 18 Jun 19 02:26:07 +0000 unixtime:1.560824767e+09" {
		t.Error("GetTime returned invalid response")
	}

}
