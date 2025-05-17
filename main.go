package main

import (
	"fmt"
	"os"
	"regexp"
	"encoding/csv"
	"strings"
	"strconv"
)

// TYPES

type Record struct {
	RideId string `json:ride_id`
	PaymentId string `json:payment_id`
	ExternalReference string `json:external_reference`
	NetAmount float64 `json:net_amount`
	GrossAmount float64 `json:gross_amount`
	Fee string `json:fee`
	Currency string `json:currency`
	Status string `json:status`
	StatusDetail string `json:status_detail`
	IssuerId string `json:issuer_id`
	TransportOperator string `json:transport_operator`
	Debt string `json:debt`
	Forced string `json:forced`
	FeatureFlags string `json:feature_flags`
	ScannedAt string `json:scanned_at`
	CreatedAt string `json:created_at`
	ProcessedAt string `json:processed_at`
}

type Results struct {
	InputCount int
	OutputCount int
	GrossAmount float64
	NetAmount float64
}

// MAIN

func main() {
	fmt.Println("START")

	inputFileName := getFileName()
	results := readAndWriteCSV(inputFileName)
	writeResults(inputFileName, results)

	fmt.Println("END")
}

// FUNCTIONS

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

func isApproved(record Record) bool {
	return strings.HasPrefix(record.Status, "APPROVED")
}

func readAndWriteCSV(inputFileName string) Results {
	fmt.Println("Procesando " + inputFileName)
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

	header := []string{
		"ride_id",
		"payment_id",
		"external_reference",
		"net_amount",
		"gross_amount",
		"fee",
		"currency",
		"status",
		"status_detail",
		"issuer_id",
		"transport_operator",
		"debt",
		"forced",
		"feature_flags",
		"scanned_at",
		"created_at",
		"processed_at",
	}
	_ = writer.Write(header)


	// Process
	inputCount := 0
	outputCount := 0
	grossAmount := 0.0
	netAmount := 0.0

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

		recordNetAmount, _ := strconv.ParseFloat(row[3], 64)
		recordGrossAmount, _ := strconv.ParseFloat(row[4], 64)
		record := Record{
			RideId:            row[0],
			PaymentId:         row[1],
			ExternalReference: row[2],
			NetAmount:         recordNetAmount,
			GrossAmount:       recordGrossAmount,
			Fee:               row[5],
			Currency:          row[6],
			Status:            row[7],
			StatusDetail:      row[8],
			IssuerId:          row[9],
			TransportOperator: row[10],
			Debt:              row[11],
			Forced:            row[12],
			FeatureFlags:      row[13],
			ScannedAt:         row[14],
			CreatedAt:         row[15],
			ProcessedAt:       row[16],
		}

		if isApproved(record) {
			if err := writer.Write(row); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outputCount++
			netAmount += record.NetAmount
			grossAmount += record.GrossAmount
		}
	}

	fmt.Println("Se procesaron CSVs")
	results := Results{
		InputCount: inputCount,
		OutputCount: outputCount,
		GrossAmount: grossAmount,
		NetAmount: netAmount,
	}
	return results
}

func writeResults(inputFileName string, results Results) {
	resultsFile, err := os.Create(strings.ReplaceAll(inputFileName, "report.csv", "wallet_report-RESULTS.txt"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultsFile.Close()

	// CSVs info
	_, _ = fmt.Fprintf(
		resultsFile,
		"Input count= %d\nOutput count= %d\nGross amount= %f\nNet amount=%f\n\n",
		results.InputCount,
		results.OutputCount,
		results.GrossAmount,
		results.NetAmount,
	)

	// Payment info
	fee := results.GrossAmount * 0.004
	iva := fee * 0.21
	totalToPay := results.GrossAmount - fee - iva

	_, _ = fmt.Fprintf(resultsFile, "Comisión= %f (0.4%% = bruto * 0.004)\n", fee)
	_, _ = fmt.Fprintf(resultsFile, "IVA= %f (21%% = comisión * 0.21)\n", iva)
	_, _ = fmt.Fprintf(resultsFile, "Total a pagar= %f (bruto - comisión - IVA)\n", totalToPay)

	// Logs
	fmt.Printf("Input count= %d\n", results.InputCount)
	fmt.Printf("Output count= %d\n", results.OutputCount)
	fmt.Printf("Gross amount= %f\n", results.GrossAmount)
	fmt.Printf("Net amount= %f\n\n", results.NetAmount)
	fmt.Printf("Fee= %f\n", fee)
	fmt.Printf("IVA= %f\n", iva)
	fmt.Printf("Total to pay= %f\n", totalToPay)
}