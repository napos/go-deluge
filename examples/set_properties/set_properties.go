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

	// Wait for torrent to be added/started
	time.Sleep(5 * time.Second)

	fmt.Printf("Setting torrent label..\n")
	err = c.SetTorrentLabel("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5", "OS")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Setting torrent Seed Ratio..\n")
	err = c.SetTorrentSeedRatio("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5", 5.2)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Setting torrent queue priority..\n")
	err = c.QueueTop("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Setting torrent queue priority..\n")
	err = c.QueueDown("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Setting torrent queue priority..\n")
	err = c.QueueUp("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Setting torrent queue priority..\n")
	err = c.QueueBottom("c3a41b13a4607f0b3188063aa5fb8a50e02ac4f5")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
