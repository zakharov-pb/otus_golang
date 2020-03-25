package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"strings"
)

var regWords = regexp.MustCompile(`[A-Za-zА-Яа-я0-9ёЁ]+(-*[A-Za-zА-Яа-яёЁ0-9])*`)

func getWords(text string) (words []string) {
	if len(text) == 0 {
		return nil
	}
	words = regWords.FindAllString(text, -1)
	length := len(words)
	if length == 0 {
		return nil
	}
	for i := 0; i < length; i++ {
		words[i] = strings.ToLower(words[i])
	}
	return words
}

func topN(text string, maxCount int) []string {
	words := getWords(text)
	if words == nil {
		return nil
	}

	indexes := map[string]int{}
	for _, w := range words {
		indexes[w]++
	}
	length := len(indexes)
	words = words[:length]
	index := 0
	for k := range indexes {
		words[index] = k
		index++
	}

	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ { //nolint:gomnd,stylecheck
			if indexes[words[i]] < indexes[words[j]] {
				words[i], words[j] = words[j], words[i]
			}
		}
	}
	if length <= maxCount {
		return words
	}
	return words[:maxCount]
}

func Top10(text string) []string {
	return topN(text, 10)
}
