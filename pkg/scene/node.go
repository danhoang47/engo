package scene

type NodeType int

const (
	ElementNode NodeType = iota
	TextNode
)

type Node struct {
	ID uint32

	Parent   *Node
	Children []*Node
	Type     NodeType
}

// Manipulate child nodes functions
// TODO: Construct a Fiber tree after changes
func (n *Node) InsertBefore(newChild, beforeChild *Node) {
	for i, c := range n.Children {
		if c == beforeChild {
			newChild.Parent = n
			n.Children = append(
				n.Children[:i], append([]*Node{newChild}, n.Children[i:]...)...)
			return
		}
	}
}

func (n *Node) AppendChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

func (n *Node) RemoveChild(child *Node) {
	for i, c := range n.Children {
		if c == child {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			child.Parent = nil
			return
		}
	}
}

func (n *Node) ReplaceChild(oldChild, newChild *Node) bool {
	for i, c := range n.Children {
		if c == oldChild {
			newChild.Parent = n
			oldChild.Parent = nil
			n.Children[i] = newChild
			return true
		}
	}
	return false
}

// Contains reports whether node is a descendant of n.
// Note: this will cause a stack overflow if the tree contains cycles.
func (n *Node) Contains(node *Node) bool {
	if n == node {
		return true
	}

	for _, child := range n.Children {
		if child.Contains(node) {
			return true
		}
	}

	return false
}

func (n *Node) Clone() *Node {
	clone := &Node{
		// the ID should not the same
		// maybe will add a global incremental ID generator later
		ID:   n.ID + 1,
		Type: n.Type,
	}

	for _, child := range n.Children {
		childClone := child.Clone()
		childClone.Parent = clone
		clone.Children = append(clone.Children, childClone)
	}

	return clone
}

func (n *Node) HasChildNodes() bool {
	return len(n.Children) > 0
}

// Implement EventTarget interface
// TODO: Implement
func (n *Node) DispatchEvent(event Event) {

}

// TODO: Implement
func (n *Node) AddEventListener(eventType EventType, listener EventListener) {

}
