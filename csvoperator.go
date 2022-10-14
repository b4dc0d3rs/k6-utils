package k6utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
)

func (k6utils *K6Utils) TakeRandomRow() (interface{}, error) {
	if k6utils.csvRecords == nil || len(k6utils.csvRecords) == 0 {
		return nil, fmt.Errorf("%s", "List of CSV records not loaded, call Load method first")
	}

	randomIndex := rand.Intn(len(k6utils.csvRecords))
	return k6utils.csvRecords[randomIndex], nil
}

func (k6utils *K6Utils) Load(filePath string, separator []byte) (interface{}, error) {

	if len(separator) == 0 {;
		return nil, fmt.Errorf("%s", "Separator not valid")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	reader := csv.NewReader(file)

	header, err :=	reader.Read();

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	k6utils.csvRecords = []map[string]string{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if len(record) == 0 {
			continue
		}

		dict := map[string]string{}
		for i := range record {
			dict[header[i]] = record[i]
		}
		k6utils.csvRecords = append(k6utils.csvRecords, dict)

	}

	file.Close()

	return k6utils.csvRecords, nil
}
