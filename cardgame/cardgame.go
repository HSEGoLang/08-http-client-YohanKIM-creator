package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://deckofcardsapi.com/api/deck"

type NewDeckResponse struct {
	Success   bool   `json:"success"`
	DeckID    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

type Card struct {
	Code  string `json:"code"`
	Image string `json:"image"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
}

type DrawResponse struct {
	Success   bool   `json:"success"`
	DeckID    string `json:"deck_id"`
	Cards     []Card `json:"cards"`
	Remaining int    `json:"remaining"`
}

func main() {
	// проверка аргумента
	if len(os.Args) < 2 {
		os.Exit(2)
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n < 1 {
		os.Exit(2)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	newDeckURL := fmt.Sprintf("%s/new/shuffle/?deck_count=1", baseURL)
	deckResp, err := client.Get(newDeckURL)
	if err != nil {
		os.Exit(1)
	}
	defer deckResp.Body.Close()

	var deck NewDeckResponse
	if err := json.NewDecoder(deckResp.Body).Decode(&deck); err != nil || !deck.Success {
		os.Exit(1)
	}

	firstQueenPos := 0
	pos := 0
	remaining := deck.Remaining

	for remaining > 0 {
		drawURL := fmt.Sprintf("%s/%s/draw/?count=1", baseURL, deck.DeckID)
		dr, err := client.Get(drawURL)
		if err != nil {
			os.Exit(1)
		}

		var draw DrawResponse
		if err := json.NewDecoder(dr.Body).Decode(&draw); err != nil || !draw.Success || len(draw.Cards) != 1 {
			dr.Body.Close()
			os.Exit(1)
		}
		dr.Body.Close()

		pos++
		if draw.Cards[0].Value == "QUEEN" && firstQueenPos == 0 {
			firstQueenPos = pos
			break
		}
		remaining = draw.Remaining
	}

	if firstQueenPos == n {
		fmt.Print("You win")
	} else {
		fmt.Print("You lose")
	}
}
