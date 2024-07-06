package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

var urlDB = make(map[string]URL)

func GenerateShortURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL)) // It converts the originalURL string to a byte slice
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func CreateURL(originalURL string) string {
	shortUrl := GenerateShortURL(originalURL)
	id := shortUrl
	urlDB[id] = URL{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    shortUrl,
		CreatedAt:   time.Now(),
	}
	return shortUrl
}

func GetURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to URL shortener!")
}

func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body!", http.StatusBadRequest)
		return
	}
	shortURL := CreateURL(request.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := GetURL(id)
	if err != nil {
		http.Error(w, "Invalid request!", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	// Register the handler function to handle all requests to the root URL
	http.HandleFunc("/", Handler)
	http.HandleFunc("/shorten", ShortUrlHandler)
	http.HandleFunc("/redirect/", RedirectUrlHandler)

	// Start the HTTP server
	fmt.Println("Starting server on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error on starting server: ", err)
	}
}
