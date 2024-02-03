package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type LXCStatusResponse struct {
	Data struct {
		Uptime int     `json:"uptime"`
		Status string  `json:"status"`
		CPU    float64 `json:"cpu"`
		Swap   int     `json:"swap"`
		Maxmem int     `json:"maxmem"`
		Tags   string  `json:"tags"`
		Pid    int     `json:"pid"`
		Ha     struct {
			Managed int `json:"managed"`
		} `json:"ha"`
		Type      string `json:"type"`
		Disk      int    `json:"disk"`
		Maxswap   int    `json:"maxswap"`
		Name      string `json:"name"`
		Cpus      int    `json:"cpus"`
		Diskwrite int    `json:"diskwrite"`
		Netout    int    `json:"netout"`
		Diskread  int64  `json:"diskread"`
		Netin     int    `json:"netin"`
		Maxdisk   int64  `json:"maxdisk"`
		Vmid      int    `json:"vmid"`
		Mem       int    `json:"mem"`
	} `json:"data"`
}

// Get Status of lxc id given
func (client *Client) Status(node string, lxc string) (LXCStatusResponse, error) {
	// Get current status
	APIStatusURL := fmt.Sprintf("/nodes/%s/lxc/%s/status/current", node, lxc)

	// Craft new URL
	NewURL, err := url.JoinPath(client.BaseURL, APIStatusURL)
	req, _ := http.NewRequest("GET", NewURL, nil)

	// Set Authorization Header
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", client.APITokenID, client.APITokenSecret))

	// Send request
	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return LXCStatusResponse{}, err
	}

	defer res.Body.Close()

	var post LXCStatusResponse

	// Decode the JSON reponse
	if err := json.NewDecoder(res.Body).Decode(&post); err != nil {
		fmt.Print(res.Body)
		log.Fatalf("An error occured while decode the JSON : %s", err)
	}

	return post, nil
}

func (client *Client) Start(node string, lxc string) error {
	// Start de LXC container
	APIStopURL := fmt.Sprintf("/nodes/%s/lxc/%s/status/start", node, lxc)

	// Craft new URL
	NewURL, err := url.JoinPath(client.BaseURL, APIStopURL)
	req, _ := http.NewRequest("POST", NewURL, nil)

	// Set Authorization Header
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", client.APITokenID, client.APITokenSecret))

	// Send request
	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

// Shutdown the LXC container given in parameters
func (client *Client) Shutdown(node string, lxc string) error {
	// Stop the LXC container
	APIStopURL := fmt.Sprintf("/nodes/%s/lxc/%s/status/shutdown", node, lxc)

	// Craft new URL
	NewURL, err := url.JoinPath(client.BaseURL, APIStopURL)
	req, _ := http.NewRequest("POST", NewURL, nil)

	// Set Authorization Header
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", client.APITokenID, client.APITokenSecret))

	// Send request
	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
