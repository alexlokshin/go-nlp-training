package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/euskadi31/go-tokenizer"
)

type KeyValue struct {
	Key   string
	Value int
}

func tokenize(tokenizer tokenizer.Tokenizer, text string) []string {
	tokens := tokenizer.Tokenize(strings.ToLower(text))
	return tokens
}

func collect(phrase string, phrases map[string]int) {
	value, present := phrases[phrase]
	if present {
		phrases[phrase] = value + 1
	} else {
		phrases[phrase] = 1
	}
}
func isIncluded(m map[string]int, text string) bool {
	_, v := m[text]
	return v
}
func generatePhrases(stopwords map[string]int, keywords []string, maxLength int) map[string]int {
	phrases := make(map[string]int)
	for i := 0; i < len(keywords); i++ {

		if !isIncluded(stopwords, keywords[i]) && !valid.IsInt(keywords[i]) {
			collect(keywords[i], phrases)
			if i < len(keywords)-1 {
				phrase := keywords[i]
				for j := 1; j < maxLength && j+i < len(keywords); j++ {
					if !isIncluded(stopwords, keywords[i+j]) {
						phrase = phrase + " " + keywords[i+j]
						collect(phrase, phrases)
					}
				}
			}
		}
	}
	return phrases
}
func accumulate(phrases map[string]int, temp map[string]int) {
	for k, v := range temp {
		value, present := phrases[k]
		if present {
			phrases[k] = value + v
		} else {
			phrases[k] = v
		}
	}
}

func readKnownPhrases(fileName string) map[string]int {
	phrases := make(map[string]int)
	file, err := os.Open(fileName)
	if err != nil {
		return phrases
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		text := scanner.Text()
		parts := strings.Split(text, "|")
		if len(text) > 0 {
			v, err := strconv.Atoi(parts[1])
			if err == nil {
				phrases[parts[0]] = v
			}
		}
	}
	return phrases
}

func readStopWords(fileName string) map[string]int {
	stopwords := make(map[string]int)
	file, err := os.Open(fileName)
	if err != nil {
		return stopwords
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		text := scanner.Text()
		stopwords[text] = 0
	}
	return stopwords
}

func processDataSet(fileName string, phrases map[string]int, stopwords map[string]int, t tokenizer.Tokenizer) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	product := Product{}
	for scanner.Scan() {

		text := scanner.Text()
		if len(text) == 0 {
			if len(product) > 0 {
				tokens := tokenize(t, product["Title"])
				subPhrases := generatePhrases(stopwords, tokens, 6)
				accumulate(phrases, subPhrases)

				tokens = tokenize(t, product["Feature"])
				subPhrases = generatePhrases(stopwords, tokens, 6)
				accumulate(phrases, subPhrases)
			}
			product = nil
			product = make(Product)
		}
		pair := strings.SplitN(text, "=", 2)
		if len(pair) == 1 {
			//productName := text
			//fmt.Printf("%s\n", productName)
		} else {
			product[pair[0]] = pair[1]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type Product map[string]string

// Data can be found https://github.com/dariusk/corpora/tree/master/data
// String tokenizer https://blog.gopheracademy.com/advent-2017/lexmachine-advent/
// Amazon training datasets https://github.com/SamTube405/Amazon-E-commerce-Data-set
func main() {
	knownPhrases := readKnownPhrases("phrases.txt")

	stopwords := readStopWords("stopwords.txt")

	t := tokenizer.NewWithSeparator("\t\n\r ,.:?\"!;()\\/\\-\\+\\&\\[\\]")
	phrases := make(map[string]int)

	files, err := ioutil.ReadDir("./")
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "dataset") {
			fmt.Printf("Processing %s...\n", f.Name())
			processDataSet("./"+f.Name(), phrases, stopwords, t)
		}
	}

	var ss []KeyValue
	for k, v := range phrases {
		ss = append(ss, KeyValue{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	f, err := os.OpenFile("phrases.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, kv := range ss {
		_, present := knownPhrases[kv.Key]
		if present {
			continue
		}

		fmt.Printf("%s?\n", kv.Key)
		var input string
		fmt.Scanln(&input)

		if "y" == input {
			f.WriteString(fmt.Sprintf("%s|%d|accept\n", kv.Key, kv.Value))
		}
		if "n" == input {
			f.WriteString(fmt.Sprintf("%s|%d|reject\n", kv.Key, kv.Value))
		}
		if "exit" == input {
			break
		}
	}

	fmt.Printf("Done.\n")
}
