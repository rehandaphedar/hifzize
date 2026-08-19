package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"

	"github.com/npcnixel/genanki-go"
)

const (
	minDeckID int64 = 1_000_000_000
	maxDeckID int64 = 9_999_999_999
)

func loadJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return nil
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	return string(data), nil
}

func getDeck(decks map[int64]*genanki.Deck, id int64, name string, description string) *genanki.Deck {
	if deck, ok := decks[id]; ok {
		return deck
	}

	deck := genanki.NewDeck(id, name, description)
	decks[id] = deck
	return deck
}

func compareDecks(a, b *genanki.Deck) int {
	return cmp.Compare(a.Name, b.Name)
}

func deckID(id int64) int64 {
	return minDeckID + id%(maxDeckID-minDeckID+1)
}
