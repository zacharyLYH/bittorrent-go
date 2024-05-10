package main

import (
	// Uncomment this line to pass the first stage

	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	bencode "github.com/jackpal/bencode-go"
)

// TorrentInfo represents the "info" field in the Bencoded file
type TorrentInfo struct {
	Length      int    `json:"length"`
	Name        string `json:"name"`
	PieceLength int    `json:"piece length"`
	Pieces      string `json:"pieces"`
}

// TorrentFile represents the entire Bencoded file
type TorrentFile struct {
	Announce  string      `json:"announce"`
	CreatedBy string      `json:"created by"`
	Info      TorrentInfo `json:"info"`
}

func decodeList(bencodedString string) ([]any, error) {
	list, err := bencode.Decode(strings.NewReader(bencodedString))
	return list.([]any), err
}

func decodeDict(bencodedString string) (map[string]interface{}, error) {
	list, err := bencode.Decode(strings.NewReader(bencodedString))
	return list.(map[string]interface{}), err
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error) {
	if bencodedString[0] == 'l' {
		return decodeList(bencodedString)
	} else if bencodedString[0] == 'd' {
		return decodeDict(bencodedString)
	} else if bencodedString[0] == 'i' && bencodedString[len(bencodedString)-1] == 'e' {
		numberStr := bencodedString[1 : len(bencodedString)-1]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return nil, err
		}
		return number, nil
	} else if unicode.IsDigit(rune(bencodedString[0])) {
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else {
		return "", fmt.Errorf("only strings are supported at the moment")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {

		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)

		handleErrorGeneric(err)

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))

	} else if command == "info" {

		torrentFile := os.Args[2]

		fileData, err := readFileReturnBytes(torrentFile)

		handleErrorGeneric(err)

		decoded, err := decodeBencodeToTorrentFile(string(fileData))

		handleErrorGeneric(err)

		fmt.Printf("Tracker URL: %s\n", decoded.Announce)
		fmt.Printf("Length: %d", decoded.Info.Length)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
