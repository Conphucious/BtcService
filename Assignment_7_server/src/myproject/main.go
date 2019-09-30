package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"encoding/json"

	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

type BtcDb struct {
	RequestTime time.Time
	Items       []BtcItem
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

type Endpt struct {
	Table       string
	RecordCount int
}

func main() {
	r := mux.NewRouter()
	print("Server Started!")
	r.HandleFunc("/pnguyen3/{status}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		status := vars["status"]

		items := readData(dynamodb.New(newSession()))
		print("\nLENGTH: ", len(items))

		if status == "all" {
			dbItems := BtcDb{RequestTime: time.Now(), Items: items}
			jsonFile, _ := json.Marshal(dbItems)
			fmt.Fprintf(w, "%s", jsonFile)
		} else if status == "status" {
			endPt := Endpt{Table: "golang_btc", RecordCount: len(items)}
			jsonFile, _ := json.Marshal(endPt)
			fmt.Fprintf(w, "%s", jsonFile)
		} else if status != "all" || status != "status" {
			fmt.Fprintf(w, "BOOTY HOLE R GUD BUT DIS LINK AINT")
		}
	})

	http.ListenAndServe(":8080", r)
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

func readData(svc *dynamodb.DynamoDB) []BtcItem {
	// Make the DynamoDB Query API call

	params := &dynamodb.ScanInput{
		TableName: aws.String("golang_btc"),
	}

	result, err := svc.Scan(params)
	if err != nil {
		fmt.Errorf("failed to make Query API call, %v", err)
	}

	items := []BtcItem{}

	// Unmarshal the Items field in the result value to the Item Go type.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		fmt.Errorf("failed to unmarshal Query result items, %v", err)
	}

	// Print out the items returned
	// for i, item := range items {
	// 	fmt.Printf("%d: Date: %s, Rate: %s\n", i, item.Date, item.Rate)
	// }

	return items
}
