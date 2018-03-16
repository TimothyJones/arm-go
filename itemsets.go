package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
)

func treeConstuct(c chan []int) {
	defer wg.Done()
	root := emptyTree()
	for transaction := range c {
		root.Count++
		current := root
		for _, item := range transaction {
			current = current.insert(item)
		}
	}
	log.Println("Created tree")
	//root.print()
}

func countWords() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Panicf("Can not open file %s: %v", os.Args[1], err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	tokenID := 1
	for scanner.Scan() {
		for _, token := range strings.Split(scanner.Text(), ",") {
			if _, ok := tokens[token]; !ok {
				tokens[token] = tokenID
				itemStrings[tokenID] = token
				tokenID++
			}
			id := tokens[token]
			counts[id] = counts[id] + 1
		}
	}

	log.Println("Counted items")
}

func readTransactions(toTree chan []int) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Panicf("Can not open file %s: %v", os.Args[1], err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		transaction := make([]int, len(s))
		j := 0
		for _, token := range s {
			id := tokens[token]
			if counts[id] >= minSupport {
				transaction[j] = id
				j++
			}
		}
		transaction = transaction[:j]
		sort.Sort(byCounts(transaction))
		toTree <- transaction
	}
	close(toTree)

	log.Println("Read transactions")
}
