package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/chrislusf/glow/flow"
)

var counts map[int]int
var tokens map[string]int
var itemStrings map[int]string
var minSupport int

var wg sync.WaitGroup

type byCounts []int

func (s byCounts) Len() int {
	return len(s)
}
func (s byCounts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byCounts) Less(i, j int) bool {
	return counts[s[i]] > counts[s[j]]
}

type tree struct {
	Value    int
	Count    int
	Children map[int]*tree
}

func emptyTree() *tree {
	return &tree{0, 0, make(map[int]*tree)}
}

func newNode(value int) *tree {
	return &tree{value, 0, make(map[int]*tree)}
}

func (t *tree) insert(value int) *tree {
	node, ok := t.Children[value]
	if !ok {
		node = newNode(value)
		t.Children[value] = node
	}
	node.Count++
	return node
}

func (t *tree) print() {
	/*	if t.Count < minSupport {
		return
	}*/
	if t.Value == 0 {
		fmt.Printf("( Root count: %d ", t.Count)
	} else {
		fmt.Printf("( item: %s, count: %d ", itemStrings[t.Value], t.Count)
	}
	for _, v := range t.Children {
		v.print()
	}
	fmt.Print(")")

}

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
	root.print()
}

func main() {
	log.Println("Beginning")
	counts = make(map[int]int)
	tokens = make(map[string]int)
	itemStrings = make(map[int]string)
	minSupport = 4

	toTree := make(chan []int)
	wg.Add(1)
	go treeConstuct(toTree)

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Panicf("Can not open file %s: %v", os.Args[1], err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	tokenId := 1
	for scanner.Scan() {
		for _, token := range strings.Split(scanner.Text(), ",") {
			if _, ok := tokens[token]; !ok {
				tokens[token] = tokenId
				itemStrings[tokenId] = token
				tokenId++
			}
			id := tokens[token]
			counts[id] = counts[id] + 1
		}
	}

	log.Println("Counted items")

	flow.New().TextFile(
		os.Args[1], 3,
	).Map(func(line string) []int {
		s := strings.Split(line, ",")
		transaction := make([]int, len(s))
		for i, item := range s {
			transaction[i] = tokens[item]
		}
		sort.Sort(byCounts(transaction))
		return transaction
	}).Map(func(transaction []int) []int {
		for i, s := range transaction {
			if counts[s] < minSupport {
				return transaction[:i]
			}
		}
		return transaction
	}).Map(func(transaction []int) {
		/*var buffer bytes.Buffer
		for _, s := range transaction {
			buffer.WriteString(fmt.Sprintf("(%s,%d),", s, counts[s]))
		}
		fmt.Println(buffer.String())*/
		toTree <- transaction
	}).Run()
	close(toTree)
	log.Println("Read transactions")

	wg.Wait()
	/*ufoar key := range counts {
		fmt.Printf("(%s): %d\n", key, counts[key])
	}*/
}
