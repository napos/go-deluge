go-deluge
=====

A lightweight deluge library for Go

Example
-------

```go
package main

import (
	"fmt"
	"os"

	"github.com/naposproject/go-deluge"
)

func main() {
	c, err := deluge.NewClient(&deluge.Client{
		API:      "http://localhost:8085/gui",
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
	}

	for _, torrent := range torrents {
		fmt.Printf("Name: %s, Added: %d, Completed: %d, Filepath: %s\n", torrent.Name, torrent.AddedOn, torrent.CompletedOn, torrent.FilePath)
	}

	os.Exit(0)
}

```
