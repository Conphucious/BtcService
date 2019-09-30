package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//	"os"
)

type BpiEntry struct {
	//	Time    string `json:"time:"` // time in UTC use updated?
	//	UsdRate string `json:"bpi:USD:rate"`
	//	GbpRate string `json:"bpi:GBP:rate"`
	//	EurRate string `json:"bpi:EUR:rate"`

	ChartName  string `json:"chartName"`
	Bpi        string `json:"code"`
	Time       string `json:"time"`
	Disclaimer string `json:"disclaimer"`
}

// access encoded field within struct so Capital letter like Name string
// marshall into struct using encoding/json
//Time string 'json: "name"'

func main() {
	str := getResponse("https://api.coindesk.com/v1/bpi/currentprice.json")

	parseJson(str)
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

	//_, err = os.Stdout.Write(body)

	if err != nil {
		log.Fatal(err)
	}

	return body
}

func parseJson(resp []byte) {
	//byteValue, _ := ioutil.ReadAll(resp)

	bpi = make(map[string]BpiEntry)

	err := json.Unmarshal(resp, &bpi)

	fmt.Println(bpi, err)

	//	for i := 0; i < len(bpi.BpiEntry); i++ {
	//	fmt.Println("BPI " + bpi.Bpi[i])
	//fmt.Println("BPI Time: " + bpi.Bpi[i].Code)
	//	fmt.Println("User Age: " + strconv.Itoa(users.Users[i].Age))
	//	fmt.Println("User Name: " + users.Users[i].Name)
	//	fmt.Println("Facebook Url: " + users.Users[i].Social.Facebook)
	//	}
}
