package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// type PageViewParams struct {
// 	accept string
// 	owner  string
// 	repo   string
// 	per    string
// }

type PageView struct {
	Count   int32 `json:"count"`
	Uniques int32 `json:"uniques"`
	Views   []struct {
		TimeStamp string `json:"timestamp"`
		Count     int32  `json:"count"`
		Uniques   int32  `json:"uniques"`
	}
}

type Referrers []struct {
	Referrer string `json:"referrer"`
	Count    int32  `json:"count"`
	Uniques  int32  `json:"uniques"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	getPageViews()

}

func apiRequest(method string, url string) []byte {
	// Handles all requests to the GitHub API
	// Requires 2 params: HTTP method and the url for the endpoint
	var client http.Client
	// Get values for auth from .env using godotenv
	authToken := os.Getenv("AUTH")
	username := os.Getenv("USERNAME")

	req, err := http.NewRequest(method, "https://api.github.com/"+url, nil)
	if err != nil {
		panic(err)
	}
	// Use generated token with basic auth for authentication
	req.SetBasicAuth(username, authToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// Close the body when we are done with it
	defer resp.Body.Close()

	// Read the body data and return it
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body

}

func getPageViews() {
	// Accesses the metrics endpoint to get page views
	// Get the total number of views and breakdown per day or week for the last 14 days
	// Per: Day(default), Month
	url := "repos/aerovize/H1-CLI/traffic/views"

	ghJSON := apiRequest("GET", url)

	// Parse the json response into the PageViews struct
	data := PageView{}
	err := json.Unmarshal(ghJSON, &data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data.Views)
}

func getReferrs() {
	url := "repos/aerovize/H1-CLI/traffic/popular/referrers"

	// Accesses the metrics endpoint to get the top referral sources
	// Get the top 10 referrers over the last 14 days.
	ghJSON := apiRequest("GET", url)

	// Parse the json response into the Refferrer struct
	data := Referrers{}
	err := json.Unmarshal(ghJSON, &data)
	if err != nil {
		panic(err)
	}
}

// TODO: Push to github repo for storage & server to show page with the data or CLI
