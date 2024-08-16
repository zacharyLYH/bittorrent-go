package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	bencode "github.com/jackpal/bencode-go"
)

// TorrentInfo represents the "info" field in the Bencoded file
type TorrentInfo struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

// TorrentFile represents the entire Bencoded file
type TorrentFile struct {
	Announce string      `bencode:"announce"`
	Info     TorrentInfo `bencode:"info"`
}

func handleErrorGeneric(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printPieceHashes(pieces []byte) {
	const pieceSize = 20 // Each piece is 20 bytes
	numPieces := len(pieces) / pieceSize
	fmt.Println("Piece Hashes: ")
	for i := 0; i < numPieces; i++ {
		start := i * pieceSize
		end := start + pieceSize
		piece := pieces[start:end]
		fmt.Printf("%s\n", hex.EncodeToString(piece))
	}
}

func bencodeUnmarshall(fileName string) TorrentFile {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	var meta TorrentFile
	if err := bencode.Unmarshal(f, &meta); err != nil {
		panic(err)
	}
	return meta
}

func getInfoHash(info TorrentInfo) []byte {
	h := sha1.New()
	if err := bencode.Marshal(h, info); err != nil {
		panic(err)
	}
	return h.Sum(nil)
}

func getRequest(endpoint string, parameters map[string]string) []byte {
	params := url.Values{}
	for k, v := range parameters {
		params.Add(k, v)
	}
	finalURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	resp, err := http.Get(finalURL)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	return body

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {

		bencodedValue := os.Args[2]

		decoded, err := bencode.Decode(bytes.NewReader([]byte(bencodedValue)))

		handleErrorGeneric(err)

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))

	} else if command == "info" {
		fileName := os.Args[2]
		meta := bencodeUnmarshall(fileName)
		fmt.Println("Tracker URL:", meta.Announce)
		fmt.Println("Length:", meta.Info.Length)
		fmt.Printf("Info Hash: %x\n", getInfoHash(meta.Info))
		fmt.Println("Piece Length:", meta.Info.PieceLength)
		printPieceHashes([]byte(meta.Info.Pieces))
	} else if command == "peers" {
		fileName := os.Args[2]
		meta := bencodeUnmarshall(fileName)
		length := strconv.Itoa(meta.Info.Length)
		infoHash := fmt.Sprintf("%x", getInfoHash(meta.Info))
		infoHashBytes, _ := hex.DecodeString(infoHash)
		params := map[string]string{
			"info_hash":  string(infoHashBytes),
			"peer_id":    "00112233445566778899",
			"port":       "6881",
			"uploaded":   "0",
			"downloaded": "0",
			"left":       length,
			"compact":    "1",
		}
		body := getRequest(meta.Announce, params)
		decoded, err := bencode.Decode(bytes.NewReader(body))
		handleErrorGeneric(err)
		decodedMap, ok := decoded.(map[string]interface{})
		if !ok {
			fmt.Println("Decoded value is not a map")
		}
		peerData := []byte(decodedMap["peers"].(string))
		for i := 0; i < len(peerData); i += 6 {
			ip := net.IPv4(peerData[i], peerData[i+1], peerData[i+2], peerData[i+3])
			port := int(peerData[i+4])<<8 | int(peerData[i+5])
			fmt.Printf("%s:%d\n", ip, port)
		}
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
