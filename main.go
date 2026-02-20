package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"sort"

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
	{Dirt, -3, 1, 0},
	{Dirt, -3, 2, 0},
	{Dirt, -3, 3, 0},
}

var vs = []*Vector3{
	{0, 0, 0},
	{1, 0, 0},
	{1, -1, 0},
	{0, -1, 0},

	{0, 0, 1},
	{1, 0, 1},
	{1, -1, 1},
	{0, -1, 1},
}

var fs = [][]int{
	{0, 1, 2, 3},
	{1, 2, 6, 5},
	{2, 3, 7, 6},
	{3, 0, 4, 7},
	{4, 0, 1, 5},
	{6, 7, 4, 5},
}

var centers = []*Vector3{
	{0.5, -0.5, 0.0}, // 0
	{1.0, -0.5, 0.5}, // 1
	{0.5, -1.0, 0.5}, // 2
	{0.0, -0.5, 0.5}, // 3
	{0.5, 0.0, 0.5},  // 4
	{0.5, -0.5, 1.0}, // 5
}

func (b *Block) Draw(screen *ebiten.Image, g *Game) {
	bv := &Vector3{
		float64(b.X),
		float64(b.Y),
		float64(b.Z),
	}

	// 1: compute all faces
	// 2: sort
	// 3: draw (last 3)

	type Face struct {
		Vs []*Vector3
		D  float64
		N  *Vector3
	}

	faces := make([]Face, 0, 12)

	maxD := float64(0.0)

	for i, r := range fs {
		skip := false

		faceVs := make([]*Vector3, 0, 16)
		for _, j := range r {
			v := vs[j]
			v1 := v.Add(bv)
			v1 = v1.Sub(g.PlayerPosition).Rotate(g.PlayerRotation)

			p := v1.Project().ToScreen(screen)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", j), int(p.X), int(p.Y))

			if v1.Z < 0 {
				skip = true
				break
			}

			faceVs = append(faceVs, v1)
		}
		if skip {
			continue
		}

		n := centers[i].Add(bv)

		d := math.Abs(n.Sub(g.PlayerPosition).Rotate(g.PlayerRotation).GetSqrDistance(g.PlayerPosition))

		maxD = math.Max(d, maxD)

		n = n.Sub(g.PlayerPosition).Rotate(g.PlayerRotation)
		faces = append(faces, Face{
			Vs: faceVs,
			D:  d,
			N:  n,
		})
	}

	sort.Slice(faces, func(i, j int) bool {
		return faces[i].D > faces[j].D
	})

	for _, fs := range faces {
		// if i >= 3 {
		// 	break
		// }
		pth := &vector.Path{}

		for j, v := range fs.Vs {
			p := v.Project().ToScreen(screen)

			if j == 0 {
				pth.MoveTo(float32(p.X), float32(p.Y))
			} else {
				pth.LineTo(float32(p.X), float32(p.Y))
			}
		}

		pth.Close()

		cv := uint8(0)

		if b.Type == Grass {
			cv = 150
		} else if b.Type == Dirt {
			cv = 100
		}

		// col := color.RGBA{cv, cv, cv, 255}

		//nc := uint8((float32(i)/float32(10))*255)
		nc := uint8((fs.D / maxD) * float64(cv))
		cc := color.RGBA{nc, nc, nc, 255}

		c := ebiten.ColorScale{}
		c.ScaleWithColor(cc)
		fillOpts := &vector.FillOptions{}
		strokeOpts := &vector.StrokeOptions{Width: 2}
		pathOpts := &vector.DrawPathOptions{ColorScale: c}
		if 1 == 1 {
			vector.FillPath(screen, pth, fillOpts, pathOpts)
		}

		c.ScaleWithColor(green)
		pathOpts = &vector.DrawPathOptions{ColorScale: c}
		vector.StrokePath(screen, pth, strokeOpts, pathOpts)

		pn := fs.N.Project().ToScreen(screen)

		s := float32(8)

		cc = color.RGBA{nc, nc, 255, 255}

		vector.FillRect(screen, float32(pn.X)-s, float32(pn.Y)-s, s, s, cc, false)
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
		PlayerPosition: &Vector3{5, 5, 5},
		PlayerRotation: &Vector2{0, 0},
		Cursor:         &Vector2{0, 0},
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawMoon(screen)
	g.DrawSun(screen)

	sort.Slice(blocks, func(i, j int) bool {
		a := blocks[i]
		b := blocks[j]

		av := &Vector3{
			float64(a.X),
			float64(a.Y),
			float64(a.Z),
		}

		bv := &Vector3{
			float64(b.X),
			float64(b.Y),
			float64(b.Z),
		}

		return av.GetSqrDistance(g.PlayerPosition) < bv.GetSqrDistance(g.PlayerPosition)
	})

	for _, b := range blocks {
		b.Draw(screen, g)
	}

	// g.DrawWorldAxis(screen)

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

	vector.StrokeRect(screen, 0, 0, W, H, 2, green, false)
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
