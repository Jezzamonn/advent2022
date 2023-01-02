package main

import (
	"container/heap"
)

type NodeDistPair struct {
	Node *Node
	Dist int
}

type NodePriorityQueue []NodeDistPair

func (pq NodePriorityQueue) Len() int { return len(pq) }

func (pq NodePriorityQueue) Less(i, j int) bool {
	return pq[i].Dist < pq[j].Dist
}

func (pq NodePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *NodePriorityQueue) Push(x interface{}) {
	item := x.(NodeDistPair)
	*pq = append(*pq, item)
}

func (pq *NodePriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

var _ heap.Interface = &NodePriorityQueue{}
