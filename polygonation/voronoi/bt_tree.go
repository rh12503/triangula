// MIT License: See https://github.com/pzsz/voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Port of Raymond Hill's (rhill@raymondhill.net) javascript implementation 
// of Steven  Forune's algorithm to compute Voronoi diagrams

package voronoi

// Red-Black tree code (based on C version of "rbtree" by Franck Bui-Huu
// https://github.com/fbuihuu/libtree/blob/master/rb.c
type rbTree struct {
	root *rbNode
}

type rbNodeValue interface {
	bindToNode(node *rbNode)
	getNode() *rbNode
}

type rbNode struct {
	value    rbNodeValue
	left     *rbNode
	right    *rbNode
	parent   *rbNode
	previous *rbNode
	next     *rbNode
	red      bool
}

func (t *rbTree) insertSuccessor(node *rbNode, vsuccessor rbNodeValue) {
	successor := &rbNode{value: vsuccessor}
	vsuccessor.bindToNode(successor)

	var parent *rbNode
	if node != nil {
		// >>> rhill 2011-05-27: Performance: cache previous/next nodes
		successor.previous = node
		successor.next = node.next
		if node.next != nil {
			node.next.previous = successor
		}
		node.next = successor
		// <<<
		if node.right != nil {
			// in-place expansion of node.rbRight.getFirst()
			node = node.right
			for ; node.left != nil; node = node.left {
			}
			node.left = successor
		} else {
			node.right = successor
		}
		parent = node

		// rhill 2011-06-07: if node is null, successor must be inserted
		// to the left-most part of the tree
	} else if t.root != nil {
		node = t.getFirst(t.root)
		// >>> Performance: cache previous/next nodes
		successor.previous = nil
		successor.next = node
		node.previous = successor
		// <<<
		node.left = successor
		parent = node
	} else {
		// >>> Performance: cache previous/next nodes
		successor.previous = nil
		successor.next = nil
		// <<<
		t.root = successor
		parent = nil
	}
	successor.left = nil
	successor.right = nil
	successor.parent = parent
	successor.red = true
	// Fixup the modified tree by recoloring nodes and performing
	// rotations (2 at most) hence the red-black tree properties are
	// preserved.
	var grandpa, uncle *rbNode
	node = successor
	for parent != nil && parent.red {
		grandpa = parent.parent
		if parent == grandpa.left {
			uncle = grandpa.right
			if uncle != nil && uncle.red {
				parent.red = false
				uncle.red = false
				grandpa.red = true
				node = grandpa
			} else {
				if node == parent.right {
					t.rotateLeft(parent)
					node = parent
					parent = node.parent
				}
				parent.red = false
				grandpa.red = true
				t.rotateRight(grandpa)
			}
		} else {
			uncle = grandpa.left
			if uncle != nil && uncle.red {
				parent.red = false
				uncle.red = false
				grandpa.red = true
				node = grandpa
			} else {
				if node == parent.left {
					t.rotateRight(parent)
					node = parent
					parent = node.parent
				}
				parent.red = false
				grandpa.red = true
				t.rotateLeft(grandpa)
			}
		}
		parent = node.parent
	}
	t.root.red = false
}

func (t *rbTree) removeNode(node *rbNode) {
	// >>> rhill 2011-05-27: Performance: cache previous/next nodes
	if node.next != nil {
		node.next.previous = node.previous
	}
	if node.previous != nil {
		node.previous.next = node.next
	}
	node.next = nil
	node.previous = nil
	// <<<
	var parent = node.parent
	var left = node.left
	var right = node.right
	var next *rbNode
	if left == nil {
		next = right
	} else if right == nil {
		next = left
	} else {
		next = t.getFirst(right)
	}
	if parent != nil {
		if parent.left == node {
			parent.left = next
		} else {
			parent.right = next
		}
	} else {
		t.root = next
	}
	// enforce red-black rules
	isRed := false
	if left != nil && right != nil {
		isRed = next.red
		next.red = node.red
		next.left = left
		left.parent = next
		if next != right {
			parent = next.parent
			next.parent = node.parent
			node = next.right
			parent.left = node
			next.right = right
			right.parent = next
		} else {
			next.parent = parent
			parent = next
			node = next.right
		}
	} else {
		isRed = node.red
		node = next
	}
	// 'node' is now the sole successor's child and 'parent' its
	// new parent (since the successor can have been moved)
	if node != nil {
		node.parent = parent
	}
	// the 'easy' cases
	if isRed {
		return
	}
	if node != nil && node.red {
		node.red = false
		return
	}
	// the other cases
	var sibling *rbNode
	for {
		if node == t.root {
			break
		}
		if node == parent.left {
			sibling = parent.right
			if sibling.red {
				sibling.red = false
				parent.red = true
				t.rotateLeft(parent)
				sibling = parent.right
			}
			if (sibling.left != nil && sibling.left.red) || (sibling.right != nil && sibling.right.red) {
				if sibling.right == nil || !sibling.right.red {
					sibling.left.red = false
					sibling.red = true
					t.rotateRight(sibling)
					sibling = parent.right
				}
				sibling.red = parent.red
				parent.red = false
				sibling.right.red = false
				t.rotateLeft(parent)
				node = t.root
				break
			}
		} else {
			sibling = parent.left
			if sibling.red {
				sibling.red = false
				parent.red = true
				t.rotateRight(parent)
				sibling = parent.left
			}
			if (sibling.left != nil && sibling.left.red) || (sibling.right != nil && sibling.right.red) {
				if sibling.left == nil || !sibling.left.red {
					sibling.right.red = false
					sibling.red = true
					t.rotateLeft(sibling)
					sibling = parent.left
				}
				sibling.red = parent.red
				parent.red = false
				sibling.left.red = false
				t.rotateRight(parent)
				node = t.root
				break
			}
		}
		sibling.red = true
		node = parent
		parent = parent.parent
		if node.red {
			break
		}
	}
	if node != nil {
		node.red = false
	}
}

func (t *rbTree) rotateLeft(node *rbNode) {
	var p = node
	var q = node.right // can't be null
	var parent = p.parent
	if parent != nil {
		if parent.left == p {
			parent.left = q
		} else {
			parent.right = q
		}
	} else {
		t.root = q
	}
	q.parent = parent
	p.parent = q
	p.right = q.left
	if p.right != nil {
		p.right.parent = p
	}
	q.left = p
}

func (t *rbTree) rotateRight(node *rbNode) {
	var p = node
	var q = node.left // can't be null
	var parent = p.parent
	if parent != nil {
		if parent.left == p {
			parent.left = q
		} else {
			parent.right = q
		}
	} else {
		t.root = q
	}
	q.parent = parent
	p.parent = q
	p.left = q.right
	if p.left != nil {
		p.left.parent = p
	}
	q.right = p
}

func (t *rbTree) getFirst(node *rbNode) *rbNode {
	for node.left != nil {
		node = node.left
	}
	return node
}

func (t *rbTree) getLast(node *rbNode) *rbNode {
	for node.right != nil {
		node = node.right
	}
	return node
}
