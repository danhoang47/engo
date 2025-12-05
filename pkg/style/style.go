package style

// TODO: Convert all of these into bit fields
type Position int

const (
	PositionRelative Position = iota
	PositionAbsolute
	PositionSticky
	PositionFixed
)

type Overflow int

const (
	OverflowVisible Overflow = iota
	OverflowHidden
)

// This should hold both numeric, percentage and auto value
// TODO: should use this
type Length string

type Style struct {
	// Dimensions
	Height, Width float32

	// Positioned
	Top, Right, Bottom, Left float32

	Margin, Padding float32
	BorderWidth     float32
	Position        Position
	Overflow        Overflow
	Opacity         float32
	// TODO: Add transform
}

type Align int

const (
	AlignCenter Align = iota
	AlignLeft
	AlignRight
	AlignTop
	AlignBottom
)

type FontStyle struct {
	FontSize   uint8
	FontFamily string
	TextAlign  Align
}

func NewStyle() *Style {
	return nil
}
