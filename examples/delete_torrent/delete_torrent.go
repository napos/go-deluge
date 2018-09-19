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

	fmt.Printf("Adding torrent via URL..\n")
	err = c.AddTorrent("http://releases.ubuntu.com/18.04/ubuntu-18.04.1-desktop-amd64.iso.torrent")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)

	fmt.Printf("Deleting torrent...\n")
	err = c.RemoveTorrent("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)

	fmt.Printf("ReAdding torrent via URL..\n")
	err = c.AddTorrent("http://releases.ubuntu.com/18.04/ubuntu-18.04.1-desktop-amd64.iso.torrent")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	time.Sleep(10 * time.Second)

	fmt.Printf("Deleting torrent and data...\n")
	err = c.RemoveTorrentAndData("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
