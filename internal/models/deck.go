package models

import (
	"crypto/rand"
	"errors"
	"math"
)

type Deck[T any] struct {
	cards []T
}

func NewDeck[T any](cards []T) Deck[T] {
	return Deck[T]{cards: cards}
}

func (d *Deck[T]) Len() uint {
	return uint(len(d.cards))
}

func (d *Deck[T]) Insert(t T) {
	d.cards = append(d.cards, t)
}

func (d *Deck[T]) isEmpty() bool {
	return len(d.cards) == 0
}

func (d *Deck[T]) Draw() (T, error) {
	var card T

	if d.isEmpty() {
		return card, errors.New("deck is empty")
	}

	randomItem, err := randomUint(int(d.Len() - 1))
	if err != nil {
		return card, err
	}

	return d.cards[randomItem], nil
}

func randomUint(max int) (int, error) {
	if max == 0 {
		return 0, nil
	}

	bytesNeeded := int(math.Log2(float64(max))/8) + 1
	if bytesNeeded > 8 {
		bytesNeeded = 8
	}

	b := make([]byte, bytesNeeded)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}

	var random int
	for i := 0; i < bytesNeeded; i++ {
		random = random<<8 | int(b[i])
	}

	return random % (max + 1), nil
}
