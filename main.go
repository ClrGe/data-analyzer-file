package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Station []struct {
	CodeUic   string `json:"code_uic_complet"`
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

// convert .csv file
func convertData(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data/gares.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(file)
	df.SetNames("dataset", "record", "2018", "2021", "2015", "seg_drg", "2020", "2016", "2018", "voy2017", "2017", "2019", "code_uic", "cp", "2020", "voy2019", "voy2021", "voy2016", "voy2015", "nom_gare", "date")

	data := df.Select([]int{2, 3, 4, 6, 7, 8, 10, 11, 14})

	fmt.Fprint(w, data)

}

// serve converted data
func serveJson(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "data/output.json")
}

// retrieve stations ref-data from API
func refData(w http.ResponseWriter, r *http.Request) {

	// retrieve request parameters
	uic := r.URL.Query()["uic"]
	zip := r.URL.Query()["zip"]

	//var id, err = strconv.Atoi(uic[0])

	//  url from which the referential data will be fetched
	base := "https://lab.jmg-conseil.eu/db/search?uiccode=%s"
	url := fmt.Sprintf("%s%s", base, zip)

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
	if err := json.NewDecoder(resp.Body).Decode(&stationf); err != nil {
		log.Println(err)
	}

	if uic != nil {
		fmt.Fprintf(w, "Parametre de recherche :  %s\n", uic[0])
		fmt.Fprint(w, stationf)

	}

	// check for a matching zipcode
	for i := 0; i < len(stationf); i++ {
		if stationf[i][2] == zip[0] {
			fmt.Fprint(w, stationf[i])
		}
	}

}

// retrieve stations ref-data from API
func crowdData(w http.ResponseWriter, r *http.Request) {

	// retrieve request parameters
	uiccode := r.URL.Query()["uic"]
	zipcode := r.URL.Query()["zipcode"]

	// url from which the data will be fetched
	url := "https://lab.jmg-conseil.eu/cell"

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
	/*if err := json.NewDecoder(resp.Body).Decode(&station); err != nil {
		log.Println(err)
	}*/

	if uiccode != nil {
		fmt.Printf("Parametre de recherche : Code UIC %s\n\n", uiccode)
	}

	if zipcode != nil {
		fmt.Printf("Parametre de recherche : Code postal %s\n\n", zipcode)
	}
	if station != nil {
		fmt.Printf("Parametre de recherche : Code postal %s\n\n", zipcode)
	}

	/*	for i := 0; i < len(station); i++ {
		if station[i].CodeUic == uiccode[0] {
			//fmt.Fprint(w, station[i].Total2015, station[i].Total2016, station[i].Total2017, station[i].Total2018, station[i].Total2019, station[i].Total2020, station[i].Total2021)
			fmt.Println(refData)
			fmt.Fprint(w, station[i].Total2015, station[i].Total2016, station[i].Total2017, station[i].Total2018, station[i].Total2019, station[i].Total2020, station[i].Total2021)
			return
		}
	}*/

	refData(w, r)
}

// main function
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/cell/station", crowdData).Methods("GET")
	router.HandleFunc("/cell/csv", csvReader).Methods("GET")
	router.HandleFunc("/cell", serveJson).Methods("GET")
	router.HandleFunc("/cell/api", refData).Methods("GET")

	log.Fatal(http.ListenAndServe(":8200", router))

}
