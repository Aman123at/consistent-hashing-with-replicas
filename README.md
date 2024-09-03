# Consistent Hashing in Go

This Go implementation demonstrates a consistent hashing mechanism. Consistent hashing is a technique used in distributed systems to distribute data across multiple nodes efficiently. This implementation includes the following key features:

- Adding and removing nodes dynamically
- Evenly distributing keys across nodes
- Finding the nearest node for a given key

## Code Overview

### Package Imports

The package imports include:
- `hash/crc32`: For generating a checksum for the keys.
- `log`: For logging output to the console.
- `sort`: For sorting slices.
- `strconv`: For converting integers to strings.
- `sync`: For thread-safe operations on shared data.

### HashRing Struct

The `HashRing` struct holds the necessary data for the hash ring:
- `nodes`: A map storing node hashes to node names.
- `sortedHashes`: A slice of sorted hash values for quick access.
- `mu`: A mutex to ensure thread-safe read/write operations.
- `replicas`: The number of virtual nodes (replicas) per physical node.

```go
type HashRing struct {
    nodes        map[uint32]string
    sortedHashes []uint32
    mu           sync.RWMutex
    replicas     int
}
```

## Conclusion

This implementation of consistent hashing in Go is efficient for distributing data across multiple nodes in a distributed system. By using replicas and sorting, the `HashRing` ensures that keys are evenly distributed and can be quickly retrieved or reassigned when nodes are added or removed.
