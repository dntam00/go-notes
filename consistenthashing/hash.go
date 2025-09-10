package consistenthashing

import (
	"cmp"
	"errors"
	"log"
	"slices"
)
import "hash/crc32"

// https://huizhou92.com/p/distributed-cornerstone-algorithm-consistent-hash/

var (
	ErrUnexpected = errors.New("unexpected error")
	ErrNotFound   = errors.New("not found")
)

type Hashing interface {
	AddServer(string)
	RemoveServer(string)
	Serve(request string) string
}

type Node struct {
	key string
	val string
}

type ring struct {
	size     int
	servers  []*Server
	hashFunc func(data []byte) uint32
}

type Server struct {
	id        string
	databases map[string]*Node
	hashVal   uint32
}

func NewServer(id string, hash uint32) *Server {
	return &Server{
		id:        id,
		databases: make(map[string]*Node),
		hashVal:   hash,
	}
}

func initRing(size int, servers []string) *ring {
	r := &ring{
		size:     size,
		servers:  make([]*Server, 0),
		hashFunc: crc32.ChecksumIEEE,
	}
	for _, serverId := range servers {
		server := NewServer(serverId, r.hashFunc([]byte(serverId)))
		r.servers = append(r.servers, server)
	}
	slices.SortFunc(r.servers, func(i, j *Server) int {
		return cmp.Compare(i.hashVal, j.hashVal)
	})
	return r
}

func (r *ring) hash(val int) int {
	return val % r.size
}

func (r *ring) AddServer(serverKey string) (int, error) {
	serverHashKey := r.hashFunc([]byte(serverKey))
	server := NewServer(serverKey, serverHashKey)
	r.servers = append(r.servers, server)
	slices.SortFunc(r.servers, func(i, j *Server) int {
		return cmp.Compare(i.hashVal, j.hashVal)
	})
	index, server, err := r.findServer(serverKey)
	if err != nil {
		return 0, ErrUnexpected
	}
	movedKeys := 0
	nextServer := r.servers[(index+1)%len(r.servers)]
	for _, node := range nextServer.databases {
		if r.hashFunc([]byte(node.key)) <= serverHashKey {
			server.databases[node.key] = node
			delete(nextServer.databases, node.key)
			movedKeys++
			log.Printf("move key %v from server %v to server %v\n", node.key, nextServer.id, serverKey)
		}
	}
	return movedKeys, nil
}

func (r *ring) RemoveServer(serverKey string) error {
	index, server, err := r.findServer(serverKey)
	if err != nil {
		return err
	}
	nextServer := r.servers[(index+1)%len(r.servers)]

	for key, node := range server.databases {
		nextServer.databases[key] = node
	}
	r.servers = append(r.servers[:index], r.servers[index+1:]...)
	return nil
}

func (r *ring) Store(key string, val string) (*Server, error) {
	_, server, err := r.findServer(key)
	if err != nil {
		return nil, err
	}
	server.databases[key] = &Node{key: key, val: val}
	return server, nil
}

func (r *ring) Get(key string) (string, error) {
	_, server, err := r.findServer(key)
	if err != nil {
		return "", err
	}
	node, ok := server.databases[key]
	if !ok {
		return "", ErrNotFound
	}
	return node.val, nil
}

func (r *ring) findServer(key string) (int, *Server, error) {
	keyHash := r.hashFunc([]byte(key))
	index, server := binarySearch(r.servers, keyHash)
	if server == nil {
		return index, server, ErrNotFound
	}
	return index, server, nil
}

// find the first server greater than target
func binarySearch(servers []*Server, target uint32) (int, *Server) {
	low := 0
	high := len(servers) - 1
	for low <= high {
		mid := (low + high) / 2
		serverHash := servers[mid].hashVal
		if target == serverHash {
			return mid, servers[mid]
		}
		// find left most
		if target < serverHash {
			high = mid - 1
		}
		if target > serverHash {
			low = mid + 1
		}
	}
	// if hash is between the last and first server
	if low >= len(servers) {
		return 0, servers[0]
	}
	return low, servers[low]
}
