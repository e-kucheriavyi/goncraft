package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vector2 struct {
	X float64
	Y float64
}

func (v *Vector2) Rotate(angle float64) *Vector2 {
	c := math.Cos(angle)
	s := math.Sin(angle)

	return &Vector2{
		X: v.X*c - v.Y*s,
		Y: v.Y*c + v.X*s,
	}
}

func (v *Vector2) Add(b *Vector2) *Vector2 {
	return &Vector2{
		X: v.X + b.X,
		Y: v.Y + b.Y,
	}
}

func (v *Vector2) Mul(b *Vector2) *Vector2 {
	return &Vector2{
		X: v.X * b.X,
		Y: v.Y * b.Y,
	}
}

func (v *Vector2) Sub(b *Vector2) *Vector2 {
	return &Vector2{
		X: v.X - b.X,
		Y: v.Y - b.Y,
	}
}

func (v *Vector2) Div(b *Vector2) *Vector2 {
	return &Vector2{
		X: v.X / b.X,
		Y: v.Y / b.Y,
	}
}

func (v *Vector2) ToScreen(screen *ebiten.Image) *Vector2 {
	w := float64(screen.Bounds().Dx())
	h := float64(screen.Bounds().Dy())

	return &Vector2{
		X: ((v.X + 1) / 2) * w,
		Y: (1 - ((v.Y + 1) / 2)) * h,
	}
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func (v *Vector3) Project() *Vector2 {
	return &Vector2{
		X: v.X / v.Z,
		Y: v.Y / v.Z,
	}
}

func (v *Vector3) Rotate(p *Vector2) *Vector3 {
	xc := math.Cos(p.X)
	xs := math.Sin(p.X)

	yc := math.Cos(p.Y)
	ys := math.Sin(p.Y)

	x := v.X*xc - v.Z*xs
	y := v.Y
	z := v.X*xs + v.Z*xc

	return &Vector3{
		X: x,
		Y: y*yc - z*ys,
		Z: y*ys + z*yc,
	}
}

func (v *Vector3) Add(b *Vector3) *Vector3 {
	return &Vector3{
		X: v.X + b.X,
		Y: v.Y + b.Y,
		Z: v.Z + b.Z,
	}
}

func (v *Vector3) Mul(b *Vector3) *Vector3 {
	return &Vector3{
		X: v.X * b.X,
		Y: v.Y * b.Y,
		Z: v.Z * b.Z,
	}
}

func (v *Vector3) Sub(b *Vector3) *Vector3 {
	return &Vector3{
		X: v.X - b.X,
		Y: v.Y - b.Y,
		Z: v.Z - b.Z,
	}
}

func (v *Vector3) Div(b *Vector3) *Vector3 {
	return &Vector3{
		X: v.X / b.X,
		Y: v.Y / b.Y,
		Z: v.Z / b.Z,
	}
}

func Vector3ToGameScreen(v *Vector3, g *Game, screen *ebiten.Image) *Vector2 {
	return v.Add(g.PlayerPosition).Rotate(g.PlayerRotation).Project().ToScreen(screen)
}
