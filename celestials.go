package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	CELESTIAL_SIZE = 64
)

var (
	SunColor  = color.RGBA{255, 200, 0, 255}
	MoonColor = color.RGBA{0xDD, 0xE8, 0xEA, 0xFF}
)

func (g *Game) DrawSun(screen *ebiten.Image) {
	v1 := g.SunPosition.Add(g.PlayerPosition)

	v1 = v1.Rotate(g.PlayerRotation)

	if v1.Z > 0 {
		return
	}

	sp := v1.Project().ToScreen(screen)

	vector.FillRect(screen, float32(sp.X), float32(sp.Y), CELESTIAL_SIZE, CELESTIAL_SIZE, SunColor, false)
}

func (g *Game) DrawMoon(screen *ebiten.Image) {
	v1 := &Vector3{
		g.SunPosition.X,
		g.SunPosition.Y * -1,
		g.SunPosition.Z,
	}

	v1 = v1.Add(g.PlayerPosition)

	v1 = v1.Rotate(g.PlayerRotation)

	if v1.Z > 0 {
		return
	}

	sp := v1.Project().ToScreen(screen)

	vector.FillRect(screen, float32(sp.X), float32(sp.Y), CELESTIAL_SIZE, CELESTIAL_SIZE, MoonColor, false)
}
