package math

import "math"

type Rect struct {
	Min Coord
	Max Coord
}

func NewRect(min, max Coord) *Rect {
	return &Rect{
		Min: min,
		Max: max,
	}
}

func (r *Rect) Area() float64 {
	return (r.Max.X - r.Min.X) * (r.Max.Y - r.Min.Y)
}

// Return minimum align-axis bounding rect
func (r *Rect) Union(other *Rect) *Rect {
	minX := math.Min(r.Min.X, other.Min.X)
	minY := math.Min(r.Min.Y, other.Min.Y)

	maxX := math.Max(r.Max.X, other.Max.X)
	maxY := math.Max(r.Max.Y, other.Max.Y)

	return &Rect{
		Min: Coord{minX, minY},
		Max: Coord{maxX, maxY},
	}
}

func (r *Rect) Contains(other *Rect) bool {
	if r.Min.X > other.Min.X ||
		r.Min.Y < other.Min.Y ||
		r.Max.X < other.Max.X ||
		r.Max.Y < other.Max.Y {
		return false
	}

	return true
}
