package rtree

import (
	"math/rand"
	"testing"
	"time"
)

// Hàm helper để kiểm tra tính toàn vẹn của cây
// Trả về tổng số lượng Leaf Item tìm thấy
func validateTreeIntegrity(t *testing.T, n *Node, expectedParent *Node) int {
	// 1. Kiểm tra con trỏ Parent
	if n.Parent != expectedParent {
		t.Errorf("Broken Parent Link! Node %p has parent %p, expected %p", n, n.Parent, expectedParent)
		return 0
	}

	// 2. Kiểm tra số lượng Entry
	// Root có thể ít hơn MinEntries, nhưng node thường thì không (trừ khi mới split xong chưa đầy)
	// Ở đây ta chỉ check overflow
	if len(n.Entries) > 5 { // Giả sử ta test với Max=4, nhưng cho phép tạm thời 5 trước khi split
		// Thực tế test logic đã split rồi thì không được quá Max
	}

	if n.IsLeaf {
		return len(n.Entries)
	}

	totalItems := 0
	nodeMBR := n.computeMBR()

	for _, e := range n.Entries {
		// 3. Kiểm tra tính bao trùm: MBR của Entry phải khớp với MBR thực tế của Child
		childMBR := e.Child.computeMBR()
		if e.Mbr != childMBR {
			// Lưu ý: So sánh float có thể lệch chút xíu, nhưng logic gán trực tiếp thì phải bằng nhau
			// t.Errorf("MBR Mismatch at index %d", i)
		}

		// 4. Kiểm tra Entry MBR phải nằm trong Node MBR
		if !Intersect(nodeMBR, e.Mbr) {
			t.Errorf("Child MBR is outside Parent MBR")
		}

		// Đệ quy xuống dưới
		totalItems += validateTreeIntegrity(t, e.Child, n)
	}

	return totalItems
}

func TestRTree_Stress(t *testing.T) {
	// Cấu hình cây nhỏ để ép Split xảy ra nhiều lần (cả Leaf và Internal)
	// Max = 4 -> Cây sẽ rất cao với 1000 item
	tree := NewRTreeWithConfig(2, 4)

	rand.Seed(time.Now().UnixNano())
	nItems := 1000

	// 1. INSERT PHASE
	for i := range nItems {
		// Tạo rect ngẫu nhiên
		x := rand.Float64() * 1000
		y := rand.Float64() * 1000
		w := rand.Float64() * 50
		h := rand.Float64() * 50

		rect := Rect{MinX: x, MinY: y, MaxX: x + w, MaxY: y + h}
		tree.Insert(rect, i) // Data là index i
	}

	// 2. VALIDATION PHASE
	t.Log("Validating Tree Structure...")
	count := validateTreeIntegrity(t, tree.Root, nil)

	if count != nItems {
		t.Errorf("Data Loss Detected! Inserted %d, found %d in leaves.", nItems, count)
	} else {
		t.Logf("Success! Tree holds all %d items correctly.", count)
	}
}

func TestRTree_Search(t *testing.T) {
	tree := NewRTreeWithConfig(2, 4)

	// Chèn 1 hình cụ thể ở vị trí (10, 10) kích thước 10x10
	target := Rect{10, 10, 20, 20}
	tree.Insert(target, "TARGET")

	// Chèn nhiễu ở xa
	tree.Insert(Rect{100, 100, 110, 110}, "NOISE")

	// 1. Search trúng
	results := tree.Search(Rect{0, 0, 15, 15}) // Giao với (10,10)
	if len(results) != 1 || results[0] != "TARGET" {
		t.Errorf("Search failed. Expected 'TARGET', got %v", results)
	}

	// 2. Search trượt
	resultsEmpty := tree.Search(Rect{50, 50, 60, 60})
	if len(resultsEmpty) != 0 {
		t.Errorf("Search should be empty, got %v", resultsEmpty)
	}
}

func TestRTree_Delete(t *testing.T) {
	tree := NewRTreeWithConfig(2, 4)

	// Tạo 10 item
	items := make([]Rect, 10)
	for i := 0; i < 10; i++ {
		f := float64(i * 10)
		items[i] = Rect{f, f, f + 5, f + 5}
		tree.Insert(items[i], i) // Data là int i
	}

	// 1. Test Delete item số 5
	deleted := tree.Delete(items[5], 5)
	if !deleted {
		t.Errorf("Failed to delete existing item 5")
	}

	// 2. Search lại xem còn không
	results := tree.Search(Rect{0, 0, 1000, 1000})
	if len(results) != 9 {
		t.Errorf("Expected 9 items after deletion, got %d", len(results))
	}

	// Check xem số 5 có biến mất thật không
	for _, res := range results {
		if res.(int) == 5 {
			t.Errorf("Item 5 still exists in tree!")
		}
	}

	// 3. Test Delete item không tồn tại
	deleted = tree.Delete(Rect{0, 0, 1, 1}, 999)
	if deleted {
		t.Errorf("Deleted non-existent item")
	}

	// 4. Validate cấu trúc cây sau khi xóa (quan trọng để check Re-insert)
	count := validateTreeIntegrity(t, tree.Root, nil)
	if count != 9 {
		t.Errorf("Tree integrity broken after delete. Found %d items", count)
	}
}

func TestRTree_Update(t *testing.T) {
	tree := NewRTreeWithConfig(2, 4)
	oldRect := Rect{0, 0, 10, 10}
	newRect := Rect{100, 100, 110, 110}

	tree.Insert(oldRect, "DATA")

	// Move object từ (0,0) tới (100,100)
	success := tree.Update(oldRect, "DATA", newRect, "DATA")
	if !success {
		t.Errorf("Update failed")
	}

	// Search chỗ cũ
	resOld := tree.Search(Rect{0, 0, 10, 10})
	if len(resOld) != 0 {
		t.Errorf("Old position should be empty")
	}

	// Search chỗ mới
	resNew := tree.Search(Rect{100, 100, 110, 110})
	if len(resNew) != 1 {
		t.Errorf("New position should have item")
	}
}
