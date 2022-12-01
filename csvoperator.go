package k6utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
)

func (k6utils *K6Utils) TakeRandomRow() (map[string]string, error) {
	if k6utils.csvRecords == nil || len(k6utils.csvRecords) == 0 {
		return nil, fmt.Errorf("%s", "List of CSV records not loaded, call Load method first")
	}

	randomIndex := rand.Intn(len(k6utils.csvRecords))

	record := k6utils.csvRecords[randomIndex]
	dict := map[string]string{}
	for i := range record {
		dict[k6utils.header[i]] = record[i]
	}

	return dict, nil
}

func (k6utils *K6Utils) PollRandomRow() (map[string]string, error) {
	if k6utils.csvRecords == nil || len(k6utils.csvRecords) == 0 {
		return nil, fmt.Errorf("%s", "List of CSV records not loaded, call Load method first")
	}

	randomIndex := rand.Intn(len(k6utils.csvRecords))

	record := k6utils.csvRecords[randomIndex]
	k6utils.csvRecords = removeElementByIndex(k6utils.csvRecords, randomIndex)

	dict := map[string]string{}
	for i := range record {
		dict[k6utils.header[i]] = record[i]
	}

	return dict, nil
}

func removeElementByIndex[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func (k6utils *K6Utils) TakeRowByIndex(index int) (map[string]string, error) {
	if k6utils.csvRecords == nil || len(k6utils.csvRecords) == 0 {
		return nil, fmt.Errorf("%s", "List of CSV records not loaded, call Load method first")
	}

	if index < 0 || index > len(k6utils.csvRecords) {
		return nil, fmt.Errorf("Invalid index %d. CSV feed has %d rows", index, len(k6utils.csvRecords))
	}

	record := k6utils.csvRecords[index]
	dict := map[string]string{}
	for i := range record {
		dict[k6utils.header[i]] = record[i]
	}

	return dict, nil
}

func (k6utils *K6Utils) Load(filePath string, separator []byte) (interface{}, error) {

	if len(separator) == 0 {
		return nil, fmt.Errorf("%s", "Separator not valid")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	reader := csv.NewReader(file)

	header, err := reader.Read()
	k6utils.header = header

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	k6utils.csvRecords = [][]string{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}

		// skip empty lines
		if len(record) == 0 {
			continue
		}

		k6utils.csvRecords = append(k6utils.csvRecords, record)
	}

	file.Close()

	return k6utils.csvRecords, nil
}
