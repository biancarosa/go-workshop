package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type SearchResponse struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		WrapperType            string    `json:"wrapperType"`
		Kind                   string    `json:"kind"`
		ArtistID               int       `json:"artistId"`
		CollectionID           int       `json:"collectionId"`
		TrackID                int       `json:"trackId"`
		ArtistName             string    `json:"artistName"`
		CollectionName         string    `json:"collectionName"`
		TrackName              string    `json:"trackName"`
		CollectionCensoredName string    `json:"collectionCensoredName"`
		TrackCensoredName      string    `json:"trackCensoredName"`
		ArtistViewURL          string    `json:"artistViewUrl"`
		CollectionViewURL      string    `json:"collectionViewUrl"`
		TrackViewURL           string    `json:"trackViewUrl"`
		PreviewURL             string    `json:"previewUrl"`
		ArtworkURL30           string    `json:"artworkUrl30"`
		ArtworkURL60           string    `json:"artworkUrl60"`
		ArtworkURL100          string    `json:"artworkUrl100"`
		CollectionPrice        float64   `json:"collectionPrice"`
		TrackPrice             float64   `json:"trackPrice"`
		ReleaseDate            time.Time `json:"releaseDate"`
		CollectionExplicitness string    `json:"collectionExplicitness"`
		TrackExplicitness      string    `json:"trackExplicitness"`
		DiscCount              int       `json:"discCount"`
		DiscNumber             int       `json:"discNumber"`
		TrackCount             int       `json:"trackCount"`
		TrackNumber            int       `json:"trackNumber"`
		TrackTimeMillis        int       `json:"trackTimeMillis"`
		Country                string    `json:"country"`
		Currency               string    `json:"currency"`
		PrimaryGenreName       string    `json:"primaryGenreName"`
		IsStreamable           bool      `json:"isStreamable"`
		CollectionArtistName   string    `json:"collectionArtistName,omitempty"`
		CollectionArtistID     int       `json:"collectionArtistId,omitempty"`
	} `json:"results"`
}

func search(query string) (*SearchResponse, error) {
	uri := "https://itunes.apple.com/search?term=%s"

	client := new(http.Client)
	req, err := http.NewRequest("GET", fmt.Sprintf(uri, url.QueryEscape(query)), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("An error happened while creating the request.  The request has not been sent.")
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("An error happened while executing the request. The request could have been sent.")
		return nil, err
	}

	log.WithFields(log.Fields{
		"code": resp.StatusCode,
	}).Info("Search finished.")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("An error happened while reading the request.  The request has not been sent.")
		return nil, err
	}

	var s SearchResponse
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("An error happened while unmarshalling the request.  The request has not been sent.")
		return nil, err
	}
	return &s, err
}

func main() {
	// Setup Logrus
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		queries := r.URL.Query()["query"]
		list := make([]*SearchResponse, len(queries))
		for i, q := range queries {
			go func(i int) {
				s, err := search(q)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(struct {
						Error string `json:"error"`
					}{
						Error: err.Error(),
					})
					return
				}
				list[i] = s
			}(i)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
