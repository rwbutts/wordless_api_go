package main

import (
	//	"fmt"
	"encoding/json"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/rwbutts/wordless_api_go/words"
)

type healthCheckResponse struct {
	Alive bool `json:"alive"`
}

type getWordResponse struct {
	Word string `json:"word"`
}

type checkWordResponse struct {
	Exists bool `json:"exists"`
}

type getMatchCountResponse struct {
	Count int `json:"count"`
}

type queryMatchCountRequest struct {
	Answer  string   `json:"answer"`
	Guesses []string `json:"guesses"`
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sendJSON(w, healthCheckResponse{Alive: true})
}

func handleGetRandomWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	word := words.WordList[rand.Intn(len(words.WordList))]

	sendJSON(w, getWordResponse{Word: word})
}

var getWordURLRegexp = regexp.MustCompile(`^/api/getword/(-?\d+)$`)

func handleGetWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	match := getWordURLRegexp.FindStringSubmatch(r.URL.Path)
	if match == nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	//daysAgoParam := match[1]
	word := words.WordList[rand.Intn(len(words.WordList))]

	sendJSON(w, getWordResponse{Word: word})
}

var checkWordURLRegexp = regexp.MustCompile(`^/api/checkword/(\w+)$`)

func handleCheckWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	match := checkWordURLRegexp.FindStringSubmatch(r.URL.Path)
	if match == nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	word := strings.ToLower(match[1])
	index := sort.SearchStrings(words.WordList, word)
	exists := index < len(words.WordList) && words.WordList[index] == word
	sendJSON(w, checkWordResponse{Exists: exists})
}

func handleMatchCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params queryMatchCountRequest

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Malformed request JSON", http.StatusBadRequest)
		return
	}

	matchCount := words.CountMatches(params.Answer, params.Guesses)

	sendJSON(w, getMatchCountResponse{Count: matchCount})
}

func sendJSON(w http.ResponseWriter, payload any) {
	jsonResp, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
