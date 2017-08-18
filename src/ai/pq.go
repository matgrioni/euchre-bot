package ai

import "container/heap"

// This file provides a MaxHeap / PriorityQueue for a PQItem. To use the
// PriorityQueue, create a struct that implements PQItem.
// TODO: Include example of initialization and usage code.

type PQItem interface {
    GetValue() interface{}
    Value(v interface{})

    GetPriority() float64
    Priority(priority float64)

    GetIndex() int
    Index(i int)
}

type PriorityQueue []PQItem

func (pq PriorityQueue) Len() int {
    return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
    // Since this is a MaxHeap, and container/heap implements a MinHeap, we need
    // to use greater than.
    return pq[i].GetPriority() > pq[j].GetPriority()
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].Index(i)
    pq[j].Index(j)
}

func (pq *PriorityQueue) Push(item interface{}) {
    n := len(*pq)
    x := item.(PQItem)
    x.Index(n)
    *pq = append(*pq, x)
}

func (pq *PriorityQueue) Poll() PQItem {
    return (*pq)[0]
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)

    item := old[n - 1]
    item.Index(-1)
    *pq = old[:n - 1]

    return item
}

func (pq *PriorityQueue) Update(item PQItem, value interface{}, priority float64) {
    item.Value(value)
    item.Priority(priority)
    heap.Fix(pq, item.GetIndex())
}
