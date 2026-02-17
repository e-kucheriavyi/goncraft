package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GIZMO_W = 4.0
)

var (
	red  = color.RGBA{255, 0, 0, 255}
	blue = color.RGBA{0, 0, 255, 255}
)

var (
	x1 = &Vector3{-1000, 0, 0}
	x2 = &Vector3{1000, 0, 0}
	y1 = &Vector3{0, 1000, 0}
	y2 = &Vector3{0, -1000, 0}
	z1 = &Vector3{0, 0, -1000}
	z2 = &Vector3{0, 0, 1000}
)

func (g *Game) DrawWorldAxis(screen *ebiten.Image) {
	x1p := Vector3ToGameScreen(x1, g, screen)
	x2p := Vector3ToGameScreen(x2, g, screen)
	StrokeLine(screen, x1p, x2p, GIZMO_W, red)

	y1p := Vector3ToGameScreen(y1, g, screen)
	y2p := Vector3ToGameScreen(y2, g, screen)
	StrokeLine(screen, y1p, y2p, GIZMO_W, green)

	z1p := Vector3ToGameScreen(z1, g, screen)
	z2p := Vector3ToGameScreen(z2, g, screen)
	StrokeLine(screen, z1p, z2p, GIZMO_W, blue)
}
