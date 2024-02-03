package api

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Client is in charge of interacting with a proxmox server
type Client struct {
	APITokenID     string
	APITokenSecret string
	BaseURL        string
	HTTPClient     *http.Client
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func NewClient(baseurl string, tokenid string, tokensecret string) *Client {
	return &Client{
		APITokenID:     tokenid,
		APITokenSecret: tokensecret,
		BaseURL:        baseurl,
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: tr,
		},
	}
}
