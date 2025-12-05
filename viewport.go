package main

type Viewport struct {
	Scale      float32
	TranslateX float32
	TranslateY float32
}

const (
	DefaultScale      = 1.0
	DefaultTranslateX = 0.0
	DefaultTranslateY = 0.0
)

func NewViewport() *Viewport {
	return &Viewport{
		Scale:      DefaultScale,
		TranslateX: DefaultTranslateX,
		TranslateY: DefaultTranslateY,
	}
}

func NewViewPortWithArgs(scale, translateX, translateY float32) *Viewport {
	return &Viewport{
		Scale:      scale,
		TranslateX: translateX,
		TranslateY: translateY,
	}
}
