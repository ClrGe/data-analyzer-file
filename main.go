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
type Stations []struct {
	Datasetid string `json:"datasetid"`
	Fields    struct {
		GareAliasLibelleNoncontraint string `json:"gare_alias_libelle_noncontraint"`
		CommuneLibellemin            string `json:"commune_libellemin"`
		UicCode                      string `json:"uic_code"`
		Gare                         string `json:"gare"`
		AdresseCp                    string `json:"adresse_cp"`
		DepartementLibellemin        string `json:"departement_libellemin"`
		GareRegionsncfLibelle        string `json:"gare_regionsncf_libelle"`
	} `json:"fields"`
}

// defining the station struct
type StationData [][]string

// retrieve stations ref-data from API
func getRef(w http.ResponseWriter, r *http.Request) {

	//zipcode := "76000"
	url := "http://127.0.0.1:8200/cell/test"

	var station Stations

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
	if err := json.NewDecoder(resp.Body).Decode(&station); err != nil {
		log.Println(err)
	}

	fmt.Fprintf(w, "data: %s\n", station)
}

// retrieve stations ref-data from API
func getApi(w http.ResponseWriter, r *http.Request) {

	//zipcode := "76000"
	url := "https://lab.jmg-conseil.eu/db/all"

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

	fmt.Fprintf(w, "data: %s\n", stationf)
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

// serve converted data
func test(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "data/referentiel-gares-voyageurs.json")
}

// main function
func main() {
	convertData()

	router := mux.NewRouter()
	router.HandleFunc("/cell", serveJson).Methods("GET")
	router.HandleFunc("/cell/test", test).Methods("GET")
	router.HandleFunc("/cell/api", getApi).Methods("GET")
	router.HandleFunc("/cell/ref", getRef).Methods("GET")

	log.Fatal(http.ListenAndServe(":8200", router))

}
