package fiber

import (
	"engo/internal/protocol"
	"engo/pkg/scene"
)

type EffectTag uint32

const (
	EffectNone      EffectTag = 0
	EffectPlacement EffectTag = 1 << 0 // Node mới, cần Insert vào cha
	EffectUpdate    EffectTag = 1 << 1 // Props thay đổi, cần vẽ lại
	EffectDeletion  EffectTag = 1 << 2 // Node bị xóa
	EffectLayout    EffectTag = 1 << 3 // Cần tính lại Flexbox
)

type Fiber struct {
	Node *scene.Node

	// For compare if the fiber node order is change
	Index uint8

	Props any

	Parent  *Fiber
	Sibling *Fiber
	Child   *Fiber

	Alternate *Fiber
	Deletions []*Fiber

	Flags        EffectTag
	SubtreeFlags EffectTag

	Tag protocol.OpCode
	Key uint32

	LayoutX, LayoutY, LayoutW, LayoutH float32
}

// CreateWorkInProgress tạo ra (hoặc tái sử dụng) một Fiber cho cây WIP
// dựa trên Fiber hiện tại (Current).
func CreateWorkInProgress(current *Fiber, pendingProps interface{}) *Fiber {
	var workInProgress *Fiber

	if current.Alternate == nil {
		// A. Nếu chưa có bản sao -> Tạo mới
		// (Sau này nên dùng Object Pool ở đây thay vì new)
		workInProgress = &Fiber{
			Node:      current.Node,
			Tag:       current.Tag,
			Key:       current.Key,
			Alternate: current,
		}
		current.Alternate = workInProgress
	} else {
		// B. Nếu đã có bản sao -> Tái sử dụng (Reset dữ liệu cũ)
		workInProgress = current.Alternate

		// Reset các cờ hiệu và dữ liệu layout cũ
		workInProgress.Flags = EffectNone
		workInProgress.SubtreeFlags = EffectNone
		workInProgress.Sibling = nil
		workInProgress.Child = nil

		// Cập nhật dữ liệu mới nhất
		workInProgress.Node = current.Node
	}

	// Copy các props mới vào
	workInProgress.Props = pendingProps // Hoặc lấy từ current.Node nếu props nằm trong Node

	return workInProgress
}

// IsSameNode kiểm tra xem Fiber này và Node dữ liệu mới có tương thích không?
// (Dùng để quyết định Reuse hay đập đi xây lại)
func (f *Fiber) IsSameNode(node *scene.Node) bool {
	// 1. So sánh Key (Quan trọng nhất cho danh sách)
	if f.Key != node.ID {
		return false
	}

	// 2. So sánh Loại (Type/Tag)
	// Ví dụ: Không thể tái sử dụng Fiber của "Rect" cho "Text"
	// Giả sử bạn có hàm helper chuyển String Type sang OpCode Tag
	nodeTag := protocol.GetTagFromStr(node.Type)
	if f.Tag != nodeTag {
		return false
	}

	return true
}

func (f *Fiber) HasPropsChanged(newNode *scene.Node) bool {
	return newNode.Flags != scene.FlagNone
}

func (f *Fiber) ComputeMatrix(parentMatrix protocol.Matrix) {
	// 1. Lấy Local Matrix từ Scene Node (x, y, rotation, scale)
	// (Giả sử Node có method GetLocalMatrix)
	localMat := f.Node.GetLocalMatrix()

	// 2. Nhân với ma trận của cha
	// World = Parent * Local
	f.GlobalMatrix = parentMatrix.Multiply(localMat)

	// 3. Nếu node này là ScrollView, có thể cần nhân thêm ScrollOffset
	// ...
}

// ComputeLayout (Nếu dùng Flexbox)
func (f *Fiber) ComputeLayout() {
	// Nếu node này là Flex container, chạy thuật toán layout đơn giản
	// hoặc gọi tới Yoga engine để tính toán x,y,w,h cho con cái
}

func (f *Fiber) MarkUpdate() {
	f.Flags |= EffectUpdate
}

func (f *Fiber) MarkPlacement() {
	f.Flags |= EffectPlacement
}

func (f *Fiber) MarkDeletion() {
	f.Flags |= EffectDeletion
}

func (f *Fiber) BubbleFlags() {
	if f.Parent != nil {
		f.Parent.SubtreeFlags |= f.SubtreeFlags | f.Flags
	}
}
