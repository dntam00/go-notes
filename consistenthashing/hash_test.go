package consistenthashing

import "testing"

func TestInitRing(t *testing.T) {
	// Test data
	size := 100
	servers := []string{"server1", "server2", "server3", "server4"}

	// Call the fixed function
	r := initRing(size, servers)

	// Verify ring is initialized correctly
	if r.size != size {
		t.Errorf("Expected ring size to be %d, got %d", size, r.size)
	}

	// Verify all servers are added
	if len(r.servers) != len(servers) {
		t.Errorf("Expected %d servers, got %d", len(servers), len(r.servers))
	}

	// Verify server ranges
	d := size / len(servers)
	from := 0
	for i, server := range servers {
		expected := []int{from, from + d}
		actual := r.servers[server]

		if len(actual) != 2 || actual[0] != expected[0] || actual[1] != expected[1] {
			t.Errorf("Server %d (%s): expected range %v, got %v",
				i, server, expected, actual)
		}

		from = from + d
	}
}

func TestRingServeAndGet(t *testing.T) {
	// Initialize the ring
	size := 100
	servers := []string{"server1", "server2", "server3", "server4"}
	r := initRing(size, servers)

	// Test values to add to the ring
	testValues := []int{5, 25, 55, 75, 95}

	// Test serving values to the ring
	for _, val := range testValues {
		err := r.serve(val)
		if err != nil {
			t.Errorf("Failed to serve value %d: %v", val, err)
		}
	}

	// Test getting values from the ring
	for _, val := range testValues {
		node, err := r.get(val)
		if err != nil {
			t.Errorf("Failed to get value %d: %v", val, err)
		}
		if node.val != val {
			t.Errorf("Expected node value %d, got %d", val, node.val)
		}
	}

	// Test getting a non-existent value
	nonExistentVal := 999
	_, err := r.get(nonExistentVal)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound for non-existent value, got %v", err)
	}

	// Test value distribution across servers
	serverCount := make(map[string]int)
	for i := 0; i < size; i++ {
		// Serve each value in the range
		err := r.serve(i)
		if err != nil {
			t.Errorf("Failed to serve value %d: %v", i, err)
		}

		// Find which server it went to
		hashedVal := r.hash(i)
		for server, rangeVals := range r.servers {
			if hashedVal >= rangeVals[0] && hashedVal < rangeVals[1] {
				serverCount[server]++
				break
			}
		}
	}

	// Verify each server got approximately the expected number of values
	expectedPerServer := size / len(servers)
	for server, count := range serverCount {
		if count < expectedPerServer-5 || count > expectedPerServer+5 {
			t.Errorf("Server %s: expected ~%d values, got %d",
				server, expectedPerServer, count)
		}
	}
}
