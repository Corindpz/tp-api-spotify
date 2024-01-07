package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var accessToken = ""
var clientID     = "d185fdc02ecb4374bc31fa44146eaa6f"
var clientSecret = "2ddcf2a0d2324e2991105b64f51fabdc"


type datasSon struct {
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Album struct {
		ReleaseDate string `json:"release_date"`
		Name        string `json:"name"`
		Images      []struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"album"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

type datasAlbum struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	ReleaseDate string `json:"release_date"`
	TotalTracks int    `json:"total_tracks"`
}

func infosdm(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	soundID := "0EzNyXyU7gHzj2TN8qYThj"

	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", soundID)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	var datasSon datasSon
	json.NewDecoder(resp.Body).Decode(&datasSon)

	tmpl.ExecuteTemplate(w, "sdm", datasSon)
}

func infojul(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	artistID := "3IW7ScrzXmPvZhB27hmfgy"

	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums", artistID)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	var albums struct {
		Items []datasAlbum `json:"items"`
	}

	json.NewDecoder(resp.Body).Decode(&albums)

	tmpl.ExecuteTemplate(w, "jul", albums.Items)
}

func tokenaccess() (string) {
	clientCreds := fmt.Sprintf("%s:%s", clientID, clientSecret)
	clientCredsB64 := base64.StdEncoding.EncodeToString([]byte(clientCreds))

	tokenURL := "https://accounts.spotify.com/api/token"

	tokenData := strings.NewReader("grant_type=client_credentials")

	tokenHeaders := map[string]string{
		"Authorization": "Basic " + clientCredsB64,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	req, _ := http.NewRequest("POST", tokenURL, tokenData)


	for key, value := range tokenHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var tokenResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tokenResponse)

		accessToken, _ := tokenResponse["access_token"].(string)
		return accessToken
	}

	return ""
}

func main() {

	accessToken = tokenaccess()

	css := http.FileServer(http.Dir("./styles"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index", nil)
	})

	http.HandleFunc("/album/jul", infojul)

	http.HandleFunc("/track/sdm", infosdm)

	fmt.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
