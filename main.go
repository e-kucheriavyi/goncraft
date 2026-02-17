package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	W = 800
	H = 800
)

var (
	green = color.RGBA{0, 150, 0, 255}
	brown = color.RGBA{150, 150, 0, 255}
)

type BlockType byte

const (
	Dirt BlockType = iota
	Grass
)

type Block struct {
	Type BlockType
	X    int
	Y    int
	Z    int
}

var blocks = []*Block{
	{Grass, 0, 0, 0},
	{Grass, 0, 0, 1},
	{Grass, 0, 0, 2},
	{Grass, 0, 0, 3},
	{Grass, 1, 0, 0},
	{Grass, 2, 0, 0},
	{Grass, 3, 0, 0},
	{Grass, 0, 0, -1},
	{Grass, 0, 0, -2},
	{Grass, 0, 0, -3},
	{Grass, -1, 0, 0},
	{Grass, -2, 0, 0},
	{Grass, -3, 0, 0},
	{Grass, -3, 1, 0},
	{Grass, -3, 2, 0},
	{Grass, -3, 3, 0},
}

var vs = []*Vector3{
	{0.5, 0.5, 0.5},
	{-0.5, 0.5, 0.5},
	{-0.5, -0.5, 0.5},
	{0.5, -0.5, 0.5},

	{0.5, 0.5, -0.5},
	{-0.5, 0.5, -0.5},
	{-0.5, -0.5, -0.5},
	{0.5, -0.5, -0.5},
}

var fs = [][]int{
	{0, 1, 2, 3},
	{4, 5, 6, 7},
	{0, 4, 7, 3},
	{1, 2, 6, 5},
	{0, 1, 5, 4},
	{2, 3, 7, 6},
}

func (b *Block) Draw(screen *ebiten.Image, g *Game) {
	// S := 1.0

	bv := &Vector3{
		float64(b.X),
		float64(b.Y),
		float64(b.Z),
	}

	// for i, v := range vs {
	//	p := v.Add(bv).Add(g.PlayerPosition).Rotate(g.PlayerRotation).Project()
	//	p = p.ToScreen(screen)
	//	vector.FillRect(
	//		screen,
	//		float32(p.X-S*0.5),
	//		float32(p.Y-S*0.5),
	//		float32(S),
	//		float32(S),
	//		brown,
	//		false,
	//	)
	//	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), int(p.X), int(p.Y))
	//}

	for _, r := range fs {
		pth := &vector.Path{}

		sd := &Vector3{}

		hasLessThanZero := false

		var ysd float64 = 0

		for i, j := range r {
			v1 := vs[j]
			v1 = v1.Add(bv)

			sd.X += v1.X
			sd.Y += v1.Y
			sd.Z += v1.Z

			ysd += v1.Y

			v1 = v1.Add(g.PlayerPosition).Rotate(g.PlayerRotation)

			p1 := v1.Project()
			p1 = p1.ToScreen(screen)

			if v1.Z < 0 {
				hasLessThanZero = true
				break
			}

			if i == 0 {
				pth.MoveTo(float32(p1.X), float32(p1.Y))
			} else {
				pth.LineTo(float32(p1.X), float32(p1.Y))
			}

			var v2 *Vector3

			if i == len(r)-1 {
				v2 = vs[r[0]]
			} else {
				v2 = vs[r[i+1]]
			}

			v2 = v2.Add(bv).Add(g.PlayerPosition).Rotate(g.PlayerRotation)

			if v2.Z < 0 {
				hasLessThanZero = true
				break
			}

			p2 := v2.Project()
			p2 = p2.ToScreen(screen)

			vector.StrokeLine(
				screen,
				float32(p1.X), float32(p1.Y), float32(p2.X), float32(p2.Y),
				1.0, green, false,
			)
		}

		if hasLessThanZero {
			continue
		}

		sd.X = sd.X / 4
		sd.Y = sd.Y / 4
		sd.Z = sd.Z / 4

		ysd /= 4

		//cv := uint8((sd.X + sd.Y + sd.Z) / 3 * 255)

		cv := uint8(ysd * 255)

		col := color.RGBA{cv, cv, cv, 255}

		c := ebiten.ColorScale{}
		c.ScaleWithColor(col)
		fillOpts := &vector.FillOptions{}
		pathOpts := &vector.DrawPathOptions{ColorScale: c}
		vector.FillPath(screen, pth, fillOpts, pathOpts)
	}
}

type Game struct {
	W              float64
	H              float64
	SunPosition    *Vector3
	PlayerPosition *Vector3
	PlayerRotation *Vector2
	Cursor         *Vector2
}

func NewGame() *Game {
	return &Game{
		SunPosition:    &Vector3{0, 10, 0},
		PlayerPosition: &Vector3{1, 1, 1},
		PlayerRotation: &Vector2{0, 0},
		Cursor:         &Vector2{0, 0},
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawMoon(screen)
	g.DrawSun(screen)

	for _, b := range blocks {
		b.Draw(screen, g)
	}

	g.DrawWorldAxis(screen)

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"POS: X: %.02f Y: %.02f Z: %.02f\nROT: X: %.02f Y: %.02f",
			g.PlayerPosition.X,
			g.PlayerPosition.Y,
			g.PlayerPosition.Z,
			g.PlayerRotation.X,
			g.PlayerRotation.Y,
		),
	)
}

func (g *Game) Update() error {
	g.HandleMouseMove()
	g.HandleMovement()
	g.HandleJump()

	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	g.W = float64(w)
	g.H = float64(h)
	return w, h
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(W, H)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
