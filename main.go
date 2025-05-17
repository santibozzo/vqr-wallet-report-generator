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

type Results struct {
	InputCount int
	OuputCount int
}

func main() {
	fmt.Println("START")

	inputFileName := getFileName()
	results := readAndWriteCSV(inputFileName)

	fmt.Println(results.OuputCount)
	fmt.Println("END")
}

func getFileName() string {
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

func readAndWriteCSV(inputFileName string) Results {
	// Reader
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'


	// Writer
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


	// Process
	inputCount := 0
	outputCount := 0
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if inputCount == 0 {
			inputCount++
			continue
		}
		inputCount++

		// record := Record{
		// 	Id: row[0],
		// 	Name: row[1],
		// }

		// TODO cuentas/validaciones


		if err := writer.Write(row); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		outputCount++
	}

	fmt.Println("Se procesaron CSVs")
	results := Results{
		InputCount: inputCount,
		OuputCount: outputCount,
	}
	return results
}