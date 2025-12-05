package layout

import (
	"fmt"
	"strconv"
	"strings"
)

// Enum định nghĩa loại đơn vị
type Unit uint8

const (
	UnitPixel   Unit = 0 // Giá trị tuyệt đối (px)
	UnitPercent Unit = 1 // Giá trị tương đối (%)
	UnitAuto    Unit = 2 // Tự động tính toán
	// Có thể thêm UnitUndefined nếu cần
)

// Dimension: Struct đại diện cho Width, Height, Margin, Padding...
// Kích thước: 4 bytes (float32) + 1 byte (uint8) + 3 bytes (padding) = 8 bytes
type Dimension struct {
	Value float32
	Unit  Unit
}

// --- Các Helper Constructors ---

func Px(v float32) Dimension {
	return Dimension{Value: v, Unit: UnitPixel}
}

func Pct(v float32) Dimension {
	return Dimension{Value: v, Unit: UnitPercent}
}

func Auto() Dimension {
	return Dimension{Value: 0, Unit: UnitAuto} // Value không quan trọng khi Auto
}

// Resolve chuyển đổi Dimension thành số pixel cụ thể.
// parentSize: Kích thước tham chiếu của cha (Width hoặc Height).
func (d Dimension) Resolve(parentSize float32) float32 {
	switch d.Unit {
	case UnitPixel:
		return d.Value
	case UnitPercent:
		return (d.Value / 100.0) * parentSize
	case UnitAuto:
		// Với Auto, thường Layout Engine sẽ xử lý riêng.
		// Nhưng nếu cần số, Auto thường mặc định là 0 hoặc content-size.
		// FIXME: should consider about this
		return 0
	}
	return 0
}

// UnmarshalJSON: Parse dữ liệu từ JS/File
// Hỗ trợ: 100 (số), "100px", "50%", "auto"
func (d *Dimension) UnmarshalJSON(data []byte) error {
	str := string(data)

	// Trường hợp 1: "auto"
	if strings.Contains(str, "auto") {
		d.Unit = UnitAuto
		d.Value = 0
		return nil
	}

	// Trường hợp 2: Percentage ("50%")
	if strings.Contains(str, "%") {
		d.Unit = UnitPercent
		valStr := strings.Trim(str, `"%`) // Bỏ dấu ngoặc kép và %
		val, err := strconv.ParseFloat(valStr, 32)
		if err != nil {
			return err
		}
		d.Value = float32(val)
		return nil
	}

	// Trường hợp 3: Number (100 hoặc "100px")
	// Xử lý chuỗi JSON số (không có ngoặc kép) hoặc chuỗi có "px"
	cleanStr := strings.Trim(str, `"px`)
	val, err := strconv.ParseFloat(cleanStr, 32)
	if err != nil {
		return err
	}

	d.Unit = UnitPixel
	d.Value = float32(val)
	return nil
}

// MarshalJSON: Xuất ra JSON
func (d Dimension) MarshalJSON() ([]byte, error) {
	switch d.Unit {
	case UnitAuto:
		return []byte(`"auto"`), nil
	case UnitPercent:
		s := fmt.Sprintf(`"%.2f%%"`, d.Value)
		return []byte(s), nil
	case UnitPixel:
		// Xuất ra số thuần túy để tiết kiệm dung lượng
		s := fmt.Sprintf(`%.2f`, d.Value)
		return []byte(s), nil
	}
	return []byte("0"), nil
}
