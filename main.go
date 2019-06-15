package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func search(query string) {
	url := "https://itunes.apple.com/search?term=%s"

	client := new(http.Client)

	req, err := http.NewRequest("GET", fmt.Sprint(url, query), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("An error happened while creating the request.  The request has not been sent.")
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("An error happened while executing the request. The request could have been sent.")
		return
	}

	log.WithFields(log.Fields{
		"response": resp.Body,
		"code":     resp.StatusCode,
	}).Info("Search finished.")
}

func main() {
	// Setup Logrus
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		search("lady+gaga")
		fmt.Fprintf(w, "Ok")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
