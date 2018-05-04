package checker

import (
	"encoding/csv"
	"log"
	"os"
)

func openCSV(filepath string) [][]string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
	}
	return data
}
