package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Station []struct {
	CodeUic   int    `json:"code_uic_complet"`
	NomGare   string `json:"nom_gare"`
	Total2015 int    `json:"total_voyageurs_non_voyageurs_2015"`
	Total2016 int    `json:"total_voyageurs_non_voyageurs_2016"`
	Total2017 int    `json:"total_voyageurs_non_voyageurs_2017"`
	Total2018 int    `json:"total_voyageurs_non_voyageurs_2018"`
	Total2019 int    `json:"total_voyageurs_non_voyageurs_2019"`
	Total2020 int    `json:"total_voyageurs_non_voyageurs_2020"`
	Total2021 int    `json:"total_voyageurs_non_voyageurs_2021"`
}

// defining the station struct
type StationData [][]string

func csvReader(w http.ResponseWriter, r *http.Request) {
	// 1. Open the file
	recordFile, err := os.Open("data/frequentation-gares.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
	} // 2. Initialize the reader
	reader := csv.NewReader(recordFile)
	records, _ := reader.ReadAll()
	fmt.Fprint(w, records)
}

// serve converted data
func serveJson(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "data/output.json")
}

// API documentation
func serveRawDoc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger.json")
}

// retrieve stations ref-data from API
func sendData(w http.ResponseWriter, r *http.Request) {

	// retrieve request parameters
	uiccode := r.URL.Query()["uic"]
	zipcode := r.URL.Query()["zipcode"]

	// url from which the data will be fetched
	url := "http://localhost:8200/cell"

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

	var station Station

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&station); err != nil {
		log.Println(err)
	}

	if uiccode != nil {
		fmt.Printf("Parametre de recherche : Code UIC %s\n\n", uiccode)
	}

	if zipcode != nil {
		fmt.Printf("Parametre de recherche : Code postal %s\n\n", zipcode)
	}

	for i := 0; i < len(uiccode); i++ {
		var id, err = strconv.Atoi(uiccode[i])
		if err != nil {
			log.Fatal("NewRequest: ", err)
			return
		}
		for s := 0; s < len(station); s++ {
			if station[s].CodeUic == id {
				fmt.Fprint(w, station[i].Total2015, station[i].Total2016, station[i].Total2017, station[i].Total2018, station[i].Total2019, station[i].Total2020, station[i].Total2021, "&")
			}
		}
	}

}

// main function
func main() {

	router := mux.NewRouter()

	router.HandleFunc("/cell/station", sendData).Methods("GET")
	router.HandleFunc("/cell/csv", csvReader).Methods("GET")
	router.HandleFunc("/cell", serveJson).Methods("GET")
	router.HandleFunc("/cell/raw", serveRawDoc).Methods("GET")

	log.Fatal(http.ListenAndServe(":8200", router))

}
