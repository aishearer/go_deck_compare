package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
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
	priceInCommon float32
}

type ComparisonResults []ComparisonResult

func (slice ComparisonResults) Len() int {
	return len(slice)
}

func (slice ComparisonResults) Less(i, j int) bool {
	return slice[i].priceInCommon < slice[j].priceInCommon
}

func (slice ComparisonResults) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type CardPriceDatabase map[string]float32

func (cardDb CardPriceDatabase) load(path string) error {
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

func (cardDb CardPriceDatabase) getPrice(card string) float32 {
	price, exists := cardDb[card]
	if !exists {
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
		//TODO something
		fmt.Println("File error")
		return deckMap, err
	}

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

func compareDecks(deck1 map[string]int, deck2 map[string]int, cardPriceDb CardPriceDatabase) float32 {
	var priceInCommon float32 = 0

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

func main() {
	priceDb := CardPriceDatabase{}
	err := priceDb.load("EveryCard_lithiumDucky_2017.April.20.csv")
	if err != nil {
		fmt.Println("Error loading card price file: " + err.Error())
		return
	}

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
			priceInCommon := compareDecks(deck1.cardMap, deck2.cardMap, priceDb)
			results = append(results, ComparisonResult{deck1.name, deck2.name, priceInCommon})
		}
	}

	sort.Sort(results)
	for _, result := range results {
		fmt.Printf("%f %s %s\n", result.priceInCommon, result.deck1, result.deck2)
	}

}
