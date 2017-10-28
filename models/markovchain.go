package models

import (
	"database/sql"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

type WordChain struct {
	Word1 string
	Word2 string
	Word3 string
}

type MarkovChain struct {
	Word1 string
	Word2 string
	Word3 string
	Time  int
	ID    int
}

var word_count_spam_threshold int = 50
var recently_learned []string
var recently_learned_amount int = 10

func MarkovCount() int {
	row := db.QueryRow("SELECT COUNT(*) FROM Markov")
	var count int
	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}
	return count
}

func saveNewChain(chain WordChain, timestamp int) {
	insertSql, err := db.Prepare("INSERT INTO Markov(id,word1,word2,word3,time) VALUES (NULL,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	defer insertSql.Close()
	_, err = insertSql.Exec(chain.Word1, chain.Word2, chain.Word3, timestamp)
	if err != nil {
		panic(err)
	}
}

func getMatchingChains(word2 string) []MarkovChain {
	rows, err := db.Query("SELECT id,word1,word2,word3,time FROM Markov WHERE word2 = ?", word2)
	var matchChains []MarkovChain
	switch {
	case err == sql.ErrNoRows:
		return matchChains
	case err != nil:
		panic(err)
	default:
	}

	for rows.Next() {
		markovChain := MarkovChain{}
		rows.Scan(&markovChain.ID, &markovChain.Word1, &markovChain.Word2, &markovChain.Word3, &markovChain.Time)
		matchChains = append(matchChains, markovChain)
	}
	return matchChains
}

func getPrefixChains(word2 string, word3 string) []MarkovChain {
	rows, err := db.Query("SELECT id,word1,word2,word3,time FROM Markov Where word2 = ? AND word3 = ?", word2, word3)
	var matchChains []MarkovChain
	switch {
	case err == sql.ErrNoRows:
		return matchChains
	case err != nil:
		panic(err)
	default:
	}

	for rows.Next() {
		markovChain := MarkovChain{}
		rows.Scan(&markovChain.ID, &markovChain.Word1, &markovChain.Word2, &markovChain.Word3, &markovChain.Time)
		matchChains = append(matchChains, markovChain)
	}
	return matchChains
}

func getSuffixChains(word1 string, word2 string) []MarkovChain {
	rows, err := db.Query("SELECT id,word1,word2,word3,time FROM Markov Where word1 = ? AND word2 = ?", word1, word2)
	var matchChains []MarkovChain
	switch {
	case err == sql.ErrNoRows:
		return matchChains
	case err != nil:
		panic(err)
	default:
	}

	for rows.Next() {
		markovChain := MarkovChain{}
		rows.Scan(&markovChain.ID, &markovChain.Word1, &markovChain.Word2, &markovChain.Word3, &markovChain.Time)
		matchChains = append(matchChains, markovChain)
	}
	return matchChains
}

func isSpam(words []string, message string) bool {
	if len(words) > word_count_spam_threshold {
		return true
	}

	//Generate a map of word:count for all of the words
	wordMap := make(map[string]int)
	for _, word := range words {
		count, exist := wordMap[word]
		if exist {
			wordMap[word] = count + 1
		} else {
			wordMap[word] = 1
		}
	}

	//Generate a list of all of the counts from the word map and a sum of all of the values
	var values []int
	sumValues := 0
	for _, count := range wordMap {
		values = append(values, count)
		sumValues += count
	}
	sort.Ints(values)

	most_common_val, highest_count, highest_val, last_value, count_this_val := 0, 0, 0, 0, 0
	for _, value := range values {
		if value != last_value {
			count_this_val = 1
		} else {
			count_this_val += 1
		}

		if count_this_val > highest_count {
			highest_count = count_this_val
			most_common_val = value
		}
		if value > highest_val {
			highest_val = value
		}
		last_value = value
	}

	if most_common_val != 1 {
		return true
	}

	if highest_val > 1+len(words)/4 {
		return true
	}

	if sumValues/len(values) >= 2 {
		return true
	}

	//Check if any single word is longer than 30 characters after removing links
	linkfind := regexp.MustCompile(`(http://|https://)([\S]+)`)
	for _, word := range words {
		if len(word) > 30 {
			if linkfind.FindString(word) == "" {
				return true
			}
		}
	}

	//Remove links from message, then get the number of sentences in the message and reject if 5 or more
	linkremoved := regexp.MustCompile(`(http://|https://)([\S]+)`)
	messagenolink := linkremoved.ReplaceAllString(message, "")
	sentenceregex := regexp.MustCompile(`([\S\s]+?)[.?!]`)
	if len(sentenceregex.FindAllStringSubmatch(messagenolink, -1)) >= 5 {
		return true
	}

	return false
}

func isRecentlyLearned(message string) bool {
	if len(recently_learned) > recently_learned_amount {
		_, recently_learned = recently_learned[0], recently_learned[1:]
	}

	for _, oldmsg := range recently_learned {
		if message == oldmsg {
			return true
		}
	}

	recently_learned = append(recently_learned, message)
	return false
}

func getWord(words []string, i int) string {
	if i < 0 || i >= len(words) {
		return ""
	}
	return words[i]
}

func wordChaining(words []string) []WordChain {
	var chains []WordChain

	for i, word := range words {
		newChain := WordChain{}
		newChain.Word1 = getWord(words, i-1)
		newChain.Word2 = word
		newChain.Word3 = getWord(words, i+1)
		chains = append(chains, newChain)
	}

	return chains
}

func messageSplit(message string) []string {
	words := strings.Split(message, " ")

	//Empty message check
	if len(words) == 0 {
		return nil
	}

	if isRecentlyLearned(message) {
		return nil
	}

	if isSpam(words, message) {
		return nil
	}

	return words
}

func LearnMarkov(message string) {
	words := messageSplit(message)
	if len(words) == 0 {
		return
	}

	timestamp := int(time.Now().Unix())
	chains := wordChaining(words)

	for _, chain := range chains {
		saveNewChain(chain, timestamp)
	}
}

func hasPuncSuffix(word string) bool {
	switch {
	case strings.HasSuffix(word, "!"):
		return true
	case strings.HasSuffix(word, "."):
		return true
	case strings.HasSuffix(word, ":"):
		return true
	case strings.HasSuffix(word, ";"):
		return true
	case strings.HasSuffix(word, ","):
		return true
	case strings.HasSuffix(word, "?"):
		return true
	default:
		return false
	}
}

func isJoiner(word string) bool {
	if len(word) > 1 && hasPuncSuffix(word) {
		return true
	}

	joiners := map[string]struct{}{"a": {}, "if": {}, "its": {}, "it's": {}, "and": {}, "or": {}, "because": {}, "with": {}, "when": {},
		"like": {}, "then": {}, "than": {}, "after": {}, "also": {}, "before": {}}

	_, isJoin := joiners[word]
	return isJoin
}

func GenerateMarkovResponse(startingWord string) string {
	chains := getMatchingChains(startingWord)
	if len(chains) == 0 {
		return ""
	}

	start := chains[rand.Intn(len(chains))]

	if start.Word1 == "" {
		start = chains[rand.Intn(len(chains))]
	}

	current := start
	var head []MarkovChain
	var tail []MarkovChain

	for {
		chains = getPrefixChains(current.Word1, current.Word2)
		if len(chains) == 0 {
			break
		}
		current = chains[rand.Intn(len(chains))]

		if isJoiner(current.Word2) {
			chains = getMatchingChains(current.Word2)
			current = chains[rand.Intn(len(chains))]
		}

		head = append(head, current)

		if current.Word1 == "" {
			break
		}
	}

	current = start
	for {
		chains = getSuffixChains(current.Word2, current.Word3)
		if len(chains) == 0 {
			break
		}
		current = chains[rand.Intn(len(chains))]

		if isJoiner(current.Word2) {
			chains = getMatchingChains(current.Word2)
			current = chains[rand.Intn(len(chains))]
		}

		tail = append(tail, current)

		if current.Word3 == "" {
			break
		}
	}

	var words []string

	for i := len(head) - 1; i >= 0; i-- {
		words = append(words, head[i].Word2)
	}
	words = append(words, start.Word2)
	for _, word := range tail {
		words = append(words, word.Word2)
	}

	return strings.Join(words, " ")
}

func RandomResponse(message string) string {
	words := strings.Split(message, " ")
	if len(words) == 0 {
		return ""
	}

	word := words[rand.Intn(len(words))]
	return GenerateMarkovResponse(word)
}
