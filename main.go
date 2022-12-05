package main

import (
	"fmt"
	"log"
	"os"
    "github.com/gorilla/mux"
    "net/http"
	"github.com/go-gota/gota/dataframe"
)

// defining the station struct
type Stations struct {
	Nom []string `csv:nom_gare`
	UIC []int `csv:code_uic_complet`
	Total2015 []int `json:total_voyageurs_non_voyageurs_2019`
	Total2016 []int `csv:total_voyageurs_non_voyageurs_2019`
	Total2017 []int `csv:total_voyageurs_non_voyageurs_2019`
	Total2018 []int `csv:total_voyageurs_non_voyageurs_2019`
	Total2019 []int `csv:total_voyageurs_non_voyageurs_2019`
	Total2020 []int `csv:total_voyageurs_non_voyageurs_2019`
	Total2021 []int `csv:total_voyageurs_non_voyageurs_2019`
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
	log.Fatal(http.ListenAndServe(":8200", router))
}
