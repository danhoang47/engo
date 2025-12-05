package engine

import (
	"engo/internal/algo/rtree"
	"engo/internal/protocol"
	"engo/pkg/fiber"
	"engo/pkg/layout"
	"engo/pkg/scene"
)

var id = 0

type Engine struct {
	commandBuffer *protocol.CommandBuffer
	// viewport      *Viewport
	spatial  *rtree.RTree
	rootNode *scene.Node
	// current selected node, rootNode if nil
	selection  []*scene.Node
	reconciler *fiber.Reconciler
}

var engine *Engine

//go:export
func Init() {
	rootNode := scene.NewNode(scene.Page, nil)

	engine = &Engine{
		commandBuffer: protocol.NewCommandBuffer(),
		// viewport:      NewViewport(),
		spatial:    rtree.NewRTreeWithConfig(5, 10),
		reconciler: fiber.NewReconciler(rootNode, nil),
	}

	engine.reconciler.ScheduleUpdate(rootNode)
}

// TODO: Implement this, really important
func InitWithJSON() {

}

//go:export
func Resize(width, height, dpr float32) {

}

// Lấy địa chỉ con trỏ của Shared Memory (Input State)
// Để JS biết chỗ mà ghi tọa độ chuột vào
//
// export
func GetInputStatePtr() uintptr {

}

// Chạy logic tính toán (Layout, Physics, Animation)
// dt: Delta time (giây)
//
// export
func Update(dt float32)

// Lấy địa chỉ bắt đầu của Command Buffer (Mảng uint32 chứa OpCode)
//
//go:export
func GetRenderBufferPtr() uintptr {
	return uintptr(engine.commandBuffer.GetPtr())
}

// Lấy kích thước hiện tại của Buffer (số lượng phần tử uint32)
func GetRenderBufferSize() int32 {
	return int32(engine.commandBuffer.GetSize())
}

// Xử lý Click/Mousedown/Up
// button: 0 (Left), 1 (Middle), 2 (Right)
// action: 0 (Down), 1 (Up)
// modifiers: Bitmask (Ctrl, Shift, Alt)
func OnMouseAction(button int32, action int32, modifiers int32)

// Xử lý bàn phím
// key_code: Mã ASCII hoặc KeyCode của JS
func OnKeyAction(keyCode int32, action int32, modifiers int32)

// Xử lý Zoom/Pan (nếu không dùng Shared Memory cho cái này)
func OnWheel(deltaX, deltaY float32, isZoom bool)

func CreateNode(nodeType scene.NodeType, budgetMs int64) uint32 {
	parentNode := engine.GetInsertTarget()

	parentNode.MarkDirty()

	// Add node to scene graph
	node := scene.NewNode(nodeType, parentNode)

	// Set current fiber root to parentNode
	engine.reconciler.ScheduleUpdate(parentNode)

	hasMoreWork := engine.reconciler.WorkLoop(budgetMs)

	if hasMoreWork {

	} else {

	}

	engine.commandBuffer.WriteEof()

	engine.SetSelection(node, false)
}

// Xóa Node
func RemoveNode(id uint32)

func UpdateNodeType(id uint32) {

}

// --- GENERIC SETTERS (Nhanh hơn truyền string) ---

// Set thuộc tính số thực (X, Y, W, H, Opacity, Radius...)
// propCode: Enum (1=X, 2=Y, 3=W, 4=H...)
func SetFloatProp(id uint32, propCode int32, value float32)

// Set thuộc tính số nguyên (Color, Visibility, Z-Index)
// value: Chứa cả màu RGBA nén lại hoặc Enum ID
func SetIntProp(id uint32, propCode int32, value int32)

// Set chuỗi (Text Content, Name)
// Vì truyền string qua Wasm tốn kém, ta truyền ptr và length
func SetStringProp(id uint32, propCode int32, strPtr uint32, len int32)

// Báo cho Go biết đã load xong ảnh
// JS gửi kích thước thật để Go tính layout
func RegisterImage(imgID uint32, width float32, height float32)

// Đăng ký Font metrics (để Go tính đo độ rộng chữ)
// Go cần biết Ascent, Descent, AdvanceWidth trung bình...
func RegisterFont(fontID uint32, ptrMetrics uint32)

// Lấy cây Scene Graph dạng JSON (để save file hoặc hiện Layer Tree bên React)
// Hàm này trả về pointer tới vùng nhớ chứa chuỗi JSON
func ExportJSON() uintptr

// Lấy độ dài chuỗi JSON
func GetJSONSize() int32

// Nhận một mảng các thay đổi: [NodeID, PropID, Value, NodeID, PropID, Value...]
func BatchUpdateFloats(ptr uint32, count int32)

//go:export setDimensionProp
func SetDimensionProp(nodeID uint32, propID int32, value float32, unit int32) {
	// Maybe page node
	node := engine.selection[0]

	dim := layout.Dimension{
		Value: value,
		Unit:  layout.Unit(unit),
	}

	switch propID {
	case PROP_WIDTH:
		node.Style.Width = dim
	case PROP_HEIGHT:
		node.Style.Height = dim
	}

	node.MarkDirty()
}
