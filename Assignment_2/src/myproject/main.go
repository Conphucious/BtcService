package main

import (
	loggly "github.com/jamespearly/loggly"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type BtcEntry struct {
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	Bpi        map[string]struct {
		Code        string  `json:"code"`
		Symbol      string  `json:"symbol"`
		Rate        string  `json:"rate"`
		Description string  `json:"description"`
		Rate_Float  float32 `json:"rate_float"`
	} `json:"bpi"`
}

func main() {
	for {
		go work()
		time.Sleep(10 * time.Minute)
	}
}

func work() {
	apiJson := getResponse("https://api.coindesk.com/v1/bpi/currentprice.json")
	parseJson(apiJson)
	logglyConnection(apiJson)
}

func getResponse(url string) []byte {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return body
}

func parseJson(resp []byte) BtcEntry {
	var btc BtcEntry
	json.Unmarshal(resp, &btc)

	for _, data := range btc.Bpi {
		fmt.Println("Name: ", btc.ChartName,
			"\nDisclaimer: ", btc.Disclaimer,
			"\nCode: [", data.Code, "]",
			"\nSym: [", data.Symbol, "]",
			"\nRate: $", data.Rate,
			"\nRate in Float: $", data.Rate_Float,
			"\nDescription: ", data.Description,
			"\n\n")
	}

	//	fmt.Println("USD Rate: ", btc.Bpi["USD"].Rate)

	return btc
}

func logglyConnection(resp []byte) {
	// I don't care that my token is here
	os.Setenv("LOGGLY_TOKEN", "8224b513-31a1-437c-924a-5b0e420f55ec")

	// new Client
	client := loggly.New("GoLoggly")

	// Valid EchoSend (message echoed to console and no error returned)
	//err := client.EchoSend("info", string(resp))
	//fmt.Println("err:", err)

	err := client.Send("info", string(resp))

	byteValue := len(resp)
	fmt.Println("Succcessfully sent JSON data of", byteValue, "bytes at", time.Now(), ".\n")

	if err != nil {
		log.Fatal(err)
	}
}
