package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jamespearly/loggly"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type BtcEntry struct {
	Date       time.Time
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

type BtcItem struct {
	Date        time.Time
	Disclaimer  string
	ChartName   string
	Code        string
	Symbol      string
	Rate        string
	Description string
	Rate_Float  float32
}

func main() {
	for {
		data := work()
		svc := dynamodb.New(newSession())
		addRecord(svc, data)

		fmt.Println("Next pull in minutes at: ", time.Now().Local().Add(time.Hour*time.Duration(0)+time.Minute*time.Duration(10)+time.Second*time.Duration(0)), "\n")

		time.Sleep(10 * time.Minute)
		fmt.Println("---------------------------------------------------------------")
	}
}

func work() BtcItem {
	apiJson := getResponse("https://api.coindesk.com/v1/bpi/currentprice.json")
	data := parseJson(apiJson)
	logglyConnection(apiJson)

	return data
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

func parseJson(resp []byte) BtcItem {
	var btc BtcEntry
	json.Unmarshal(resp, &btc)

	btc.Date = time.Now()
	fmt.Println("\nKEY: ", btc.Date, "\n")

	var btcI BtcItem

	for _, data := range btc.Bpi {
		if data.Code == "USD" {
			btcI.Date = btc.Date
			btcI.ChartName = btc.ChartName
			btcI.Disclaimer = btc.Disclaimer
			btcI.Code = data.Code
			btcI.Symbol = data.Symbol
			btcI.Rate = data.Rate
			btcI.Rate_Float = data.Rate_Float
			btcI.Description = data.Description

			// fmt.Println("Name: ", btcI.ChartName,
			// 	"\nDisclaimer: ", btcI.Disclaimer,
			// 	"\nCode: [", btcI.Code, "]",
			// 	"\nSym: [", btcI.Symbol, "]",
			// 	"\nRate: $", btcI.Rate,
			// 	"\nRate in Float: $", btcI.Rate_Float,
			// 	"\nDescription: ", btcI.Description,
			// 	"\n\n")
		}
	}

	return btcI
}

func logglyConnection(resp []byte) {
	// I don't care that my token is here

	os.Setenv("LOGGLY_TOKEN", "8224b513-31a1-437c-924a-5b0e420f55ec")

	client := loggly.New("GoLoggly")
	err := client.Send("info", string(resp))

	byteValue := len(resp)
	fmt.Println("Successfully sent JSON data of", byteValue)

	if err != nil {
		log.Fatal(err)
	}
}

func newSession() *session.Session {
	// James P Early Login
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIA34XNLPJYKSVEN2G2")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return sess
}

func addRecord(svc *dynamodb.DynamoDB, data BtcItem) {
	dbItem, err := dynamodbattribute.MarshalMap(data)

	input := &dynamodb.PutItemInput{
		Item:      dbItem,
		TableName: aws.String("golang_btc"),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Successfully added record to golang_btc tbl\n")
}
