//TODO
// Use a curency package instead of float32

package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type deck struct {
	name    string
	cardMap map[string]int
	price   float32
}

type comparisonResult struct {
	deck1, deck2  string
	priceInCommon float32
}

type comparisonResults []comparisonResult

func (slice comparisonResults) Len() int {
	return len(slice)
}

func (slice comparisonResults) Less(i, j int) bool {
	return slice[i].priceInCommon < slice[j].priceInCommon
}

func (slice comparisonResults) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type cardPriceDatabase map[string]float32

func (cardDb cardPriceDatabase) load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	//Remove the csv file header
	lines = lines[1:]

	for _, line := range lines {
		price, err := strconv.ParseFloat(line[0], 32)
		if err != nil {
			return err
		}
		cardDb[line[1]] = float32(price)
	}

	return nil
}

func (cardDb cardPriceDatabase) getPrice(card string) float32 {
	price, exists := cardDb[card]
	if !exists {
		fmt.Printf("Can't find %s in price list\n", card)
		return 0.25
	}
	return price
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func readDecklist(path string) (map[string]int, error) {
	deckMap := make(map[string]int)

	lines, err := readLines(path)
	if err != nil {
		fmt.Printf("Error reading file %s\n", path)
		return deckMap, err
	}

	for _, line := range lines {
		if line == "" {
			continue
		}
		split := strings.SplitN(line, " ", 2)
		if len(split) != 2 {
			return deckMap, errors.New("Error parsing line: " + line)
		}
		cards, err := strconv.Atoi(split[0])
		if err != nil {
			fmt.Printf("On line: %s\n --- Error reading card number: %s\n", line, split[0])
			return deckMap, err
		}
		cardName := split[1]
		deckMap[cardName] = cards
	}

	return deckMap, nil
}

func compareDecks(deck1 map[string]int, deck2 map[string]int, cardPriceDb cardPriceDatabase) float32 {
	var priceInCommon float32

	for cardName, deck1Number := range deck1 {
		deck2Number, exists := deck2[cardName]
		if exists {
			numInCommon := min(deck1Number, deck2Number)
			cardPrice := cardPriceDb.getPrice(cardName)
			priceInCommon += (cardPrice * float32(numInCommon))
		}
	}

	return priceInCommon
}

func calculateDeckPrice(deck map[string]int, cardPriceDb cardPriceDatabase) float32 {
	var total float32
	for card, cardCount := range deck {
		total += float32(cardCount) * cardPriceDb.getPrice(card)
	}
	return total
}

func compareEveryDeckTogether(decks []deck, priceDb cardPriceDatabase) comparisonResults {
	results := comparisonResults{}
	// For every combination of decks
	for i, deck1 := range decks {
		for _, deck2 := range decks[i+1:] {
			priceInCommon := compareDecks(deck1.cardMap, deck2.cardMap, priceDb)
			results = append(results, comparisonResult{deck1.name, deck2.name, priceInCommon})
		}
	}
	return results
}

func ouputResults(results comparisonResults) {
	sort.Sort(results)
	for _, result := range results {
		fmt.Printf("%f %s %s\n", result.priceInCommon, result.deck1, result.deck2)
	}
}

func main() {
	var pathToCardPrice = flag.String("price-file", "card_prices.csv", "The path to the card price file")
	var deckDirectory = flag.String("deck-dir", "./decks/", "The path to the directory containing the deck files")
	var pathToCardCollection = flag.String("my-collection", "my_collection.dec", "The path to your card collection in .dec format")
	flag.Parse()

	priceDb := cardPriceDatabase{}
	err := priceDb.load(*pathToCardPrice)
	if err != nil {
		fmt.Println("Error loading card price file: " + err.Error())
		return
	}

	decks := []deck{}
	files, _ := ioutil.ReadDir(*deckDirectory)
	for _, f := range files {
		fileName := f.Name()
		cardMap, err := readDecklist(*deckDirectory + fileName)
		if err != nil {
			return
		}
		deckPrice := calculateDeckPrice(cardMap, priceDb)
		decks = append(decks, deck{fileName, cardMap, deckPrice})
	}

	cardCollection, err := readDecklist(*pathToCardCollection)
	if err != nil {
		return
	}

	results := comparisonResults{}
	for _, deck := range decks {
		priceInCommon := compareDecks(deck.cardMap, cardCollection, priceDb)
		results = append(results, comparisonResult{"-", deck.name + " " + strconv.FormatFloat(float64(deck.price), 'f', 2, 32), priceInCommon})
	}

	ouputResults(results)
}
