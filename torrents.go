package deluge

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type RpcError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type BoolResponse struct {
	Id     int      `json:"id"`
	Result bool     `json:"result"`
	Error  RpcError `json:"error"`
}

type StringResponse struct {
	Id     int      `json:"id"`
	Result string   `json:"result"`
	Error  RpcError `json:"error"`
}

// TorrentsResponse is an interface that allows unmarshalling of the
// deluge/Bittorrent api into proper golang compatible Torrent structs.
type TorrentsResponse struct {
	Index       int                `json:"id"`
	RawTorrents map[string]Torrent `json:"result"`
	Torrents    []Torrent
	Error       RpcError `json:"error"`
}

// TorrentResponse is an interface that allows unmarshalling of the
// deluge/Bittorrent api into proper golang compatible Torrent structs.
type TorrentResponse struct {
	Index   int      `json:"id"`
	Torrent Torrent  `json:"result"`
	Error   RpcError `json:"error"`
}

// TorrentProperties is a string containing the json keys to grab via JSONRPC
// it vastly speeds up the call due to only grabbing the required values
var TorrentProperties string = "\"hash\", \"name\", \"total_size\", \"progress\", \"all_time_download\", \"total_uploaded\", \"ratio\", \"upload_payload_rate\", \"download_payload_rate\", \"eta\", \"label\", \"num_peers\", \"total_peers\", \"num_seeds\", \"total_seeds\", \"seeds_peers_ratio\", \"queue\", \"state\", \"time_added\", \"move_on_completed_path\""

type Torrent struct {
	Hash            string  `json:"hash"`
	StatusCode      int     `json:"status_code"`
	Name            string  `json:"name"`
	Size            int     `json:"total_size"`
	PercentProgress float64 `json:"progress"`
	Downloaded      int     `json:"all_time_download"`
	Uploaded        int     `json:"total_uploaded"`
	Ratio           float64 `json:"ratio"`
	UploadSpeed     int     `json:"upload_payload_rate"`
	DownloadSpeed   int     `json:"download_payload_rate"`
	ETA             int     `json:"eta"`
	Label           string  `json:"label"`
	PeersConnected  int     `json:"num_peers"`
	PeersTotal      int     `json:"total_peers"`
	SeedsConnected  int     `json:"num_seeds"`
	SeedsTotal      int     `json:"total_seeds"`
	Availability    float64 `json:"seeds_peers_ratio"`
	QueueOrder      int     `json:"queue"`
	Remaining       int     `json:"remaining"`
	Status          string  `json:"state"`
	AddedOn         int     `json:"added_on"`
	CompletedOn     int     `json:"completed_on"`
	FilePath        string  `json:"move_on_completed_path"`
	AddedRaw        float64 `json:"time_added"`
}

// max is a simple function used to bound the remaining value (as it can go negative)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// UnmarshallJSON is a custom unmarshaller for torrent lists. Necessary due to
// the fact uTorrent/Bittorrent does not implement a proper json api.
func (torrents *TorrentsResponse) UnmarshalJSON(b []byte) error {
	type Alias TorrentsResponse
	rawTorrents := &struct {
		*Alias
	}{
		Alias: (*Alias)(torrents),
	}

	err := json.Unmarshal(b, &rawTorrents)
	if err != nil {
		return err
	}

	for _, torrent := range rawTorrents.RawTorrents {
		torrents.Torrents = append(torrents.Torrents, Torrent{
			Hash:            torrent.Hash,
			StatusCode:      200, //OK? - Not Provided
			Name:            torrent.Name,
			Size:            torrent.Size,
			PercentProgress: torrent.PercentProgress,
			Downloaded:      torrent.Downloaded,
			Uploaded:        torrent.Uploaded,
			Ratio:           torrent.Ratio,
			UploadSpeed:     torrent.UploadSpeed,
			DownloadSpeed:   torrent.DownloadSpeed,
			ETA:             torrent.ETA,
			Label:           torrent.Label,
			PeersConnected:  torrent.PeersConnected,
			PeersTotal:      torrent.PeersTotal,
			SeedsConnected:  torrent.SeedsConnected,
			SeedsTotal:      torrent.SeedsTotal,
			Availability:    torrent.Availability,
			QueueOrder:      torrent.QueueOrder,
			Remaining:       max(torrent.Size-torrent.Downloaded, 0),
			Status:          torrent.Status,
			AddedOn:         int(torrent.AddedRaw),
			CompletedOn:     int(torrent.AddedRaw), // Not Provided
			FilePath:        torrent.FilePath + "/" + torrent.Name,
		})
	}
	return nil
}

