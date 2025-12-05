package rtree

import "math"

type Node struct {
	IsLeaf  bool
	Entries []Entry
	Parent  *Node
}

// Cập nhật NewNode: Cần biết maxEntries để cấp phát slice capacity
func NewNode(isLeaf bool, maxEntries int) *Node {
	return &Node{
		IsLeaf: isLeaf,
		// Cấp phát dư 1 slot để chứa entry tạm trước khi split
		Entries: make([]Entry, 0, maxEntries+1),
	}
}

func (n *Node) computeMBR() Rect {
	if len(n.Entries) == 0 {
		return Rect{}
	}
	mbr := n.Entries[0].Mbr
	for _, e := range n.Entries[1:] {
		mbr = Union(mbr, e.Mbr)
	}
	return mbr
}

// Cập nhật isOverflow: Cần tham số maxEntries
func (n *Node) isOverflow(maxEntries int) bool {
	return len(n.Entries) > maxEntries
}

// Cập nhật split: Cần minEntries để chạy thuật toán Quadratic,
// và maxEntries để tạo newNode
func (n *Node) split(minEntries, maxEntries int) *Node {
	newNode := NewNode(n.IsLeaf, maxEntries)

	allEntries := append([]Entry{}, n.Entries...)
	n.Entries = n.Entries[:0]

	// --- QUADRATIC SPLIT ---

	seed1, seed2 := pickSeeds(allEntries)

	n.Entries = append(n.Entries, allEntries[seed1])
	newNode.Entries = append(newNode.Entries, allEntries[seed2])

	processed := make([]bool, len(allEntries))
	processed[seed1] = true
	processed[seed2] = true
	count := 2

	mbr1 := allEntries[seed1].Mbr
	mbr2 := allEntries[seed2].Mbr

	for count < len(allEntries) {
		remaining := len(allEntries) - count

		// SỬ DỤNG BIẾN minEntries THAY VÌ CONST
		if len(n.Entries)+remaining == minEntries {
			for i, e := range allEntries {
				if !processed[i] {
					n.Entries = append(n.Entries, e)
				}
			}
			break
		}
		if len(newNode.Entries)+remaining == minEntries {
			for i, e := range allEntries {
				if !processed[i] {
					newNode.Entries = append(newNode.Entries, e)
				}
			}
			break
		}

		idx := pickNext(allEntries, processed, mbr1, mbr2)
		entry := allEntries[idx]

		d1 := Enlargement(mbr1, entry.Mbr)
		d2 := Enlargement(mbr2, entry.Mbr)

		// Logic chọn nhóm (Giữ nguyên logic cũ)
		if d1 < d2 {
			n.Entries = append(n.Entries, entry)
			mbr1 = Union(mbr1, entry.Mbr)
		} else if d2 < d1 {
			newNode.Entries = append(newNode.Entries, entry)
			mbr2 = Union(mbr2, entry.Mbr)
		} else {
			area1 := mbr1.Area()
			area2 := mbr2.Area()
			if area1 < area2 {
				n.Entries = append(n.Entries, entry)
				mbr1 = Union(mbr1, entry.Mbr)
			} else if area2 < area1 {
				newNode.Entries = append(newNode.Entries, entry)
				mbr2 = Union(mbr2, entry.Mbr)
			} else {
				if len(n.Entries) <= len(newNode.Entries) {
					n.Entries = append(n.Entries, entry)
					mbr1 = Union(mbr1, entry.Mbr)
				} else {
					newNode.Entries = append(newNode.Entries, entry)
					mbr2 = Union(mbr2, entry.Mbr)
				}
			}
		}

		processed[idx] = true
		count++
	}

	if !n.IsLeaf {
		for _, entry := range newNode.Entries {
			if entry.Child != nil {
				entry.Child.Parent = newNode
			}
		}
	}

	return newNode
}

// pickSeeds tìm 2 entry lãng phí diện tích nhất nếu gộp chung (O(N^2))
func pickSeeds(entries []Entry) (int, int) {
	// Trường hợp edge case: ít hơn 2 phần tử (thường không xảy ra do logic insert)
	if len(entries) < 2 {
		return 0, 0
	}

	maxWaste := -1.0
	seed1, seed2 := 0, 1

	for i := range entries {
		for j := i + 1; j < len(entries); j++ {
			// Diện tích bao trùm cả 2
			unionArea := Union(entries[i].Mbr, entries[j].Mbr).Area()
			// Diện tích lãng phí = Diện tích bao trùm - tổng diện tích 2 hình
			waste := unionArea - entries[i].Mbr.Area() - entries[j].Mbr.Area()

			if waste > maxWaste {
				maxWaste = waste
				seed1 = i
				seed2 = j
			}
		}
	}
	return seed1, seed2
}

func pickNext(entries []Entry, processed []bool, mbr1, mbr2 Rect) int {
	maxDiff := -1.0
	chosenIdx := -1

	for i, e := range entries {
		if processed[i] {
			continue
		}

		d1 := Enlargement(mbr1, e.Mbr)
		d2 := Enlargement(mbr2, e.Mbr)
		diff := math.Abs(d1 - d2)

		if diff > maxDiff {
			maxDiff = diff
			chosenIdx = i
		}
	}

	// Fallback nếu không tìm thấy (trường hợp còn 1 phần tử)
	if chosenIdx == -1 {
		for i, p := range processed {
			if !p {
				return i
			}
		}
	}

	return chosenIdx
}

// removeEntryAt xóa entry tại index i và giữ thứ tự (quan trọng để debug)
// Hoặc có thể swap với phần tử cuối để nhanh hơn (O(1)) nếu không quan trọng thứ tự
func (n *Node) removeEntryAt(i int) {
	// Cách nhanh: Swap với cuối rồi cắt đuôi
	lastIdx := len(n.Entries) - 1
	n.Entries[i] = n.Entries[lastIdx]
	n.Entries = n.Entries[:lastIdx]
}
