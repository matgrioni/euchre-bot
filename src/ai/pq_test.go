package ai

import (
    "container/heap"
    "testing"
)

type TestItem struct {
    value interface{}
    priority float64
    index int
}

func (t* TestItem) GetValue() interface{} {
    return t.value
}

func (t* TestItem) Value(value interface{}) {
    t.value = value
}

func (t *TestItem) GetPriority() float64 {
    return t.priority
}

func (t *TestItem) Priority(priority float64) {
    t.priority = priority
}

func (t *TestItem) GetIndex() int {
    return t.index
}

func (t *TestItem) Index(index int) {
    t.index = index
}

func TestUpdateMaxHeap(t *testing.T) {
    items := map[string]float64 {
        "banana": 3, "apple": 2, "pear": 4,
    }

    pq := make(PriorityQueue, len(items))
    i := 0
    for value, priority := range items {
        pq[i] = &TestItem {
            value,
            priority,
            i,
        }

        i++
    }
    heap.Init(&pq)

    item := &TestItem {
        "pineapple",
        5,
        0,
    }
    heap.Push(&pq, item)
    pq.Update(item, item.GetValue(), item.GetPriority())

    res := pq.Poll()
    if res != item {
        // TODO: Improve error message.
        t.Error("Expected pineapple not %s.\n", res)
    }
}
