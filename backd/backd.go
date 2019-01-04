package backd

import (
	"net"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Backd is the struct that holds the client for the service
type Backd struct {
	sling      *sling.Sling
	authURL    string
	objectsURL string
	apiKey     string
	sessionID  string
	expiresAt  int64
}

// NewClient returns an usable client to connect to an instance of Backd
func NewClient(authURL, objectsURL, apiKey string) *Backd {

	var (
		backd Backd
	)

	backd.authURL = authURL
	backd.objectsURL = objectsURL
	backd.apiKey = apiKey
	backd.ConnectionTimeouts(5, 5, 10)

	return &backd

}

// ConnectionTimeouts allow to change the client timeouts for:
//  - Dialer
//  - TLS Handshake
//  - HTTP timeout
func (b *Backd) ConnectionTimeouts(dialer, tlsHandshake, timeout time.Duration) {

	b.sling = sling.New().Client(&http.Client{
		Timeout: timeout * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: dialer * time.Second,
			}).Dial,
			TLSHandshakeTimeout: tlsHandshake * time.Second,
		},
	}).Set("User-Agent", clientName)

}
