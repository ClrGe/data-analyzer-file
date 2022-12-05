package main

import (
	"fmt"
	"log"
	"os"
	"github.com/go-gota/gota/dataframe"
)

type Stations struct {
	UIC []string `csv:[fields.code_uic]`
}

func main() {

	file, err := os.Open("data/gares.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(file)
	columns := df.Select([]string{"fields/nom_gare", "fields/code_uic_complet", "fields/total_voyageurs_non_voyageurs_2019"})
	fmt.Println(columns)

}
