package scene

import "engo/pkg/style"

type NodeType int

const (
	// Like Figma, one page has one and only one Page node
	Page NodeType = iota
	Frame
	// Placeholder, no-op node, for semantic grouping purpose
	Section
	// TODO: Node type for vector
	Vector
	// Placeholder, use to group nodes inside
	Group
	// Rect, circle, ellipse, etc...
	Polygon
	Line
	Text
	Image
	// Component-based system
	Component
	Instance
	// Conjunction, direction
	Connector
	// For note purpose
	Sticky
)

// Each type will have difference behavior
// Considering moving this into different files
type HTMLElementType int

const (
	Div HTMLElementType = iota
	Button
	Anchor
	Nav
	List
	ListItem
	Checkbox
	Radio
	Toggle
	Input
)

type NodeFlag uint32

const (
	FlagNone           NodeFlag = 0
	FlagTransformDirty NodeFlag = 1 << 0
	FlagContentDirty   NodeFlag = 1 << 1
	FlagLayoutDirty    NodeFlag = 1 << 2
	FlagSubtreeDirty   NodeFlag = 1 << 3
)

// TODO: Built-in node from element types

type Node struct {
	ID uint32

	Parent          *Node
	Children        []*Node
	Type            NodeType
	HTMLElementType HTMLElementType

	Style         *style.Style
	ComputedStyle *style.Style

	Flags NodeFlag
}

var counter uint32 = 0

func NewNode(nodeType NodeType, parent *Node) *Node {
	// Increase ID
	counter = counter + 1

	return &Node{
		Parent:          parent,
		Type:            nodeType,
		HTMLElementType: Div,
		Flags:           FlagContentDirty,
	}
}

// TODO: Implement
func NewNodeWithHTMLTag() {

}

func (n *Node) MarkDirty(flag NodeFlag) {
	if n.Flags&flag != 0 {
		return
	}

	n.Flags |= flag

	if flag == FlagTransformDirty || flag == FlagLayoutDirty {
		n.bubbleUp()
	}
}

func (n *Node) bubbleUp() {
	curr := n.Parent
	for curr != nil {
		if curr.Flags&FlagSubtreeDirty != 0 {
			break
		}

		curr.Flags |= FlagSubtreeDirty
		curr = curr.Parent
	}
}

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

func (n *Node) IsContainer() bool {
	return n.Type == Group || n.Type == Frame || n.Type == Page
}
