package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MAX_CAM_TILT = -math.Pi * 0.5
	MIN_CAM_TILT = -math.Pi * 1.5
)

func (g *Game) HandleMouseMove() {
	xt, yt := ebiten.CursorPosition()
	x := float64(xt)
	y := float64(yt)

	dt := float64(1) / float64(60)

	g.PlayerRotation.X += (((g.Cursor.X - x) * dt) * math.Pi)
	g.PlayerRotation.Y += (((g.Cursor.Y - y) * dt) * math.Pi)

	if g.PlayerRotation.Y > MAX_CAM_TILT {
		g.PlayerRotation.Y = MAX_CAM_TILT
	} else if g.PlayerRotation.Y < MIN_CAM_TILT {
		g.PlayerRotation.Y = MIN_CAM_TILT
	}

	g.Cursor.X = x
	g.Cursor.Y = y
}

func (g *Game) HandleMovement() {
	dt := float64(1) / float64(60)

	var dv *Vector3

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dv = &Vector3{1 * dt, 0, 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		dv = &Vector3{-1 * dt, 0, 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		dv = &Vector3{0, 0, 1 * dt}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		dv = &Vector3{0, 0, -1 * dt}
	}

	if dv == nil {
		return
	}

	g.PlayerPosition = g.PlayerPosition.Add(dv)
}

func (g *Game) HandleJump() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		dt := float64(1) / float64(60)
		g.PlayerPosition.Y += 1 * dt
		return
	}
}
