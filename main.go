package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// defining the station struct
type Stations struct {
	Nom       []string `json:"nom_gare"`
	UIC       []int    `json:"code_uic_complet"`
	Total2015 []int    `json:"total_voyageurs_non_voyageurs_2015"`
	Total2016 []int    `json:"total_voyageurs_non_voyageurs_2016"`
	Total2017 []int    `json:"total_voyageurs_non_voyageurs_2017"`
	Total2018 []int    `json:"total_voyageurs_non_voyageurs_2018"`
	Total2019 []int    `json:"total_voyageurs_non_voyageurs_2019"`
	Total2020 []int    `json:"total_voyageurs_non_voyageurs_2020"`
	Total2021 []int    `json:"total_voyageurs_non_voyageurs_2021"`
}

var station []Stations

// retrieve stations ref-data from API
func getRef(w http.ResponseWriter, r *http.Request) {

	//zipcode := "76000"
	url := "https://lab.jmg-conseil.eu/db/all"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(station)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body); err != nil {
		log.Println(err)
	}
}

// convert .csv file
func convertData() {
	file, err := os.Open("data/gares.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(file)
	columns := df.Select([]string{"nom_gare", "code_uic_complet", "total_voyageurs_non_voyageurs_2015", "total_voyageurs_non_voyageurs_2016", "total_voyageurs_non_voyageurs_2017", "total_voyageurs_non_voyageurs_2018", "total_voyageurs_non_voyageurs_2019", "total_voyageurs_non_voyageurs_2020", "total_voyageurs_non_voyageurs_2021"})

	fmt.Println(columns)
	nile, err := os.Create("data/output.json")
	if err != nil {
		log.Fatal(err)
	}
	columns.WriteJSON(nile)
}

// calculate sum of travellers for the period 2015-2021
func getSum() {
	file, err := os.Open("data/gares.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(file)
	columns := df.Select([]string{"nom_gare", "code_uic_complet", "total_voyageurs_non_voyageurs_2015", "total_voyageurs_non_voyageurs_2016", "total_voyageurs_non_voyageurs_2017", "total_voyageurs_non_voyageurs_2018", "total_voyageurs_non_voyageurs_2019", "total_voyageurs_non_voyageurs_2020", "total_voyageurs_non_voyageurs_2021"})

	fmt.Println(columns)
	nile, err := os.Create("data/output.json")
	if err != nil {
		log.Fatal(err)
	}
	columns.WriteJSON(nile)
}

// serve converted data
func serveJson(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "data/output.json")
}

// main function
func main() {
	convertData()

	router := mux.NewRouter()
	router.HandleFunc("/cell", serveJson).Methods("GET")
	router.HandleFunc("/cell/ref", getRef).Methods("GET")

	log.Fatal(http.ListenAndServe(":8200", router))
}
