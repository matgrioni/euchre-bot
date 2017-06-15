package ai

import (
    "container/heap"
)

// This file provides a MaxHeap / PriorityQueue for a PQItem. To use the
// PriorityQueue, create a struct that implements PQItem.

type PQItem interface {
    Value() interface{}
    Value(v interface{})

    Priority() int
    Priority(priority int)

    Index() int
    Index(i int)
}

type PriorityQueue []*PQItem

func (pq PriorityQueue) Len() int {
    return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
    // Since this is a MaxHeap, and container/heap implements a MinHeap, we need
    // to use greater than.
    return pq[i].Priority() > pq[j].Priority()
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].Index(i)
    pq[j].Index(j)
}

func (pq *PriorityQueue) Push(item PQItem) {
    n := len(*pq)
    item.Index(n)
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() PQItem {
    old := *pq
    n := len(old)

    item := old[n - 1]
    item.index = -1
    *pq = old[:n - 1]

    return item
}

func (pq *PriorityQueue) Update(item *PQItem, value interface{}, priority int) {
    item.Value(value)
    item.Priority(priority)
    heap.Fix(pq, item.Index())
}
