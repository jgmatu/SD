package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func countWords(file string, m_words map[string]int) {
	fp, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		w := scanner.Text() // word of file...
		if _, ok := m_words[w]; ok {
			m_words[w]++
		} else {
			m_words[w] = 1
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fp.Close()
}

func printwords(m_words map[string]int) {
	ws := make([]string, len(m_words))

	// Sort...
	i := 0
	for w := range m_words {
		ws[i] = w
		i++
	}
	sort.Strings(ws)

	// Print...
	for _, w := range ws {
		fmt.Printf("%12s\t%4d\n", w, m_words[w])
	}
}

func main() {
	m_words := make(map[string]int)

	for _, arg := range os.Args[1:] {
		countWords(arg, m_words)
	}
	printwords(m_words)
}
