package main

import (
	"encoding/json"
	"fmt"
	loggly "github.com/jamespearly/loggly"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type BtcEntry struct {
	//Time       time.Time `json:"time"` // improt "time"
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	Bpis       map[string]struct {
		Code        string  `json:"code"`
		Symbol      string  `json:"symbol"`
		Rate        string  `json:"rate"`
		Description string  `json:"description"`
		Rate_Float  float32 `json:"rate_float"`
	} `json:"bpi"`
}

func main() {
	//apiJson := getResponse("https://api.coindesk.com/v1/bpi/currentprice.json")
	//parseJson(apiJson)
	logglyConnection()

}

func logglyConnection() {
	// new Client
	client := loggly.New("BTC482")

	// Valid EchoSend (message echoed to console and no error returned)
	err := client.EchoSend("info", "Good morning!")
	fmt.Println("err:", err)

	// Valid Send (no error returned)
	err = client.Send("error", "Good morning! No echo.")
	fmt.Println("err:", err)

	// Invalid EchoSend -- message level error
	err = client.EchoSend("blah", "blah")
	fmt.Println("err:", err)
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

	_, err = os.Stdout.Write(body)

	if err != nil {
		log.Fatal(err)
	}

	return body

}

func parseJson(resp []byte) {
	var btc BtcEntry
	json.Unmarshal(resp, &btc)

	fmt.Println("\n\n")

	for _, data := range btc.Bpis {
		fmt.Println("Name: ", btc.ChartName,
			"\nDisclaimer: ", btc.Disclaimer,
			"\nDATA:", data, "\n\n")
	}

	fmt.Println("USD Rate: ", btc.Bpis["USD"].Rate)
}