// GetTorrents returns a list of Torrent structs containing all of the torrents
// added to the deluge/Bittorrent server
func (c *Client) GetTorrents() ([]Torrent, error) {
	var torrents TorrentsResponse
	err := c.action("core.get_torrents_status", fmt.Sprintf("{},[%s]",
		TorrentProperties), &torrents)
	if err != nil {
		return nil, fmt.Errorf("Error getting torrents: %s", err.Error())
	}

	return torrents.Torrents, nil
}

// GetTorrent gets a specific torrent by info hash
func (c *Client) GetTorrent(hash string) (Torrent, error) {
	var torrent TorrentResponse
	err := c.action("core.get_torrent_status", fmt.Sprintf("\"%s\",[%s]", hash,
		TorrentProperties), &torrent)
	if err != nil {
		return Torrent{}, fmt.Errorf("Error getting torrents: %s", err.Error())
	}

	//FixUps
	torrent.Torrent.StatusCode = 200
	torrent.Torrent.AddedOn = int(torrent.Torrent.AddedRaw)
	torrent.Torrent.CompletedOn = int(torrent.Torrent.AddedRaw)
	torrent.Torrent.Remaining = max(torrent.Torrent.Size-torrent.Torrent.Downloaded, 0)

	return torrent.Torrent, nil
}

// PauseTorrent pauses the torrent specified by info hash
func (c *Client) PauseTorrent(hash string) error {
	var res BoolResponse
	err := c.action("core.pause_torrent", fmt.Sprintf("[\"%s\"]", hash), &res)
	if err != nil {
		return fmt.Errorf("Error pausing torrent: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error pausing torrent: %s", res.Error.Message)
	}

	return nil
}

// UnPauseTorrent unpauses the torrent specified by info hash
func (c *Client) UnPauseTorrent(hash string) error {
	var res BoolResponse
	err := c.action("core.resume_torrent", fmt.Sprintf("[\"%s\"]", hash), &res)
	if err != nil {
		return fmt.Errorf("Error resuming torrent: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error resuming torrent: %s", res.Error.Message)
	}

	return nil
}

// StartTorrent starts the torrent specified by info hash
func (c *Client) StartTorrent(hash string) error {
	//Deluge has no concept of "Start/Stop" so mimicing using UnPause
	return c.UnPauseTorrent(hash)
}

// StopTorrent stops the torrent specified by info hash
func (c *Client) StopTorrent(hash string) error {
	//Deluge has no concept of "Start/Stop" so mimicing using Pause
	return c.PauseTorrent(hash)
}

// RecheckTorrent rechecks the torrent specified by info hash
func (c *Client) RecheckTorrent(hash string) error {
	var res BoolResponse
	err := c.action("core.force_recheck", fmt.Sprintf("[\"%s\"]", hash), &res)
	if err != nil {
		return fmt.Errorf("Error rechecking torrent: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error rechecking torrent: %s", res.Error.Message)
	}

	return nil
}

// RemoveTorrent removes the torrent specified by info hash
func (c *Client) RemoveTorrent(hash string) error {
	var res BoolResponse
	err := c.action("core.remove_torrent", fmt.Sprintf("\"%s\", false", hash), &res)
	if err != nil {
		return fmt.Errorf("Error removing torrent: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error removing torrent: %s", res.Error.Message)
	}

	return nil
}

// RemoveTorrentAndData removes the torrent and associated data specified by info hash
func (c *Client) RemoveTorrentAndData(hash string) error {
	var res BoolResponse
	err := c.action("core.remove_torrent", fmt.Sprintf("\"%s\", true", hash), &res)
	if err != nil {
		return fmt.Errorf("Error removing torrent: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error removing torrent: %s", res.Error.Message)
	}

	return nil
}

