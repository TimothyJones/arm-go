package main

import "fmt"

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
