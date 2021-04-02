package delaunay

type node struct {
	i    int
	t    int
	prev *node
	next *node
}

func newNode(nodes []node, i int, prev *node) *node {
	n := &nodes[i]
	n.i = i
	if prev == nil {
		n.prev = n
		n.next = n
	} else {
		n.next = prev.next
		n.prev = prev
		prev.next.prev = n
		prev.next = n
	}
	return n
}

func (n *node) remove() *node {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.i = -1
	return n.prev
}
