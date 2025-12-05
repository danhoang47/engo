package rtree

import "math"

// Default configuration nếu người dùng lười
const (
	DefaultMinEntries = 2
	DefaultMaxEntries = 4
)

type RTree struct {
	Root        *Node
	MinChildren int
	MaxChildren int
}

// Constructor mặc định
func NewRTree() *RTree {
	return NewRTreeWithConfig(DefaultMinEntries, DefaultMaxEntries)
}

// Constructor tùy chỉnh (Cái bạn cần)
func NewRTreeWithConfig(min, max int) *RTree {
	return &RTree{
		Root:        NewNode(true, max), // Truyền max để init capacity
		MinChildren: min,
		MaxChildren: max,
	}
}

// mbr stands for Minimum-Bounding Rect
// data will be scene-graph node
func (t *RTree) Insert(mbr Rect, data interface{}) {
	entry := Entry{
		Mbr:  mbr,
		Data: data,
	}

	leaf := t.chooseLeaf(t.Root, entry)
	leaf.Entries = append(leaf.Entries, entry)

	var splitNode *Node
	// Truyền MaxChildren vào để check overflow
	if leaf.isOverflow(t.MaxChildren) {
		// Truyền cấu hình vào hàm split
		splitNode = leaf.split(t.MinChildren, t.MaxChildren)
	}

	t.adjustTree(leaf, splitNode)
}

// Search tìm tất cả các object có MBR giao với searchMBR
func (t *RTree) Search(searchMBR Rect) []interface{} {
	var results []interface{}
	t.searchRecursive(t.Root, searchMBR, &results)
	return results
}

func (t *RTree) searchRecursive(n *Node, s Rect, results *[]interface{}) {
	for _, e := range n.Entries {
		if Intersect(s, e.Mbr) {
			if n.IsLeaf {
				*results = append(*results, e.Data)
			} else {
				t.searchRecursive(e.Child, s, results)
			}
		}
	}
}

func (t *RTree) adjustTree(n *Node, nn *Node) {
	if n == t.Root {
		if nn != nil {
			// Tạo root mới cần biết MaxChildren
			newRoot := NewNode(false, t.MaxChildren)
			newRoot.Entries = append(newRoot.Entries, Entry{Mbr: n.computeMBR(), Child: n})
			newRoot.Entries = append(newRoot.Entries, Entry{Mbr: nn.computeMBR(), Child: nn})
			n.Parent = newRoot
			nn.Parent = newRoot
			t.Root = newRoot
		}
		return
	}

	parent := n.Parent

	for i := range parent.Entries {
		if parent.Entries[i].Child == n {
			parent.Entries[i].Mbr = n.computeMBR()
			break
		}
	}

	if nn != nil {
		entryNN := Entry{Mbr: nn.computeMBR(), Child: nn}
		parent.Entries = append(parent.Entries, entryNN)
		nn.Parent = parent
	}

	var newParentSplit *Node
	// Truyền MaxChildren vào check overflow
	if parent.isOverflow(t.MaxChildren) {
		// Truyền cấu hình vào split
		newParentSplit = parent.split(t.MinChildren, t.MaxChildren)
	}

	t.adjustTree(parent, newParentSplit)
}

// chooseLeaf chọn node lá để chèn entry (theo tiêu chí mở rộng diện tích ít nhất)
func (t *RTree) chooseLeaf(n *Node, e Entry) *Node {
	if n.IsLeaf {
		return n
	}

	// Tìm entry con mà khi thêm e vào thì diện tích mở rộng là nhỏ nhất
	minEnlargement := math.MaxFloat64
	bestEntryIdx := -1

	for i, childEntry := range n.Entries {
		enlargement := Enlargement(childEntry.Mbr, e.Mbr)
		if enlargement < minEnlargement {
			minEnlargement = enlargement
			bestEntryIdx = i
		}

		if enlargement == minEnlargement {
			// Nếu bằng nhau, chọn cái có diện tích nhỏ hơn
			if childEntry.Mbr.Area() < n.Entries[bestEntryIdx].Mbr.Area() {
				bestEntryIdx = i
			}
		}
	}

	return t.chooseLeaf(n.Entries[bestEntryIdx].Child, e)
}

