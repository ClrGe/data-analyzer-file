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
type StationData [][]string

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

// serve converted data
func serveJson(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "data/output.json")
}

// retrieve stations ref-data from API
func getApi(w http.ResponseWriter, r *http.Request) {

	// retrieve request parameters
	uiccode := r.URL.Query()["uic"]
	zipcode := r.URL.Query()["zipcode"]

	// base url from which the data will be fetched
	url := "https://lab.jmg-conseil.eu/db/search?zipcode=" + zipcode[0]

	var stationf StationData

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
	if err := json.NewDecoder(resp.Body).Decode(&stationf); err != nil {
		log.Println(err)
	}

	if uiccode != nil {
		fmt.Fprintf(w, "Parametre de recherche : Code UIC %s\n", uiccode)
	}

	if zipcode != nil {
		fmt.Fprintf(w, "Parametre de recherche : Code postal %s\n", zipcode)
	}

	/* check for a matching zipcode
	var result bool = false
	for i := 0; i < len(stationf); i++ {
		if stationf[i][1] == zipcode {
			result = true
			fmt.Fprintf(w, "Code postal: %s\n", stationf[i])
			break
		}
	}*/
	fmt.Fprintf(w, "Commune: %s\n", stationf[0][0])
	fmt.Fprintf(w, "Département: %s\n", stationf[0][2])
	fmt.Fprintf(w, "Region: %s\n", stationf[0][1])
	fmt.Fprintf(w, "Code UIC: %s\n", stationf[0][4])
}

// main function
func main() {
	convertData()

	router := mux.NewRouter()
	router.HandleFunc("/cell", serveJson).Methods("GET")
	router.HandleFunc("/cell/api", getApi).Methods("GET")

	log.Fatal(http.ListenAndServe(":8200", router))

}
