package deck

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

type DeckManager struct {
	Hand []string
	Deck []string
}

func InitDeck() *DeckManager {
	var newDeck DeckManager
	newDeck.Hand = []string{}
	newDeck.Deck = []string{
		"The Fool(0)",
		"The Magician(I)",
		"The High Priestess(II)",
		"The Empress(III)",
		"The Emporer(IV)",
		"The Heirophant(V)",
		"The Lovers(VI)",
		"The Chariot(VII)",
		"Strength(VIII)",
		"The Hermit(IX)",
		"Wheel of Fortune(X)",
		"Justice(XI)",
		"The Hanged Man(XII)",
		"Death(XIII)",
		"Temperance(XIV)",
		"The Devil(XV)",
		"The Tower(XVI)",
		"The Star(XVII)",
		"The Moon(XVIII)",
		"The Sun(XIX)",
		"Judgement(XX)",
		"The World(XXI)",
		"King of Swords",
		"Queen of Swords",
		"Knight of Swords",
		"Page of Swords",
		"One of Swords",
		"Two of Swords",
		"Three of Swords",
		"Four of Swords",
		"Five of Swords",
		"Six of Swords",
		"Seven of Swords",
		"Eight of Swords",
		"Nine of Swords",
		"Ten of Swords",
		"King of Batons",
		"Queen of Batons",
		"Knight of Batons",
		"Page of Batons",
		"One of Batons",
		"Two of Batons",
		"Three of Batons",
		"Four of Batons",
		"Five of Batons",
		"Six of Batons",
		"Seven of Batons",
		"Eight of Batons",
		"Nine of Batons",
		"Ten of Batons",
		"King of Coins",
		"Queen of Coins",
		"Knight of Coins",
		"Page of Coins",
		"One of Coins",
		"Two of Coins",
		"Three of Coins",
		"Four of Coins",
		"Five of Coins",
		"Six of Coins",
		"Seven of Coins",
		"Eight of Coins",
		"Nine of Coins",
		"Ten of Coins",
		"King of Cups",
		"Queen of Cups",
		"Knight of Cups",
		"Page of Cups",
		"One of Cups",
		"Two of Cups",
		"Three of Cups",
		"Four of Cups",
		"Five of Cups",
		"Six of Cups",
		"Seven of Cups",
		"Eight of Cups",
		"Nine of Cups",
		"Ten of Cups",
	}
	newDeck.TarotShuffle()

	//Rust suggests
	//78uuuuuuuuuuuuuuuuu

	return &newDeck
}
func (m *DeckManager) TarotShuffle() {
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := len(m.Deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		m.Deck[i], m.Deck[j] = m.Deck[j], m.Deck[i]
	}
}
func (m *DeckManager) Draw() string {
	var inversion string
	if rand.Intn(2) == 1 {
		inversion = " inverted"
	}

	m.Hand = append(m.Hand, m.Deck[len(m.Deck)-1])
	m.Deck = m.Deck[:len(m.Deck)-1]

	fmt.Println(m.Hand)

	return m.Hand[len(m.Hand)-1] + inversion
}
func (m *DeckManager) ResetDeck() {
	m.Deck = append(m.Hand, m.Deck...)
	m.Hand = []string{}
	m.TarotShuffle()
}
