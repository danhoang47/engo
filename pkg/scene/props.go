package scene

type Props interface {
	Clone() Props
}

type RectProps struct {
	CornerRadius [4]float32
	Fill         uint32 // Màu RGBA
	Stroke       uint32
	StrokeWidth  float32
}

type TextProps struct {
	Content    string
	FontFamily string
	FontSize   float32
	Fill       uint32
	Align      uint8
}

type ImageProps struct {
	SourceURL string
	TextureID uint32
	ScaleMode uint8
}
