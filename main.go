package main

import (
	"hash/crc32"
	"log"
	"sort"
	"strconv"
	"sync"
)

// Hash ring struct to hold the circle of nodes
type HashRing struct {
	nodes        map[uint32]string
	sortedHashes []uint32
	mu           sync.RWMutex
	replicas     int
}

func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		nodes:        make(map[uint32]string),
		sortedHashes: make([]uint32, 0),
		replicas:     replicas,
	}
}

func (hr *HashRing) AddNode(node string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	// add replicas for a node
	for i := 0; i < hr.replicas; i++ {
		hash := hr.hashKey(node + strconv.Itoa(i))
		hr.nodes[hash] = node
		hr.sortedHashes = append(hr.sortedHashes, hash)
	}

	// sort the sortedHashes after adding node
	sort.Slice(hr.sortedHashes, func(i, j int) bool { return hr.sortedHashes[i] < hr.sortedHashes[j] })
}

// hashKey generates a hash value for the given key using crc32
func (hr *HashRing) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (hr *HashRing) GetNode(key string) string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()

	if len(hr.sortedHashes) == 0 {
		return ""
	}

	hash := hr.hashKey(key)

	// search for nearest node for this key
	idx := hr.search(hash)
	log.Printf("GETNODE key : %s , hash : %d, sortedHash : %d", key, hash, hr.sortedHashes[idx])

	return hr.nodes[hr.sortedHashes[idx]]
}

func (hr *HashRing) RemoveNode(node string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	// remove replicas for a node
	for i := 0; i < hr.replicas; i++ {
		hash := hr.hashKey(node + strconv.Itoa(i))
		delete(hr.nodes, hash)
	}

	// reconstruct sortedHashes because replicas also removed, need to sort this
	hr.sortedHashes = make([]uint32, 0, len(hr.nodes))
	for hash := range hr.nodes {
		hr.sortedHashes = append(hr.sortedHashes, hash)
	}

	sort.Slice(hr.sortedHashes, func(i, j int) bool { return hr.sortedHashes[i] < hr.sortedHashes[j] })
}

// search finds the closest hash in the sorted list
func (hr *HashRing) search(hash uint32) int {
	idx := sort.Search(len(hr.sortedHashes), func(i int) bool { return hr.sortedHashes[i] >= hash })

	if idx < len(hr.sortedHashes) {
		return idx
	}

	return 0
}
func main() {
	log.Println("Consistent Hashing")
	ring := NewHashRing(8)
	ring.AddNode("node1")
	ring.AddNode("node2")
	ring.AddNode("node3")
	log.Printf("Available Nodes ==> \n %v", ring.nodes)
	log.Println("newkey ==> ", ring.GetNode("newkey"))
	log.Println("newkey2 ==> ", ring.GetNode("newkey2"))
	log.Println("newkey3 ==> ", ring.GetNode("newkey3"))
	log.Println("newkey4 ==> ", ring.GetNode("newkey4"))
	log.Println("newkey5 ==> ", ring.GetNode("newkey5"))

	ring.RemoveNode("node2")
	log.Printf("Available Nodes After remove ==> \n %v", ring.nodes)

}
