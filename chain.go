package main

import (
	"math/rand"
	"strings"
)

const (
	ChainStart = "\x01S"
	ChainEnd   = "\x01E"
)

type Chain struct {
	chain map[string]map[string]int
}

func NewChain() Chain {
	return Chain{make(map[string]map[string]int)}
}

func (c *Chain) Clone(chain *Chain) {
	chain.chain = c.chain
}

func (c *Chain) addWord(rootWord, word string) {
	if c.chain[rootWord] == nil {
		c.chain[rootWord] = make(map[string]int)
	}
	c.chain[rootWord][word]++
}

func (c *Chain) AddSentence(words []string) {
	c.addWord(ChainStart, words[0])
	for i := 0; i < len(words)-1; i++ {
		c.addWord(words[i], words[i+1])
	}
	c.addWord(words[len(words)-1], ChainEnd)
}

func (c *Chain) AddText(text string) {
	text = strings.ToLower(text)
	sentences := strings.Split(text, ".")
	if sentences[len(sentences)-1] == "" {
		sentences = sentences[:len(sentences)-1]
	}
	for i := range sentences {
		sentence := strings.Split(sentences[i], " ")
		if sentence[0] == "" {
			sentence = sentence[1:]
		}
		c.AddSentence(sentence)
	}
}

func (c *Chain) RandomSentence() string {
	sentence := ""
	lastWord := ChainStart
	for {
		var possibleWords []string
		for word := range c.chain[lastWord] {
			for i := 0; i < c.chain[lastWord][word]; i++ {
				possibleWords = append(possibleWords, word)
			}
		}

		rand.Shuffle(len(possibleWords), func(i, j int) {
			possibleWords[i], possibleWords[j] = possibleWords[j], possibleWords[i]
		})

		selectedWord := rand.Intn(len(possibleWords))
		if possibleWords[selectedWord] == ChainEnd {
			sentence += "."
			break
		}
		sentence += possibleWords[selectedWord] + " "
		lastWord = possibleWords[selectedWord]
	}
	sentence = strings.Replace(sentence, " .", ".", -1)
	return sentence
}
