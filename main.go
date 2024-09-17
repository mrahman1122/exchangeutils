package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	flag.Parse()

	if flag.NArg() != 3 {
		// error msg here with info on expected cmd line parameters
		fmt.Println("cryptocalc expects 3 arguments, defined in order: an integer for USD amount, and two strings for the currency ticker symbols")
		fmt.Println("e.g.: 100 BTC ETH")
		os.Exit(1)
	}

	var firstCurrency = string(flag.Arg(1))
	var secondCurrency = string(flag.Arg(2))
	var usdAmount, err = strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Println("Error: Cannot parse USD Amount properly -- please provide an integer")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("USD: %d!\n", usdAmount)
	fmt.Printf("First Currency: %s!\n", firstCurrency)
	fmt.Printf("Second Currency: %s!\n", secondCurrency)

	url := "https://api.coinbase.com/v2/exchange-rates?currency=USD"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error: Could not execute HTTP request to Coinbase API")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error: Could not retrieve exchange data from Coinbase API")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: Could not properly read Coinbase API reaponse body")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Parsed Response Body!")
	//fmt.Println(string(body))

	var objmap map[string]json.RawMessage
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		fmt.Println("JSON Error!")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	fmt.Println("Parsed JSON!")

	var dataObj map[string]json.RawMessage
	err = json.Unmarshal(objmap["data"], &dataObj)
	if err != nil {
		fmt.Println("Data Error!")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	fmt.Println("Parsed Data!")

	var rates map[string]string
	err = json.Unmarshal(dataObj["rates"], &rates)
	if err != nil {
		fmt.Println("Rates Error!")
		fmt.Print(err.Error())
		os.Exit(1)
	}

	fmt.Println("Parsed Rates!")

	firstValue, exists := rates[firstCurrency]
	if exists {
		// valid ticker symbol
	} else {
		fmt.Printf("%s is not a valid currency ticket symbol\n", firstCurrency)
		os.Exit(1)
	}
	secondValue, exists := rates[secondCurrency]
	if exists {
		// valid ticker symbol
	} else {
		fmt.Printf("%s is not a valid currency ticket symbol\n", secondCurrency)
		os.Exit(1)
	}

	firstRate, err := strconv.ParseFloat(firstValue, 64)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("%s Rate: %f\n", firstCurrency, firstRate)
	}

	secondRate, err := strconv.ParseFloat(secondValue, 64)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("%s Rate: %f\n", secondCurrency, secondRate)
	}

	firstAmount := float64(usdAmount) * firstRate * 0.7
	fmt.Printf("$70.00 => %f %s\n", firstAmount, firstCurrency)

	secondAmount := float64(usdAmount) * secondRate * 0.3
	fmt.Printf("$30.00 => %f %s\n", secondAmount, secondCurrency)

	os.Exit(0)
}
