package reacts

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

const (
	klipyAPIKey   = "lJodCdaswwEpTPg5Lhix66aIcaFsXBTrKKGlAzX1rPSQvagQHDNczRi42lNJ6x56"
	klipyClientKey = "tronimoes-js"
	klipyLimit    = 30
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

// FindImageURL calls the KLIPY API with the given query and returns a random URL
// from the first 30 results (tinygif_transparent), or empty string on error or no results.
func FindImageURL(query string) string {
	u := fmt.Sprintf(
		"https://api.klipy.com/v2/search?q=%s&key=%s&client_key=%s&limit=%d&searchfilter=sticker",
		url.QueryEscape(query),
		url.QueryEscape(klipyAPIKey),
		url.QueryEscape(klipyClientKey),
		klipyLimit,
	)
	resp, err := http.Get(u)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	var data klipySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ""
	}
	if len(data.Results) == 0 {
		return ""
	}
	idx := rand.Intn(len(data.Results))
	gifURL := data.Results[idx].MediaFormats.TinygifTransparent.URL
	if gifURL == "" {
		return ""
	}
	return gifURL
}
