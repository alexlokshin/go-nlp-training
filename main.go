package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/euskadi31/go-tokenizer"
)

// stopwords: with, for, and,

var stopwords map[string]int

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
func isStopWord(text string) bool {
	_, v := stopwords[text]
	return v
}
func generatePhrases(keywords []string, maxLength int) map[string]int {
	phrases := make(map[string]int)
	for i := 0; i < len(keywords); i++ {
		//collect(keywords[i], phrases)
		if !isStopWord(keywords[i]) && !valid.IsInt(keywords[i]) {
			if i < len(keywords)-1 {
				phrase := keywords[i]
				for j := 1; j < maxLength && j+i < len(keywords); j++ {
					if !isStopWord(keywords[i+j]) {
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

type Product map[string]string

// Data can be found https://github.com/dariusk/corpora/tree/master/data
// String tokenizer https://blog.gopheracademy.com/advent-2017/lexmachine-advent/
func main() {
	stopwords = make(map[string]int)
	stopwords["for"] = 0
	stopwords["with"] = 0
	stopwords["to"] = 0
	stopwords["and"] = 0

	t := tokenizer.NewWithSeparator("\t\n\r ,.:?\"!;()\\/\\-\\+\\&")
	phrases := make(map[string]int)

	// title := "SMSL Sanskrit 24bit192kHz USB/Coaxial/Optical Digital To Analog Audio Decoder Converter (silver) ,by Gemini Doctor"
	// tokens := tokenize(t, title)
	// phrases := generatePhrases(tokens, 3)
	// for k, v := range phrases {
	// 	fmt.Printf("%s, %d\n", k, v)
	// }

	file, err := os.Open("dataset.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//products := []Product{}
	product := Product{}
	for scanner.Scan() {

		text := scanner.Text()
		if len(text) == 0 {
			if len(product) > 0 {
				//fmt.Printf("%s\n", product["Title"])
				tokens := tokenize(t, product["Title"])
				subPhrases := generatePhrases(tokens, 3)
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
		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var ss []KeyValue
	for k, v := range phrases {
		ss = append(ss, KeyValue{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	counter := 0
	for _, kv := range ss {
		fmt.Printf("%s, %d\n", kv.Key, kv.Value)
		if counter > 200 {
			break
		}
		counter++
	}

	fmt.Printf("Done.\n")
}
