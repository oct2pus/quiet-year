package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type card struct {
	Suit    int    `json:"Suit"` //0 spring/heart, 1 summer/diamond, 2 fall/club, 3 winter/spade
	Face    string `json:"Face"`
	Text    string `json:"Text"`    //can be empty
	Option1 string `json:"Option1"` //can be empty
	Option2 string `json:"Option2"` //can be empty
}

func (c card) Info() string {
	return fmt.Sprintf(
		"Suit: %v | Face: %v",
		c.Suit, c.Face)
}

type deck struct {
	Cards []card `json:"Cards"`
}

// Shuffle unsorts the cards.
func (d *deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		a := rand.Intn(len(d.Cards))
		b := rand.Intn(len(d.Cards))
		c := d.Cards[a]
		d.Cards[a] = d.Cards[b]
		d.Cards[b] = c
	}
}

// Info prints debug info.
func (d *deck) Info() string {
	cardInfo := make([]string, 0)
	for _, ele := range d.Cards {
		cardInfo = append(cardInfo, ele.Info())
	}
	return fmt.Sprintf("Length: %v\n Cards: %v\n", len(d.Cards), cardInfo)
}

// Draw returns a card and removes it from the deck.
func (d *deck) Draw() (card, error) {
	if len(d.Cards) == 0 {
		return card{}, errors.New("no cards available")
	}
	out := d.Cards[len(d.Cards)-1]
	newCards := make([]card, len(d.Cards)-1)
	copy(newCards, d.Cards[0:len(d.Cards)-1])
	d.Cards = newCards
	return out, nil
}

func newDeck() deck {
	var d deck
	d.Cards = make([]card, 0)
	return d
}

func PrepareDecks() *([]deck) {
	file, err := os.ReadFile("./cards.json")
	if err != nil {
		log.Printf("%v\n", err)
	}
	var pile *deck
	json.Unmarshal(file, &pile)
	d := make([]deck, 0, 4)
	d = append(d, newDeck())
	d = append(d, newDeck())
	d = append(d, newDeck())
	d = append(d, newDeck())
	for _, ele := range pile.Cards {
		d[ele.Suit].Cards = append(d[ele.Suit].Cards, ele)
	}
	return &d
}
