package protocol

type OpCode uint32

const (
	// --- GROUP 0: META & CONTROL ---
	OpEof       OpCode = 0x00 // Kết thúc frame
	OpClear     OpCode = 0x01 // Xóa màn hình
	OpPushGroup OpCode = 0x02 // Bắt đầu Layer/Group (Save state)
	OpPopGroup  OpCode = 0x03 // Kết thúc Layer/Group (Restore state)

	// --- GROUP 1: TRANSFORM & CLIPPING ---
	OpSetMatrix OpCode = 0x10 // Set ma trận tuyệt đối (a, b, c, d, tx, ty)
	OpTransform OpCode = 0x11 // Nhân ma trận (Relative)
	OpClipRect  OpCode = 0x12 // Cắt vùng nhìn hình chữ nhật
	OpResetClip OpCode = 0x13 // Hủy cắt

	// --- GROUP 2: STATE MANAGEMENT ---
	OpSetFill   OpCode = 0x20 // Set màu nền (Solid Color)
	OpSetStroke OpCode = 0x21 // Set màu viền + độ dày
	OpSetJoin   OpCode = 0x22 // Set kiểu nối góc (Miter/Round/Bevel)
	OpSetDash   OpCode = 0x23 // Set nét đứt
	OpSetShadow OpCode = 0x24 // Set đổ bóng

	// --- GROUP 3: PRIMITIVES ---
	OpDrawRect  OpCode = 0x30 // Vẽ hình chữ nhật
	OpDrawRRect OpCode = 0x31 // Vẽ hình chữ nhật bo góc (Rounded)
	OpDrawOval  OpCode = 0x32 // Vẽ hình tròn/bầu dục
	OpDrawLine  OpCode = 0x33 // Vẽ đường thẳng

	// --- GROUP 4: PATHS (Vector) ---
	OpPathBegin  OpCode = 0x40 // Bắt đầu Path mới
	OpPathMove   OpCode = 0x41 // Di chuyển điểm bút (MoveTo)
	OpPathLine   OpCode = 0x42 // Kẻ đường thẳng đến (LineTo)
	OpPathQuad   OpCode = 0x43 // Đường cong bậc 2 (Quadratic)
	OpPathCubic  OpCode = 0x44 // Đường cong bậc 3 (Cubic)
	OpPathClose  OpCode = 0x45 // Khép kín Path
	OpPathFill   OpCode = 0x46 // Tô màu Path
	OpPathStroke OpCode = 0x47 // Viền Path

	// --- GROUP 5: ASSETS (Text & Images) ---
	OpDrawImg  OpCode = 0x50 // Vẽ ảnh thường
	OpDrawImg9 OpCode = 0x51 // Vẽ ảnh 9-slice (cho UI buttons)
	OpSetFont  OpCode = 0x52 // Cài đặt Font chữ + Size
	OpDrawText OpCode = 0x53 // Vẽ Text (theo ID chuỗi)
)
