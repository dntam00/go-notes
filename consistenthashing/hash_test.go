package consistenthashing

import (
	"fmt"
	"math"
	"testing"
)

func TestInitRing(t *testing.T) {
	size := math.MaxUint32
	servers := []string{"server1", "server2", "server3", "server4"}

	r := initRing(size, servers)

	if r.size != size {
		t.Errorf("Expected ring size to be %d, got %d", size, r.size)
	}

	if len(r.servers) != len(servers) {
		t.Errorf("Expected %d servers, got %d", len(servers), len(r.servers))
	}

	if len(r.servers) != len(servers) {
		t.Errorf("Expected %d servers, got %d", len(servers), len(r.servers))
	}
}

func TestStoreAndGet(t *testing.T) {
	size := math.MaxUint32
	servers := []string{"server1", "server2", "server3", "server4"}
	r := initRing(size, servers)

	_, err := r.Store("key_1", "value_1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Store("key_2", "value_2")
	if err != nil {
		t.Fatal(err)
	}
	get, err := r.Get("key_1")
	if err != nil || get != "value_1" {
		t.Errorf("expect `value_1`, get error: %v", get)
	}
}

type requestAndExpect struct {
	hash   uint32
	expect string
}

func TestBinarySearch(t *testing.T) {
	servers := []*Server{
		{
			id:        "server_1",
			databases: nil,
			hashVal:   10,
		},
		{
			id:        "server_2",
			databases: nil,
			hashVal:   100,
		},
		{
			id:        "server_3",
			databases: nil,
			hashVal:   1000,
		},
		{
			id:        "server_4",
			databases: nil,
			hashVal:   10000,
		},
	}

	targets := []requestAndExpect{
		{1, "server_1"},
		{9, "server_1"},
		{10, "server_1"},
		{100, "server_2"},
		{1_000, "server_3"},
		{9_999, "server_4"},
		{10_000, "server_4"},
		{10_001, "server_1"},
	}

	for _, target := range targets {
		_, server := binarySearch(servers, target.hash)
		if server.id != target.expect {
			t.Errorf("Expected %s, got %s", target.expect, server.id)
		}
	}
}

func TestAddServer(t *testing.T) {
	size := math.MaxUint32
	servers := []string{"server1", "server2", "server3", "server4"}
	r := initRing(size, servers)
	for _, server := range r.servers {
		fmt.Println(server.hashVal)
	}
	tobeAddedServer := "random_key_31"
	tobeAddedHash := r.hashFunc([]byte(tobeAddedServer))
	fmt.Println(tobeAddedHash)
	inRange := findKeyInRange(r, 5, r.servers[0].hashVal, tobeAddedHash)
	for _, v := range inRange {
		_, err := r.Store(v, v)
		if err != nil {
			t.Fatal(err)
		}
	}
	movedKey, err := r.AddServer(tobeAddedServer)
	if err != nil {
		t.Fatal(err)
	}
	if movedKey != 5 {
		t.Errorf("expect moved key to be 5, got %d", movedKey)
	}
	for _, v := range inRange {
		val, err := r.Get(v)
		if err != nil {
			t.Fatal(err)
		}
		if val != v {
			t.Errorf("Expected val to be %s, got %s", v, val)
		}
	}
}

func findKeyInRange(r *ring, size int, start, end uint32) []string {
	var result []string
	attempts := 0
	maxAttempts := size * 1_000_000_000

	for len(result) < size && attempts < maxAttempts {
		// Generate random string
		key := fmt.Sprintf("random_key_%d", attempts)

		// Calculate hash value for the key
		hashVal := r.hashFunc([]byte(key))

		// Check if hash is in range [start, end]
		var inRange bool
		if start <= end {
			inRange = hashVal > start && hashVal <= end
		} else {
			// Handle wrap-around case (e.g., start=4000000000, end=1000000000)
			inRange = hashVal >= start || hashVal <= end
		}

		if inRange {
			result = append(result, key)
		}

		attempts++
	}

	return result
}

func TestRemoveServer(t *testing.T) {
	size := math.MaxUint32
	servers := []string{"server1", "server2", "server3", "server4"}
	r := initRing(size, servers)

	_, _ = r.Store("key_1", "value_1")
	_, _ = r.Store("key_2", "value_2")
	_, _ = r.Store("key_3", "value_3")
	_, _ = r.Store("key_4", "value_4")
	_, _ = r.Store("key_5", "value_5")

	for _, server := range r.servers {
		if len(server.databases) > 0 {
			err := r.RemoveServer(server.id)
			if err != nil {
				t.Errorf("Remove server %s failed: %v", server.id, err)
			}
			break
		}
	}
	total := 0
	expectedKeySize := 5
	for _, server := range r.servers {
		total += len(server.databases)
	}
	if total != expectedKeySize {
		t.Errorf("Expected %d servers, got %d", expectedKeySize, total)
	}
}
