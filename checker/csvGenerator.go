package checker

import (
	"encoding/csv"
	"log"
	"os"
)

func ToCSV(filepath string, data [][]string) {
	file, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.WriteAll(data)
	err = writer.Error()
	if err != nil {
		log.Println(err)
	}
	return
}
