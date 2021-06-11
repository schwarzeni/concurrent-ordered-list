package concurrentorderdlist

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	// STATUS_DEFAULT default node status
	STATUS_DEFAULT uint32 = 0

	// STATUS_DELETED node has been deleted
	STATUS_DELETED uint32 = 1
)

type IntList struct {
	header *IntNode
	size   size
}

func (i *IntList) Contains(value int) (exists bool) {
	i.Range(func(v int) bool {
		if value == v {
			exists = true
			return false
		} else if value < v {
			return false
		} else {
			return true
		}
	})
	return
}

func (i *IntList) Insert(value int) bool {
	for {
		prevNode := i.header
		var currNode *IntNode

		// step1: find the suitable node and it's previous node
		for {
			currNode = prevNode.Next()

			// not exists, just return false
			if currNode != nil && currNode.value == value {
				return false
			}

			// find the suitable node
			if currNode == nil || currNode.value > value {
				break
			}
			prevNode = currNode
		}

		// step2: lock the previous node
		prevNode.mu.Lock()
		// previous node has been modified concurrently, goto step1 and try again
		if prevNode.Next() != currNode || prevNode.delete() {
			prevNode.mu.Unlock()
			continue
		}

		// step3: create new node and join into the linklist
		newNode := NewIntNode(value)
		newNode.SetNext(currNode)
		prevNode.SetNext(newNode)
		i.size.Add(1)

		// step4: unlock the previous node
		prevNode.mu.Unlock()
		return true
	}
}

func (i *IntList) Delete(value int) bool {
	for {
		prevNode := i.header
		var currNode *IntNode

		// step1: find the target node and it's previous node
		for {
			currNode = prevNode.Next()
			if currNode == nil || currNode.value > value {
				return false
			}
			if currNode.value == value {
				break
			}
			prevNode = currNode
		}

		// TODO: not know whether it is right to swap step2 and step3, both can pass tests anyway
		// step2: lock the previous node and check if it is valid
		prevNode.mu.Lock()
		if prevNode.Next() != currNode || prevNode.delete() {
			// if not, go back to step1
			prevNode.mu.Unlock()
			continue
		}

		// step3: lock the targe node and check if it is valid
		currNext := currNode.Next()
		currNode.mu.Lock()
		if currNext != currNode.Next() || currNode.delete() {
			// if not, go back to step1
			prevNode.mu.Unlock()
			currNode.mu.Unlock()
			continue
		}

		// step4: delete target node
		currNode.markAsDeleted()
		prevNode.SetNext(currNext)
		i.size.Add(-1)

		// step5: unlock the target node and previous node
		prevNode.mu.Unlock()
		currNode.mu.Unlock()
		return true
	}
}

func (i *IntList) Range(f func(value int) bool) {
	for currNode := i.header.Next(); currNode != nil; currNode = currNode.Next() {
		if !f(currNode.value) {
			return
		}
	}
}

func (i *IntList) Len() int {
	return i.size.Size()
}

func NewInt() *IntList {
	return &IntList{
		header: NewIntNode(0),
	}
}

type IntNode struct {
	value  int
	marked uint32
	next   *IntNode
	mu     sync.Mutex
}

func NewIntNode(val int) *IntNode {
	return &IntNode{value: val, marked: STATUS_DEFAULT}
}

func (n *IntNode) Next() *IntNode {
	return (*IntNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next))))
}

func (n *IntNode) SetNext(next *IntNode) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.next)), unsafe.Pointer(next))
}

func (n *IntNode) delete() bool {
	return atomic.LoadUint32(&n.marked) == STATUS_DELETED
}

func (n *IntNode) markAsDeleted() {
	atomic.StoreUint32(&n.marked, STATUS_DELETED)
}

type size struct {
	data int64
}

func (s *size) Add(n int) {
	atomic.AddInt64(&s.data, int64(n))
}

func (s *size) Size() int {
	return int(atomic.LoadInt64(&s.data))
}
