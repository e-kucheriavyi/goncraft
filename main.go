package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	W = 800
	H = 800
	CELESTIAL_SIZE = 64
)

var (
	green = color.RGBA{0, 150, 0, 255}
	brown = color.RGBA{150, 150, 0, 255}
	SunColor = color.RGBA{255, 200, 0, 255}
	MoonColor = color.RGBA{0xDD, 0xE8, 0xEA, 0xFF}
)

type Point struct {
	X float64
	Y float64
}

func (p *Point) ToScreen(screen *ebiten.Image) *Point {
	w := float64(screen.Bounds().Dx())
	h := float64(screen.Bounds().Dy())

	return &Point{
		X: ((p.X + 1) / 2) * w,
		Y: (1 - ((p.Y + 1) / 2)) * h,
	}
}

type Vertex struct {
	X float64
	Y float64
	Z float64
}

func (v *Vertex) Project() *Point {
	return &Point{
		X: v.X / v.Z,
		Y: v.Y / v.Z,
	}
}

func (v *Vertex) Rotate(p *Point) *Vertex {
	xc := math.Cos(p.X)
	xs := math.Sin(p.X)

	yc := math.Cos(p.Y)
	ys := math.Sin(p.Y)

	vr := &Vertex{
		X: v.X*xc - v.Z*xs,
		Y: v.Y,
		Z: v.X*xs + v.Z*xc,
	}

	vr = &Vertex{
		X: vr.X,
		Y: vr.Y*yc - vr.Z*ys,
		Z: vr.Y*ys + vr.Z*yc,
	}

	return vr
}

func (v *Vertex) Translate(t *Vertex) *Vertex {
	return &Vertex{
		X: v.X + t.X,
		Y: v.Y + t.Y,
		Z: v.Z + t.Z,
	}
}

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
}

var vs = []*Vertex{
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

	bv := &Vertex{
		float64(b.X),
		float64(b.Y),
		float64(b.Z),
	}

	// for i, v := range vs {
	//	p := v.Translate(bv).Translate(g.PlayerPosition).Rotate(g.PlayerRotation).Project()
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

		sd := &Vertex{}

		hasLessThanZero := false

		var ysd float64 = 0

		for i, j := range r {
			v1 := vs[j]
			v1 = v1.Translate(bv)

			sd.X += v1.X
			sd.Y += v1.Y
			sd.Z += v1.Z

			ysd += v1.Y

			v1 = v1.Translate(g.PlayerPosition).Rotate(g.PlayerRotation)

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

			var v2 *Vertex

			if i == len(r)-1 {
				v2 = vs[r[0]]
			} else {
				v2 = vs[r[i+1]]
			}

			v2 = v2.Translate(bv).Translate(g.PlayerPosition).Rotate(g.PlayerRotation)

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

		cv := uint8(ysd*255)

		col := color.RGBA{cv, cv, cv, 255}

		c := ebiten.ColorScale{}
		c.ScaleWithColor(col)
		fillOpts := &vector.FillOptions{}
		pathOpts := &vector.DrawPathOptions{ColorScale: c}
		vector.FillPath(screen, pth, fillOpts, pathOpts)
	}
}

// world
// chunk

type Game struct {
	W              float64
	H              float64
	SunPosition    *Vertex
	PlayerPosition *Vertex
	PlayerRotation *Point
	Cursor         *Point
}

func NewGame() *Game {
	return &Game{
		SunPosition:    &Vertex{0, 10, 0},
		PlayerPosition: &Vertex{1, 1, 1},
		PlayerRotation: &Point{0, 0},
		Cursor:         &Point{0, 0},
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

func (g *Game) DrawSun(screen *ebiten.Image) {
	v1 := g.SunPosition.Translate(g.PlayerPosition)

	v1 = v1.Rotate(g.PlayerRotation)

	if v1.Z > 0 {
		return
	}

	sp := v1.Project().ToScreen(screen)

	vector.FillRect(screen, float32(sp.X), float32(sp.Y), CELESTIAL_SIZE, CELESTIAL_SIZE, SunColor, false)
}

func (g *Game) DrawMoon(screen *ebiten.Image) {
	v1 := &Vertex{
		g.SunPosition.X,
		g.SunPosition.Y * -1,
		g.SunPosition.Z,
	}

	v1 = v1.Translate(g.PlayerPosition)

	v1 = v1.Rotate(g.PlayerRotation)

	if v1.Z > 0 {
		return
	}

	sp := v1.Project().ToScreen(screen)

	vector.FillRect(screen, float32(sp.X), float32(sp.Y), CELESTIAL_SIZE, CELESTIAL_SIZE, MoonColor, false)
}

func (g *Game) Update() error {
	xt, yt := ebiten.CursorPosition()
	x := float64(xt)
	y := float64(yt)

	dt := float64(1) / float64(60)

	g.PlayerRotation.X += (((g.Cursor.X - x) * dt) * math.Pi)
	g.PlayerRotation.Y += (((g.Cursor.Y - y) * dt) * math.Pi)

	g.HandleMovement()
	g.HandleJump()

	g.Cursor.X = x
	g.Cursor.Y = y
	return nil
}

func (g *Game) HandleMovement() {
	dt := float64(1)/float64(60)

	var dv *Vertex

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dv = &Vertex{1*dt, 0, 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		dv = &Vertex{-1*dt, 0, 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		dv = &Vertex{0, 0, 1*dt}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		dv = &Vertex{0, 0, -1*dt}
	}

	if dv == nil {
		return
	}

	g.PlayerPosition = g.PlayerPosition.Translate(dv)
}

func (g *Game) HandleJump() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		dt := float64(1)/float64(60)
		g.PlayerPosition.Y += 1 * dt
		return
	}

	// keys := inpututil.AppendJustReleasedKeys(nil)
	//for _, key := range keys {
	//	if key == ebiten.KeySpace {
	//		g.PlayerPosition.Y += 1
	//	}
	//}
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
