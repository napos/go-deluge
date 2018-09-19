package main

import (
	"fmt"
	"os"

	"github.com/naposproject/go-deluge"
)

func main() {
	c, err := deluge.NewClient(&deluge.Client{
		API:      "http://192.168.1.101:8112/json",
		Username: "admin",
		Password: os.Getenv("TORRENT_PASSWORD"),
	})

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	fmt.Printf("Getting torrents..\n")
	torrents, err := c.GetTorrents()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	for _, torrent := range torrents {
		fmt.Printf("Hash: %s, Name: %s, Added: %d, Completed: %d, Filepath: %s\n", torrent.Hash, torrent.Name, torrent.AddedOn, torrent.CompletedOn, torrent.FilePath)
	}

	fmt.Printf("\n")
	fmt.Printf("Getting single torrent..\n")

	torrent, err := c.GetTorrent(torrents[0].Hash)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Name: %s, Added: %d, Completed: %d, Filepath: %s\n", torrent.Name, torrent.AddedOn, torrent.CompletedOn, torrent.FilePath)

	os.Exit(0)
}
