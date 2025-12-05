package rtree

import (
	"math"
)

// Rect đại diện cho một hình chữ nhật 2D (Axis-Aligned Bounding Box)
type Rect struct {
	MinX, MinY float64
	MaxX, MaxY float64
}

// Entry là một phần tử trong Node.
// Nếu ở Leaf Node: Nó chứa ID dữ liệu (hoặc con trỏ data).
// Nếu ở Internal Node: Nó chứa MBR của node con và con trỏ tới node con đó.
type Entry struct {
	Mbr   Rect
	Data  interface{} // ID dữ liệu (cho Leaf)
	Child *Node       // Con trỏ tới Node con (cho Internal)
}

// Area tính diện tích hình chữ nhật
func (r Rect) Area() float64 {
	return (r.MaxX - r.MinX) * (r.MaxY - r.MinY)
}

// Union trả về hình chữ nhật bao trùm cả r1 và r2
func Union(r1, r2 Rect) Rect {
	return Rect{
		MinX: math.Min(r1.MinX, r2.MinX),
		MinY: math.Min(r1.MinY, r2.MinY),
		MaxX: math.Max(r1.MaxX, r2.MaxX),
		MaxY: math.Max(r1.MaxY, r2.MaxY),
	}
}

// Intersect kiểm tra hai hình chữ nhật có giao nhau không
func Intersect(r1, r2 Rect) bool {
	return r1.MinX <= r2.MaxX &&
		r1.MaxX >= r2.MinX &&
		r1.MinY <= r2.MaxY &&
		r1.MaxY >= r2.MinY
}

// Enlargement tính diện tích tăng thêm nếu thêm r2 vào r1
// Dùng cho thuật toán ChooseLeaf và Quadratic Split
func Enlargement(r1, r2 Rect) float64 {
	unionRect := Union(r1, r2)
	return unionRect.Area() - r1.Area()
}
