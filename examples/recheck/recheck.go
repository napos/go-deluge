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

	fmt.Printf("Stopping torrent..\n")
	err = c.StopTorrent("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	time.Sleep(5 * time.Second)

	fmt.Printf("Rechecking torrent..\n")
	err = c.RecheckTorrent("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	time.Sleep(5 * time.Second)

	fmt.Printf("Starting torrent..\n")
	err = c.StartTorrent("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	os.Exit(0)
}
