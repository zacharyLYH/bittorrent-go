package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func decodeBencodeToTorrentFile(fileData string) (TorrentFile, error) {
	decoded, err := decodeBencode(string(fileData))
	if err != nil {
		fmt.Println(err)
		return TorrentFile{}, err
	}

	var jsonOutput TorrentFile
	jsonBytes, _ := json.Marshal(decoded)
	err = json.Unmarshal(jsonBytes, &jsonOutput)
	if err != nil {
		fmt.Println("Error mapping to struct:", err)
		return TorrentFile{}, err
	}
	return jsonOutput, nil
}

func readFileReturnBytes(fileName string) ([]byte, error) {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

func handleErrorGeneric(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
