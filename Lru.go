package LRU

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	Cap  int
	Map  map[int]*list.Element
	List *list.List
	Lock sync.RWMutex
}

type pair struct {
	K, V int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		Cap:  capacity,
		Map:  make(map[int]*list.Element),
		List: list.New(),
	}
}

func (c *LRUCache) Get(key int) int {
	// 需移动头部节点，因此加写锁
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if el, ok := c.Map[key]; ok {
		c.List.MoveToFront(el)
		return el.Value.(pair).V
	}
	return -1
}

func (c *LRUCache) Put(key int, value int) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if el, ok := c.Map[key]; ok {
		el.Value = pair{K: key, V: value}
		c.List.MoveToFront(el)
	} else {
		el := c.List.PushFront(pair{K: key, V: value})
		c.Map[key] = el
	}
	if c.List.Len() > c.Cap {
		el := c.List.Back()
		c.List.Remove(el)
		delete(c.Map, el.Value.(pair).K)
	}
}
