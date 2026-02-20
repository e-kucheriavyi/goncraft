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
	Pos  *Vector3
}

var blocks = []*Block{
	{Grass, &Vector3{0, 0, 0}},
	{Grass, &Vector3{0, 0, 1}},
	{Grass, &Vector3{0, 0, 2}},
	{Grass, &Vector3{0, 0, 3}},
	{Grass, &Vector3{1, 0, 0}},
	{Grass, &Vector3{2, 0, 0}},
	{Grass, &Vector3{3, 0, 0}},
	{Grass, &Vector3{0, 0, -1}},
	{Grass, &Vector3{0, 0, -2}},
	{Grass, &Vector3{0, 0, -3}},
	{Grass, &Vector3{-1, 0, 0}},
	{Grass, &Vector3{-2, 0, 0}},
	{Grass, &Vector3{-3, 0, 0}},
	{Dirt, &Vector3{-3, 1, 0}},
	{Dirt, &Vector3{-3, 2, 0}},
	{Dirt, &Vector3{-3, 3, 0}},
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
	{0.5, -0.5, 0.0}, // 0 side
	{1.0, -0.5, 0.5}, // 1 side
	{0.5, -1.0, 0.5}, // 2 bottom
	{0.0, -0.5, 0.5}, // 3 side
	{0.5, 0.0, 0.5},  // 4 top
	{0.5, -0.5, 1.0}, // 5 side
}

func (b *Block) Draw(screen *ebiten.Image, g *Game) {
	bv := b.Pos.Clone()

	// 1: compute all faces
	// 2: sort
	// 3: draw (last 3)

	type Face struct {
		Vs    []*Vector3
		D     float64
		N     *Vector3
		Shade float64
	}

	faces := make([]Face, 0, 12)

	maxD := float64(0.0)

	for i, r := range fs {
		skip := false

		faceVs := make([]*Vector3, 0, 16)
		for _, j := range r {
			v := vs[j].Clone()
			v1 := v.Add(bv)
			v1 = v1.Sub(g.PlayerPosition).Rotate(g.PlayerRotation)

			if v1.Z < 0 {
				skip = true
				break
			}

			p := v1.Project().ToScreen(screen)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", j), int(p.X), int(p.Y))

			faceVs = append(faceVs, v1)
		}
		if skip {
			continue
		}

		n := centers[i].Clone().Add(bv)

		d := math.Abs(n.Sub(g.PlayerPosition).Rotate(g.PlayerRotation).GetSqrDistance(g.PlayerPosition))

		maxD = math.Max(d, maxD)

		shade := 0.8

		if i == 4 {
			shade = 1.0
		} else if i == 6 {
			shade = 0.6
		}

		faces = append(faces, Face{
			Vs:    faceVs,
			D:     d,
			N:     n,
			Shade: shade,
		})
	}

	sort.Slice(faces, func(i, j int) bool {
		return faces[i].D > faces[j].D
	})

	for i, fs := range faces {
		// if i >= 3 {
		// 	break
		// }
		pth := &vector.Path{}

		skip := false

		for j, v := range fs.Vs {
			if v.Z < 0 {
				skip = true
				break
			}

			p := v.Project().ToScreen(screen)

			if j == 0 {
				pth.MoveTo(float32(p.X), float32(p.Y))
			} else {
				pth.LineTo(float32(p.X), float32(p.Y))
			}
		}

		if skip {
			continue
		}

		pth.Close()

		cv := 0.0

		if b.Type == Grass {
			cv = 150.0
		} else if b.Type == Dirt {
			cv = 100.0
		}

		shade := fs.D / maxD
		shade = faces[i].Shade
		nc := uint8(shade * cv)
		cc := color.RGBA{nc, nc, nc, 255}

		c := ebiten.ColorScale{}
		c.ScaleWithColor(cc)
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

		av := a.Pos.Clone()
		bv := b.Pos.Clone()

		return av.GetSqrDistance(g.PlayerPosition) > bv.GetSqrDistance(g.PlayerPosition)
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
