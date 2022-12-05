package main

import (
	"fmt"
	"log"
	"os"
    "github.com/gorilla/mux"
    "net/http"
	"github.com/go-gota/gota/dataframe"
)

type Stations struct {
	Nom []string `csv:fields/nom_gare`
	UIC []string `csv:fields/code_uic_complet`
	Total2019 []string `csv:fields/total_voyageurs_non_voyageurs_2019`
}


func serveJson(w http.ResponseWriter, r *http.Request) {

    // b,_ := ioutil.ReadFile("data/output.json");

    // rawIn := json.RawMessage(string(b))
    // var objmap map[string]*json.RawMessage
    // err := json.Unmarshal(rawIn, &objmap)
    // if err != nil {
    //   fmt.Println(err)
    // }
    // fmt.Println(objmap)

    // json.NewEncoder(w).Encode(objmap)

	http.ServeFile(w, r, "data/output.json")

}

func main() {

	file, err := os.Open("data/gares.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	df := dataframe.ReadCSV(file)
	columns := df.Select([]string{"fields/nom_gare", "fields/code_uic_complet", "fields/total_voyageurs_non_voyageurs_2015", "fields/total_voyageurs_non_voyageurs_2016", "fields/total_voyageurs_non_voyageurs_2017", "fields/total_voyageurs_non_voyageurs_2018", "fields/total_voyageurs_non_voyageurs_2019", "fields/total_voyageurs_non_voyageurs_2020", "fields/total_voyageurs_non_voyageurs_2021"})
	fmt.Println(columns)

	nile, err := os.Create("data/output.json")
	if err != nil {
		log.Fatal(err)
	}

	columns.WriteJSON(nile)

	router := mux.NewRouter()
    router.HandleFunc("/cell", serveJson).Methods("GET")
	log.Fatal(http.ListenAndServe(":8200", router))
}