// AddTorrent adds the torrent specified by url or magnet link
func (c *Client) AddTorrent(url string) error {
	var res StringResponse
	err := c.action("core.add_torrent_magnet", fmt.Sprintf("\"%s\",{}", url), &res)
	if err != nil {
		return fmt.Errorf("Error adding torrent: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error adding torrent: %s", res.Error.Message)
	}

	return nil
}

// AddTorrentFile adds the torrent specified by a file on disk
func (c *Client) AddTorrentFile(torrentpath string) error {
	f, err := os.Open(torrentpath)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Error opening torrent file: %s", err.Error())
	}
	blob, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		return fmt.Errorf("Error reading torrent file: %s", err.Error())
	}

	var res StringResponse
	err = c.action("core.add_torrent_file", fmt.Sprintf("\"%s\", \"%s\",{}",
		filepath.Base(torrentpath), base64.StdEncoding.EncodeToString(blob)), &res)
	if err != nil {
		return fmt.Errorf("Error adding torrent: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error adding torrent: %s", res.Error.Message)
	}

	return nil
}

// SetTorrentProperty sets a property for the given torrent.
func (c *Client) SetTorrentProperty(hash string, property string, value string) error {
	err := errors.New("TODO")
	if err != nil {
		return fmt.Errorf("Error setting torrent (%s) '%s' to '%s': %s ", hash,
			property, value, err)
	}

	return nil
}

// SetTorrentLabel sets the label for the given torrent
func (c *Client) SetTorrentLabel(hash string, label string) error {
	// err := c.SetTorrentProperty(hash, "label", label)
	err := errors.New("TODO")
	if err != nil {
		return err
	}

	return nil
}

// SetTorrentSeedRatio sets the seed ratio for the given torrent
func (c *Client) SetTorrentSeedRatio(hash string, ratio float64) error {
	// err := c.SetTorrentProperty(hash, "seed_override", "1")
	err := errors.New("TODO")
	if err != nil {
		return err
	}

	// err = c.SetTorrentProperty(hash, "seed_ratio", strconv.FormatFloat(ratio*10, 'f', 0, 64))
	if err != nil {
		return err
	}

	return nil
}

// SetTorrentSeedTime sets the seed time for the given torrent
func (c *Client) SetTorrentSeedTime(hash string, time int) error {
	//Deluge does not have a concept of stop after seeding for a specific time
	err := errors.New("Not Implemented")
	if err != nil {
		return err
	}

	return nil
}

// QueueTop sends the torrent to the top of the download queue
func (c *Client) QueueTop(hash string) error {
	var res BoolResponse
	err := c.action("core.queue_top", fmt.Sprintf("[\"%s\"]", hash), &res)
	if err != nil {
		return fmt.Errorf("Error setting torrent queue priority: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error setting torrent queue priority: %s", res.Error.Message)
	}

	return nil
}

// QueueUp moves the torrent up the download queue
func (c *Client) QueueUp(hash string) error {
	var res BoolResponse
	err := c.action("core.queue_up", fmt.Sprintf("[\"%s\"]", hash), &res)
	if err != nil {
		return fmt.Errorf("Error setting torrent queue priority: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error setting torrent queue priority: %s", res.Error.Message)
	}

	return nil
}

// QueueUp moves the torrent down the download queue
func (c *Client) QueueDown(hash string) error {
	var res BoolResponse
	err := c.action("core.queue_down", fmt.Sprintf("[\"%s\"]", hash), &res)
	if err != nil {
		return fmt.Errorf("Error setting torrent queue priority: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error setting torrent queue priority: %s", res.Error.Message)
	}

	return nil
}

// QueueTop sends the torrent to the bottom of the download queue
func (c *Client) QueueBottom(hash string) error {
	var res BoolResponse
	err := c.action("core.queue_bottom", fmt.Sprintf("[\"%s\"]", hash), &res)
	if err != nil {
		return fmt.Errorf("Error setting torrent queue priority: %s", err.Error())
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error setting torrent queue priority: %s", res.Error.Message)
	}

	return nil
}
