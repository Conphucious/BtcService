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

func work() BtcEntry {
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

func parseJson(resp []byte) BtcEntry {
	var btc BtcEntry
	json.Unmarshal(resp, &btc)

	btc.Date = time.Now()
	fmt.Println("\nKEY: ", btc.Date, "\n")

	return btc
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
	os.Setenv("AWS_ACCESS_KEY_ID", "ASIA2VZJBDU47GGI5DNH")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "uKiniSJWIbWYRuJUQe77CFiTdNqOQsqKdHgCz58u")
	os.Setenv("AWS_SESSION_TOKEN", "FQoGZXIvYXdzEHgaDKBr3+CQUxrWZDJ1TCLwAuky5vsgphy6B2YM+15w5N3m0F8kAs5r+CV8WLBzuEyNpH5N+ZeEZPMlDaYXF11SYTubfgTkghRsDt8J6nD899GYrxY/5zbf8gqDzQx47vPlM1HeP0nCI2jyt5EwGDNt9aeSH4dCzYEcQ5J/ppnhK0qBTSdSmWt7PUhGFKHhFk2HxOxfK7ZDtGX7SkLIWSYqjX1IVFjTulRPik4ShkYHW5y/vRBYrQrbnpK1skBUfZlGJNAWaITOT/mCDQKhhblsKWcPj1qwto2wItJ2m6QvenRqQOD8pXb+9HPVN6ubJwVPS/JlocB8kTorix22EMK6LVE4kk/JYNiRaaERkABXQlEDMhv+KITKJan86HIanjIp7QMqWP8p42v89DT1TDzupbj5CMsgMM8LqveQKGasF2BNFWiz4emR7pTXAbcKO3jUR4tYDJy6nDeIVnqRruJXQRatA5YFB4OwSAvNGV7IZ+w21GUdlRXGv65MpuW9vsjiKJu3/+MF")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return sess
}

func addRecord(svc *dynamodb.DynamoDB, data BtcEntry) {
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

//read recs
// docker server http request from db
// student ID
