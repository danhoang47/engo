package layout

type Direction int

const (
	DirectionRow Direction = iota
	DirectionColumn
)

type JustifyContent int

const (
	JustifyStart JustifyContent = iota
	JustifyEnd
	JustifyCenter

	// Put content into left/right edge
	JustifyBetween
	// Space of left + right = space between each item
	JustifyAround
	// Remaining space of flex container will divided by
	// space = remain space / num of item's row
	JustifyEvenly
)

type AlignItem int

const (
	// Stretch the item's height to fit the current
	// flex row height
	AlignStretch AlignItem = iota
	AlignCenter
	AlignStart
	AlignEnd
)
