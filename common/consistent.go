package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type units []uint32

func (x units) Len() int {
	return len(x)
}

func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("Hash 环没有数据.")

type Consistent struct {
	//hash circle key：hashed value：info
	Circle map[uint32]string
	//sorted hashed slice
	sortedHashes units
	//virtual node number
	VirtualNode int
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		Circle:      make(map[uint32]string),
		VirtualNode: 20,
	}
}

func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSoreHashes() {
	hashes := c.sortedHashes[:0]
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.Circle) {
		hashes = nil
	}
	for k := range c.Circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedHashes = hashes
}
func (c *Consistent) Add(element string) {
	// lock  unlock
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		c.Circle[c.hashKey(c.generateKey(element, i))] = element
	}
	c.updateSoreHashes()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.Circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSoreHashes()
}

//顺时针查找最近的节点
func (c *Consistent) search(key uint32) int {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

func (c *Consistent) Get(name string) (string, error) {
	c.Lock()
	defer c.Unlock()

	if len(c.Circle) == 0 {
		return "", errEmpty
	}
	key := c.hashKey(name)
	i := c.search(key)
	return c.Circle[c.sortedHashes[i]], nil
}
