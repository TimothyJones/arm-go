package main

import (
	"log"
	"sync"
)

var minSupport int
var wg sync.WaitGroup

func init() {
	counts = make(map[int]int)
	tokens = make(map[string]int)
	itemStrings = make(map[int]string)
}

func main() {
	log.Println("Beginning")
	minSupport = 0

	toTree := make(chan []int, 100000)

	countWords()
	wg.Add(1)
	go treeConstuct(toTree)
	readTransactions(toTree)

	wg.Wait()
}
