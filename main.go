
// This script creates a web server using the net/http package and the gorilla/mux router package,
// and exposes several routes to handle different types of requests (csv file, serving json and API documentation, retrieving data from an API..)

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
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
type StationE []struct {
	CodeUic   int    `json:"code_uic"`
	NomGare   string `json:"cp"`
	Total2015 int    `json:"a2015"`
	Total2016 int    `json:"a2016"`
	Total2017 int    `json:"a2017"`
	Total2018 int    `json:"a2018"`
	Total2019 int    `json:"2019"`
	Total2020 int    `json:"a2020"`
	Total2021 int    `json:"a2021"`
}

// Config struct holds the environment variables
type Config struct {
	PORT string `mapstructure:"PORT"`
	HOST string `mapstructure:"HOST"`
}

// LoadConfig loads the env file data to a struct
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

// csvReader function reads a CSV file and returns the records to the client
// func csvReader(w http.ResponseWriter, r *http.Request) {
// 	// retrieve request parameters
// 	uiccode := r.URL.Query()["uic"]
// 	zipcode := r.URL.Query()["zipcode"]

// 	in, err := os.Open("data/gares.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer in.Close()

// 	stations := []Station{}

// 	if uiccode != nil {
// 		fmt.Printf("Parametre de recherche : Code UIC %s\n\n", uiccode)
// 	}
// 	if zipcode != nil {
// 		fmt.Printf("Parametre de recherche : Code postal %s\n\n", zipcode)
// 	}

// 	// Unmarshal the CSV data into the stations variable
// 	gocsv.UnmarshalFile(in, stations)

// 	for _, t := range stations {
// 		json.NewEncoder(w).Encode(t)
// 	}
// }

// serveJson function serves a JSON file to the client
func serveJson(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "data/output.json")
}

// serveRawDoc function serves the API documentation to the client
func serveRawDoc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "swagger.json")
}

// sendData function retrieves data from an API based on request parameters 'uic' and 'zipcode' and returns the data in json format
func sendData(w http.ResponseWriter, r *http.Request) {

	// retrieve request parameters
	uiccode := r.URL.Query()["uic"]
	zipcode := r.URL.Query()["zipcode"]

	// url from which the data will be fetched
	url := "https://lab.jmg-conseil.eu/cell"

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Requete : ", err)
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

	// looping over the received data to find the station with the matching uic code
	for i := 0; i < len(uiccode); i++ {
		var id, err = strconv.Atoi(uiccode[i])
		if err != nil {
			log.Fatal("Requete: ", err)
			return
		}
		for s := 0; s < len(station); s++ {
			if station[s].CodeUic == id {
				fmt.Fprint(w, station[i].Total2015, station[i].Total2016, station[i].Total2017, station[i].Total2018, station[i].Total2019, station[i].Total2020, station[i].Total2021, "&")
			}
		}
	}
}

func main() {
	// load app.env file data to struct
	config, err := LoadConfig(".")

	router := mux.NewRouter()

	router.HandleFunc("/cell/station", sendData).Methods("GET")
	//router.HandleFunc("/cell/csv", csvReader).Methods("GET")
	router.HandleFunc("/cell", serveJson).Methods("GET")
	router.HandleFunc("/cell/raw", serveRawDoc).Methods("GET")
	log.Fatal(http.ListenAndServe(config.PORT, router))

	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}
}
