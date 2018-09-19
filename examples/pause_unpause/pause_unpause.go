package main

import (
	"fmt"
	"os"
	"time"

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

	fmt.Printf("Pausing torrent..\n")
	err = c.PauseTorrent("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	time.Sleep(10 * time.Second)

	fmt.Printf("Unpausing torrent..\n")
	err = c.UnPauseTorrent("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	os.Exit(0)
}
