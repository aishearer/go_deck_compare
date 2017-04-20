package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type Deck struct {
	name    string
	cardMap map[string]int
}

type ComparisonResult struct {
	deck1, deck2  string
	cardsInCommon int
}

type ComparisonResults []ComparisonResult

func (slice ComparisonResults) Len() int {
	return len(slice)
}

func (slice ComparisonResults) Less(i, j int) bool {
	return slice[i].cardsInCommon < slice[j].cardsInCommon
}

func (slice ComparisonResults) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func readDecklist(path string) (map[string]int, error) {
	deckMap := make(map[string]int)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		//TODO something
		fmt.Println("File error")
		return deckMap, err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		split := strings.SplitN(line, " ", 2)
		//TODO check split length
		cards, err := strconv.Atoi(split[0])
		if err != nil {
			//TODO
			fmt.Println("Error converting card number")
			return deckMap, err
		}
		cardName := split[1]
		deckMap[cardName] = cards
	}

	return deckMap, nil
}

func compareDecks(deck1 map[string]int, deck2 map[string]int) int {
	numInCommon := 0

	for cardName, deck1Number := range deck1 {
		deck2Number, exists := deck2[cardName]
		if exists {
			numInCommon += min(deck1Number, deck2Number)
		}
	}

	return numInCommon
}

func main() {
	decks := []Deck{}
	files, _ := ioutil.ReadDir("./decks")
	for _, f := range files {
		fileName := f.Name()
		cardMap, _ := readDecklist("decks/" + fileName)
		decks = append(decks, Deck{fileName, cardMap})
	}

	results := ComparisonResults{}
	// For every combination of decks
	for i, deck1 := range decks {
		for _, deck2 := range decks[i+1:] {
			numInCommon := compareDecks(deck1.cardMap, deck2.cardMap)
			results = append(results, ComparisonResult{deck1.name, deck2.name, numInCommon})
		}
	}

	sort.Sort(results)
	for _, result := range results {
		fmt.Printf("%d %s %s\n", result.cardsInCommon, result.deck1, result.deck2)
	}

}
