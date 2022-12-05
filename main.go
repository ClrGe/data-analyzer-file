package main

import (
	"fmt"
	"log"
	"os"
	"github.com/go-gota/gota/dataframe"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"time"
)

type Stations struct {
	UIC []string `json:[fields.code_uic]`
}

func main() {
	
	url := "https://lab.jmg-conseil.eu/db/all"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	spaceClient := http.Client{
		Timeout: time.Second * 232, // Timeout after 2 seconds
	}

	req.Header.Set("User-Agent", "service")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	station := Stations{}

	jsonErr := json.Unmarshal(body, &station)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(station.UIC)

	file, err := os.Open("data/frequentation-gares.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(file)

	fmt.Println(df)
}
