package main

import (
	"fmt"
	"log"
	"os"
	"github.com/go-gota/gota/dataframe"
)

func main() {
	
	file, err := os.Open("data/frequentation-gares.csv")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	df := dataframe.ReadCSV(file)

	fmt.Println(df)
}
