package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	red = color.RGBA{255, 0, 0, 255}
	blue = color.RGBA{0, 0, 255, 255}
)

var (
	x1 = &Vertex{-1000, 0, 0}
	x2 = &Vertex{1000, 0, 0}
	y1 = &Vertex{0, 1000, 0}
	y2 = &Vertex{0, -1000, 0}
	z1 = &Vertex{0, 0, -1000}
	z2 = &Vertex{0, 0, 1000}
)

func (g *Game) DrawWorldAxis(screen *ebiten.Image) {
	w := float32(4)

	x1p := VertexToGameScreen(x1, g, screen)
	x2p := VertexToGameScreen(x2, g, screen)
	StrokeLine(screen, x1p, x2p, w, red)

	y1p := VertexToGameScreen(y1, g, screen)
	y2p := VertexToGameScreen(y2, g, screen)
	StrokeLine(screen, y1p, y2p, w, green)

	z1p := VertexToGameScreen(z1, g, screen)
	z2p := VertexToGameScreen(z2, g, screen)
	StrokeLine(screen, z1p, z2p, w, blue)
}

