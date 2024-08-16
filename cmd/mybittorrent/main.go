package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

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
		f, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		var meta TorrentFile
		if err := bencode.Unmarshal(f, &meta); err != nil {
			panic(err)
		}
		fmt.Println("Tracker URL:", meta.Announce)
		fmt.Println("Length:", meta.Info.Length)
		h := sha1.New()
		if err := bencode.Marshal(h, meta.Info); err != nil {
			panic(err)
		}
		fmt.Printf("Info Hash: %x\n", h.Sum(nil))
		fmt.Println("Piece Length:", meta.Info.PieceLength)
		printPieceHashes([]byte(meta.Info.Pieces))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
