package layout

// Flex properties will be calculated BASED on
// current direction main axis

// struct for flex layout property
// TODO: Using bitmap for this
type Flex struct {
	Direction      Direction
	JustifyContent JustifyContent
	AlignItem      AlignItem
	Grow           int8
	Shrink         int8
}
