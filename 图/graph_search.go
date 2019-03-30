package main

import (
	"container/list"
	"fmt"
)

//adjacency table, 无向图
type Graph struct {
	adj []*list.List
	val int
}

// init graphh according to capacity
func newGraph(val int) *Graph {
	graph := &Graph{}
	graph.val = val
	graph.adj = make([]*list.List, val)
	for i := range graph.adj {
		graph.adj[i] = list.New()
	}

	return graph
}

//insert as add edge，一条边存2次
func (this *Graph) addEdge(start int, target int) {
	this.adj[start].PushBack(target)
	this.adj[target].PushBack(start)
}

// search path by BFS
func (this *Graph) BFS(start int, target int) {
	// todo
	if start == target {
		return
	}

	//init prev
	prev := make([]int, this.val)
	for index := range prev {
		prev[index] = -1
	}

	//search by queue
	var queue []int
	visited := make([]bool, this.val)
	queue = append(queue, start)
	visited[start] = true
	isFound := false
	for len(queue) > 0 && !isFound {
		top := queue[0]
		queue = queue[1:]
		linkedList := this.adj[top]
		for e := linkedList.Front(); e != nil; e = e.Next() {
			k := e.Value.(int)
			if !visited[k] {
				prev[k] = top
				if k == target {
					isFound = true
					break
				}
				queue = append(queue, k)
				visited[k] = true
			}
		}

		if isFound {
			printPrev(prev, start, target)
		} else {
			fmt.Printf("no path found from %d to %d\n", start, target)
		}
	}
}

// seatch by DFS
func (this *Graph) DFS(start int, target int) {
	prev := make([]int, this.val)
	for i := range prev {
		prev[i] = -1
	}

	visited := make([]bool, this.val)
	visited[s] = true

	isFound := false
	this.recurse(start, target, prev, visited, isFound)

	printPrev(prev, start, target)
}

//recursivly find path
func (this *Graph) recurse(start int, target int, prev []int, visited []bool, isFound bool) {
	if isFound {
		return
	}

	visited[start] = true

	if start == target {
		isFound = true
		return
	}

	linkedlist := this.adj[start]
	for e := linkedlist.Front(); e != nil; e = e.Next() {
		k := e.Value.(int)
		if !visited[k] {
			prev[k] = start
			this.recurse(k, target, prev, visited, false)
		}
	}
}

//print path recursively
func printPrev(prev []int, s int, t int) {

	if t == s || prev[t] == -1 {
		fmt.Printf("%d ", t)
	} else {
		printPrev(prev, s, prev[t])
		fmt.Printf("%d ", t)
	}

}
