package reacts

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const (
	klipyAPIKey    = "lJodCdaswwEpTPg5Lhix66aIcaFsXBTrKKGlAzX1rPSQvagQHDNczRi42lNJ6x56"
	klipyClientKey = "tronimoes-js"
	klipyLimit     = 30
)

// klipySearchResponse matches the KLIPY v2 search API response (only the fields we need).
type klipySearchResponse struct {
	Results []klipyResult `json:"results"`
}

type klipyResult struct {
	MediaFormats struct {
		TinygifTransparent struct {
			URL string `json:"url"`
		} `json:"tinygif_transparent"`
	} `json:"media_formats"`
}

// no more than one klipy query every 5 seconds.
var rateLimit = time.NewTicker(5 * time.Second)

// FindImageURL calls the KLIPY API with the given query and returns a random URL
// from the first 30 results (tinygif_transparent), or empty string on error or no results.
func FindImageURL(ctx context.Context, query string) (string, error) {
	select {
	case <-rateLimit.C:
	default:
		return "", fmt.Errorf("rate limit exceeded")
	}
	u := fmt.Sprintf(
		"https://api.klipy.com/v2/search?q=%s&key=%s&client_key=%s&limit=%d&searchfilter=sticker",
		url.QueryEscape(query),
		url.QueryEscape(klipyAPIKey),
		url.QueryEscape(klipyClientKey),
		klipyLimit,
	)
	resp, err := http.Get(u)
	if err != nil {
		return "", fmt.Errorf("error getting image URL: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}
	var data klipySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}
	if len(data.Results) == 0 {
		return "", fmt.Errorf("no results")
	}
	idx := rand.Intn(len(data.Results))
	gifURL := data.Results[idx].MediaFormats.TinygifTransparent.URL
	if gifURL == "" {
		return "", fmt.Errorf("no gif URL")
	}
	return gifURL, nil
}
