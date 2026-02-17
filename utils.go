package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func StrokeLine(screen *ebiten.Image, p1, p2 *Vector2, w float32, c color.Color) {
	vector.StrokeLine(
		screen,
		float32(p1.X), float32(p1.Y), float32(p2.X), float32(p2.Y),
		w, c, false,
	)
}
