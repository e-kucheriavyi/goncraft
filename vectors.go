package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Vector2 struct {
	X float64
	Y float64
}

func (v *Vector2) Scale(value float64) *Vector2 {
	v.X *= value
	v.Y *= value
	return v
}

func (v *Vector2) Clone() *Vector2 {
	return &Vector2{v.X, v.Y}
}

func (v *Vector2) Rotate(angle float64) *Vector2 {
	c := math.Cos(angle)
	s := math.Sin(angle)

	vr := &Vector2{
		X: v.X*c - v.Y*s,
		Y: v.Y*c + v.X*s,
	}

	v.X = vr.X
	v.Y = vr.Y

	return v
}

func (v *Vector2) Add(b *Vector2) *Vector2 {
	v.X += b.X
	v.Y += b.Y
	return v
}

func (v *Vector2) Mul(b *Vector2) *Vector2 {
	v.X *= b.X
	v.Y *= b.Y
	return v
}

func (v *Vector2) Sub(b *Vector2) *Vector2 {
	v.X -= b.X
	v.Y -= b.Y
	return v
}

func (v *Vector2) Div(b *Vector2) *Vector2 {
	v.X /= b.X
	v.Y /= b.Y
	return v
}

func (v *Vector2) ToScreen(screen *ebiten.Image) *Vector2 {
	w := float64(screen.Bounds().Dx())
	h := float64(screen.Bounds().Dy())

	return &Vector2{
		X: ((v.X + 1) / 2) * w,
		Y: ((v.Y + 1) / 2) * h,
	}
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func (v *Vector3) Clone() *Vector3 {
	return &Vector3{v.X, v.Y, v.Z}
}

func (v *Vector3) Scale(value float64) *Vector3 {
	v.X *= value
	v.Y *= value
	v.Z *= value
	return v
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

	vr := &Vector3{
		X: x,
		Y: y*yc - z*ys,
		Z: y*ys + z*yc,
	}

	v.X = vr.X
	v.Y = vr.Y
	v.Z = vr.Z

	return v
}

func (v *Vector3) Add(b *Vector3) *Vector3 {
	v.X += b.X
	v.Y += b.Y
	v.Z += b.Z

	return v
}

func (v *Vector3) Mul(b *Vector3) *Vector3 {
	v.X *= b.X
	v.Y *= b.Y
	v.Z *= b.Z

	return v
}

func (v *Vector3) Sub(b *Vector3) *Vector3 {
	v.X -= b.X
	v.Y -= b.Y
	v.Z -= b.Z

	return v
}

func (v *Vector3) Div(b *Vector3) *Vector3 {
	v.X /= b.X
	v.Y /= b.Y
	v.Z /= b.Z

	return v
}

func (v *Vector3) GetDistance(b *Vector3) float64 {
	return math.Sqrt(v.GetSqrDistance(b))
}

func (v *Vector3) GetSqrDistance(b *Vector3) float64 {
	x := b.X - v.X
	y := b.Y - v.Y
	z := b.Z - v.Z

	return x*x + y*y + z*z
}

func Vector3ToGameScreen(v *Vector3, g *Game, screen *ebiten.Image) *Vector2 {
	return v.Add(g.PlayerPosition).Rotate(g.PlayerRotation).Project().ToScreen(screen)
}
