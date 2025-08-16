package consistenthashing

import "errors"

var (
	ErrUnexpected = errors.New("unexpected error")
	ErrNotFound   = errors.New("not found")
)

type Hashing interface {
	AddNode(string)
	RemoveNode(string)
}

type node struct {
	val int
}

type ring struct {
	size     int
	database map[string]map[int][]node
	servers  map[string][]int
}

func initRing(size int, servers []string) *ring {
	r := &ring{
		size:     size,
		database: make(map[string]map[int][]node),
		servers:  make(map[string][]int),
	}
	d := size / len(servers)
	from := 0
	for i := 0; i < len(servers); i++ {
		r.database[servers[i]] = make(map[int][]node)
		next := from + d
		if i == len(servers)-1 {
			next = size
		}
		r.servers[servers[i]] = []int{from, next}
		from = from + d
	}
	return r
}

func (r *ring) hash(val int) int {
	return val % r.size
}

func (r *ring) serve(val int) error {
	hashedVal := r.hash(val)
	selected := ""
	for server, rangesVal := range r.servers {
		if hashedVal >= rangesVal[0] && hashedVal < rangesVal[1] {
			selected = server
			break
		}
	}
	if selected == "" {
		return ErrUnexpected
	}
	serverDb, ok := r.database[selected]
	if !ok {
		return ErrUnexpected
	}

	n := node{val: val}

	newNodes := make([]node, 0)

	nodes, ok := serverDb[hashedVal]
	if ok {
		newNodes = nodes
	}
	newNodes = append(newNodes, n)
	serverDb[hashedVal] = newNodes
	return nil
}

func (r *ring) get(val int) (*node, error) {
	hashedVal := r.hash(val)
	selected := ""
	for server, rangesVal := range r.servers {
		if hashedVal >= rangesVal[0] && hashedVal < rangesVal[1] {
			selected = server
			break
		}
	}
	if selected == "" {
		return nil, ErrUnexpected
	}

	datas, ok := r.database[selected]
	if !ok {
		return nil, ErrUnexpected
	}

	nodes, ok := datas[hashedVal]
	if !ok {
		return nil, ErrNotFound
	}

	for _, node := range nodes {
		if node.val == val {
			return &node, nil
		}
	}
	return nil, ErrNotFound
}
