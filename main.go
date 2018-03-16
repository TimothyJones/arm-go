package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/chrislusf/glow/flow"
)

var counts map[string]int
var minSupport int

var wg sync.WaitGroup

type byCounts []string

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
	Value    string
	Count    int
	Children map[string]*tree
}

func emptyTree() *tree {
	return &tree{"", 0, make(map[string]*tree)}
}

func newNode(value string) *tree {
	return &tree{value, 0, make(map[string]*tree)}
}

func (t *tree) insert(value string) *tree {
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
	if t.Value == "" {
		fmt.Printf("( Root count: %d ", t.Count)
	} else {
		fmt.Printf("( item: %s, count: %d ", t.Value, t.Count)
	}
	for _, v := range t.Children {
		v.print()
	}
	fmt.Print(")")

}

func treeConstuct(c chan []string) {
	defer wg.Done()
	root := emptyTree()
	for transaction := range c {
		root.Count++
		current := root
		for _, item := range transaction {
			current = current.insert(item)
		}
	}
	root.print()
}

func main() {
	counts = make(map[string]int)
	minSupport = 4

	toTree := make(chan []string)
	wg.Add(1)
	go treeConstuct(toTree)

	flow.New().TextFile(
		os.Args[1], 3,
	).Map(func(line string, ch chan string) {
		for _, token := range strings.Split(line, ",") {
			ch <- token
		}
	}).Map(func(key string) (string, int) {
		return key, 1
	}).ReduceByKey(func(x int, y int) int {
		return x + y
	}).Map(func(key string, value int) {
		counts[key] = value
	}).Run()

	flow.New().TextFile(
		os.Args[1], 3,
	).Map(func(line string) []string {
		s := strings.Split(line, ",")
		sort.Sort(byCounts(s))
		return s
	}).Map(func(transaction []string) []string {
		for i, s := range transaction {
			if counts[s] < minSupport {
				return transaction[:i]
			}
		}
		return transaction
	}).Map(func(transaction []string) {
		/*var buffer bytes.Buffer
		for _, s := range transaction {
			buffer.WriteString(fmt.Sprintf("(%s,%d),", s, counts[s]))
		}
		fmt.Println(buffer.String())*/
		toTree <- transaction
	}).Run()
	close(toTree)

	wg.Wait()
	/*ufoar key := range counts {
		fmt.Printf("(%s): %d\n", key, counts[key])
	}*/
}
