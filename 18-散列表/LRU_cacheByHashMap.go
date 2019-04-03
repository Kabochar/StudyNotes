package main

const (
	hostbit = uint64(^uint(0)) == ^uint64(0)
	LENGTH  = 100
)

type LRUNode struct {
	prev *LRUNode
	next *LRUNode

	key   int // lru key
	value int // lru value

	hnext *LRUNode // 拉链
}

type LRUCache struct {
	node []LRUNode // hash list

	head *LRUNode // lru head node
	tail *LRUNode // lru tail node

	capacity int //
	used     int //
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		node:     make([]LRUNode, LENGTH),
		head:     nil,
		tail:     nil,
		capacity: capacity,
		used:     0,
	}
}

func (this *LRUCache) Get(key int) int {
	if this.tail == nil {
		return -1
	}

	if tmp := this.searchNode(key); tmp != nil {
		this.moveToTail(tmp)
		return tmp.value
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	// 1. 首次插入数据
	// 2. 插入数据不在 LRU 中
	// 3. 插入数据在 LRU 中
	// 4. 插入数据不在 LRU 中, 并且 LRU 已满

	if tmp := this.searchNode(key); tmp != nil {
		tmp.value = value
		this.moveToTail(tmp)
		return
	}
	this.addNode(key, value)

	if this.used > this.capacity {
		this.delNode()
	}
}

func (this *LRUCache) addNode(key int, value int) {
	newNode := &LRUNode{
		key:   key,
		value: value,
	}

	tmp := &this.node[hash(key)]
	newNode.hnext = tmp.hnext
	tmp.hnext = newNode
	this.used++

	if this.tail == nil {
		this.tail, this.head = newNode, newNode
		return
	}
	this.tail.next = newNode
	newNode.prev = this.tail
	this.tail = newNode
}

func (this *LRUCache) delNode() {
	if this.head == nil {
		return
	}
	prev := &this.node[hash(this.head.key)]
	tmp := prev.hnext

	for tmp != nil && tmp.key != this.head.key {
		prev = tmp
		tmp = tmp.hnext
	}
	if tmp == nil {
		return
	}
	prev.hnext = tmp.hnext
	this.head = this.head.next
	this.head.prev = nil
	this.used--
}

func (this *LRUCache) searchNode(key int) *LRUNode {
	if this.tail == nil {
		return nil
	}

	// 查找
	tmp := this.node[hash(key)].hnext
	for tmp != nil {
		if tmp.key == key {
			return tmp
		}
		tmp = tmp.hnext
	}
	return nil
}

func (this *LRUCache) moveToTail(node *LRUNode) {
	if this.tail == node {
		return
	}
	if this.head == node {
		this.head = node.next
		this.head.prev = nil
	} else {
		node.next.prev = node.prev
		node.prev.next = node.next
	}

	node.next = nil
	this.tail.next = node
	node.prev = this.tail

	this.tail = node
}

func hash(key int) int {
	if hostbit {
		return (key ^ (key >> 32)) & (LENGTH - 1)
	}
	return (key ^ (key >> 16)) & (LENGTH - 1)
}