func (t *RTree) Update(oldMBR Rect, oldData interface{}, newMBR Rect, newData interface{}) bool {
	if t.Delete(oldMBR, oldData) {
		t.Insert(newMBR, newData)
		return true
	}
	return false
}

func (t *RTree) Delete(mbr Rect, data interface{}) bool {
	// 1. Tìm node lá chứa entry cần xóa
	leaf, idx := t.findLeaf(t.Root, mbr, data)
	if leaf == nil {
		return false // Không tìm thấy
	}

	// 2. Xóa entry khỏi node lá
	leaf.removeEntryAt(idx)

	// 3. CondenseTree: Xử lý underflow và lan truyền thay đổi lên gốc
	t.condenseTree(leaf)

	// 4. Giảm độ cao cây nếu Root chỉ còn 1 con (và con đó không phải là Leaf)
	if !t.Root.IsLeaf && len(t.Root.Entries) == 1 {
		t.Root = t.Root.Entries[0].Child
		t.Root.Parent = nil
	}

	return true
}

func (t *RTree) findLeaf(n *Node, mbr Rect, data interface{}) (*Node, int) {
	if n.IsLeaf {
		for i, e := range n.Entries {
			// So sánh MBR trước cho nhanh
			// Lưu ý: data interface{} cần so sánh được (==).
			// Nếu data là struct phức tạp, bạn cần logic so sánh riêng.
			if e.Data == data {
				return n, i
			}
		}
		return nil, -1
	}

	// Nếu là Internal Node, duyệt xuống các con có MBR bao trùm đối tượng
	for _, e := range n.Entries {
		if Intersect(e.Mbr, mbr) { // Dùng Intersect hoặc Contains đều được
			leaf, idx := t.findLeaf(e.Child, mbr, data)
			if leaf != nil {
				return leaf, idx
			}
		}
	}
	return nil, -1
}

// condenseTree xử lý việc tái cân bằng cây sau khi xóa
// Nguyên tắc: Nếu node bị thiếu hụt (Underflow), xóa node đó và Re-insert các con của nó.
func (t *RTree) condenseTree(n *Node) {
	var q []*Node // Danh sách các node bị xóa cần insert lại con
	currentNode := n

	// Duyệt từ node bị xóa ngược lên gốc
	for currentNode != t.Root {
		parent := currentNode.Parent

		// Tìm entry trỏ tới currentNode trong parent
		idx := -1
		for i, e := range parent.Entries {
			if e.Child == currentNode {
				idx = i
				break
			}
		}

		// Nếu bị Underflow (ít hơn MinChildren)
		if len(currentNode.Entries) < t.MinChildren {
			// 1. Xóa currentNode khỏi parent
			parent.removeEntryAt(idx)

			// 2. Đưa currentNode vào danh sách cần Re-insert
			q = append(q, currentNode)
		} else {
			// Nếu không underflow, chỉ cần cập nhật lại MBR của parent
			parent.Entries[idx].Mbr = currentNode.computeMBR()
		}

		currentNode = parent
	}

	// Re-insert các "trẻ mồ côi" (Orphans)
	// Lưu ý: Re-insert các entries bên trong node, chứ không phải bản thân node đó
	for _, node := range q {
		// Nếu node là lá, insert lại Data
		if node.IsLeaf {
			for _, e := range node.Entries {
				t.Insert(e.Mbr, e.Data)
			}
		} else {
			// Nếu node là internal, insert lại các nhánh con (Child Nodes)
			// Lưu ý: Level của cây phải được bảo toàn.
			// Để đơn giản hóa, ta re-insert đệ quy các entry lá của nhánh này.
			// (Cách Guttman chuẩn: Re-insert tại cùng level, nhưng phức tạp hơn.
			// Cách đơn giản và hiệu quả tương đương: Lấy hết tất cả Leaf bên dưới node này insert lại)
			leaves := t.collectAllLeaves(node)
			for _, e := range leaves {
				t.Insert(e.Mbr, e.Data)
			}
		}
	}
}

// collectAllLeaves gom tất cả entry lá con cháu của n (Helper cho Re-insert)
func (t *RTree) collectAllLeaves(n *Node) []Entry {
	var entries []Entry
	if n.IsLeaf {
		return n.Entries
	}
	for _, e := range n.Entries {
		entries = append(entries, t.collectAllLeaves(e.Child)...)
	}
	return entries
}
