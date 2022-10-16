package LRU

import (
	"container/list"
	"runtime"
	"sync"
)

/*
doc1: https://halfrost.com/lru_lfu_interview/#toc-1
doc2: https://carlclone.github.io/algos/lru-cache-go/
thread safe: use Read/Write Lock on Get() / Put() / Size()
high concurrency : use sharding , noted that this code won't pass the leetcode test ,
the reason is the capacity of this cache is approximately , not absolutely , eviction won't act as leetcode expect

for example , concurrency=3 , approximate capacity=7 , 7/3=2 , 2 segments , every segment has 7/2=3 ,
we need extra 1 capacity , so add 1 at every segment 3+1=4 capacity ,
so the absolute capacity will be 4*2=8
*/

type LRUCacheShard struct {
	Cap  int
	Map  map[int]*list.Element
	List *list.List
	Lock sync.RWMutex
}

type LRUCache struct {
	Shard map[int]*LRUCacheShard
}

type pair struct {
	K, V int
}

func Constructor(capacity int) LRUCache {
	cpuNum := runtime.NumCPU()
	shard := map[int]*LRUCacheShard{}
	capSeg := capacity / cpuNum
	// determine odd or even number
	if capacity&cpuNum == 1 {
		capSeg++
	}
	for i := 0; i < cpuNum; i++ {
		value := &LRUCacheShard{
			Cap:  capSeg,
			Map:  map[int]*list.Element{},
			List: list.New(),
		}
		shard[i] = value
	}
	return LRUCache{
		Shard: shard,
	}
}

func (c *LRUCache) getShard(key int) *LRUCacheShard {
	return c.Shard[key%runtime.NumCPU()]
}

func (c *LRUCache) Get(key int) int {
	shard := c.getShard(key)
	shard.Lock.RLock()
	defer shard.Lock.RUnlock()
	if el, ok := shard.Map[key]; ok {
		shard.List.MoveToFront(el)
		return el.Value.(pair).V
	}
	return -1
}

func (c *LRUCache) Put(key int, value int) {
	shard := c.getShard(key)
	shard.Lock.Lock()
	defer shard.Lock.Unlock()
	if el, ok := shard.Map[key]; ok {
		el.Value = pair{K: key, V: value}
		shard.List.MoveToFront(el)
	} else {
		el := shard.List.PushFront(pair{K: key, V: value})
		shard.Map[key] = el
	}
	if shard.List.Len() > shard.Cap {
		el := shard.List.Back()
		shard.List.Remove(el)
		delete(shard.Map, el.Value.(pair).K)
	}
}
