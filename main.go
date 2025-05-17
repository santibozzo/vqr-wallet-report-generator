package main

import (
	"fmt"
	"os"
	"regexp"
	"encoding/csv"
	"strings"
)

type Record struct {
	Id string `json:id`
	Name string `json:name`
}

func main() {
	fmt.Println("hola")

	// CSV reader
	inputFileName := get_file_name()
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	var recordsToSave []Record

	count := 0
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if count == 0 {
			count++
			continue
		}

		record := Record{
			Id: row[0],
			Name: row[1],
		}

		recordsToSave = append(recordsToSave, record)
	}

	// CSV writer
	outputFileName := strings.ReplaceAll(inputFileName, "report", "wallet_report")
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	writer.Comma = ';'
	defer writer.Flush()

	header := []string{"id", "name"}
	_ = writer.Write(header)

	for _, record := range recordsToSave {
		row := []string{
			record.Id,
			record.Name,
		}

		if err := writer.Write(row); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	
	fmt.Println("listo")
}

func get_file_name() string {
	regex, _ := regexp.Compile(`^.*ARS_report\.csv$`)

	files, _ := os.ReadDir(".")
	var file_name string

	for _, file := range files {
		if file.Type().IsRegular() && regex.MatchString(file.Name()) {
			file_name = file.Name()
		}
	}

	return file_name
}