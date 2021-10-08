package restful

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client client
type Client struct {
	*resty.Client
}

// New new client
func New(debug bool) *Client {
	client := resty.New()
	// Enable debug mode
	client.SetDebug(debug)
	// or One can disable security check (https)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// Set client timeout as per your need
	client.SetTimeout(1 * time.Minute)
	// Enabling Content length value for all request
	client.SetContentLength(true)

	return &Client{
		Client: client,
	}
}
