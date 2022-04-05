package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inufuzei/puzzle2048/dnd"
	"github.com/inufuzei/puzzle2048/inu"
	"github.com/inufuzei/puzzle2048/tyoco"
)

const (
	screenWidth  = 400
	screenHeight = 400
	blockSize    = uint(4)
)

type Game struct {
	keys           []ebiten.Key
	Tester2187     inu.Dog
	Tester4893     inu.Dog
	Msg            string
	count          int
	Witch          bool
	Questionlist   []tyoco.Tyoco
	Questionnumber uint

	// For Drag&Drop
	touchIDs []ebiten.TouchID
	strokes  map[*dnd.Stroke]struct{}
	sprites  []*dnd.Sprite
	Blocks   []*dnd.Block
}

func (g *Game) spriteAt(x, y int) *dnd.Sprite {
	// As the sprites are ordered from back to front,
	// search the clicked/touched sprite in reverse order.
	for i := len(g.sprites) - 1; i >= 0; i-- {
		s := g.sprites[i]
		if s.In(x, y) {
			return s
		}
	}
	return nil
}

func (g *Game) updateEachStroke(stroke *dnd.Stroke) {
	stroke.Update()
	if !stroke.IsReleased() {
		return
	}

	s := stroke.DraggingObject().(*dnd.Sprite)
	if s == nil {
		return
	}

	x, y := stroke.PositionDiff()
	s.MoveBy(screenWidth, screenHeight, x, y)

	index := -1
	for i, ss := range g.sprites {
		if ss == s {
			index = i
			break
		}
	}

	// Move the dragged sprite to the front.
	g.sprites = append(g.sprites[:index], g.sprites[index+1:]...)
	g.sprites = append(g.sprites, s)

	stroke.SetDraggingObject(nil)
}

func (g *Game) updateStrokes() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := dnd.NewStroke(&dnd.MouseStrokeSource{})
		s.SetDraggingObject(g.spriteAt(s.Position()))
		g.strokes[s] = struct{}{}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := dnd.NewStroke(&dnd.TouchStrokeSource{ID: id})
		s.SetDraggingObject(g.spriteAt(s.Position()))
		g.strokes[s] = struct{}{}
	}

	for s := range g.strokes {
		g.updateEachStroke(s)
		if s.IsReleased() {
			delete(g.strokes, s)
		}
	}
	return nil
}

func (g *Game) Update() error {
	if err := g.updateStrokes(); err != nil {
		return err
	}
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.count = g.count + 1
	if g.count < 60 {
		return nil
	}
	g.count = 0
	if g.Witch {
		g.Witch = false
		g.Msg = g.Tester2187.Hashiru()
	} else {
		g.Witch = true
		g.Msg = g.Tester4893.Hashiru()
	}
	return nil
}

func (g *Game) drawSprites(screen *ebiten.Image) {
	draggingSprites := map[*dnd.Sprite]struct{}{}
	for s := range g.strokes {
		if sprite := s.DraggingObject().(*dnd.Sprite); sprite != nil {
			draggingSprites[sprite] = struct{}{}
		}
	}

	for _, s := range g.sprites {
		if _, ok := draggingSprites[s]; ok {
			continue
		}
		s.Draw(screen, 0, 0, 1)
	}
	for s := range g.strokes {
		dx, dy := s.PositionDiff()
		if sprite := s.DraggingObject().(*dnd.Sprite); sprite != nil {
			sprite.Draw(screen, dx, dy, 0.5)
		}
	}
}

func (g *Game) drawBlocks(screen *ebiten.Image) {
	draggingBlocks := map[*dnd.Block]struct{}{}
	for s := range g.strokes {
		blocks, ok := s.DraggingObject().(*dnd.Block)
		if ok && blocks != nil {
			draggingBlocks[blocks] = struct{}{}
		}
	}

	for _, s := range g.Blocks {
		if _, ok := draggingBlocks[s]; ok {
			continue
		}
		s.Draw(screen)
	}
	for s := range g.strokes {
		dx, dy := s.PositionDiff()
		if block := s.DraggingObject().(*dnd.Sprite); block != nil {
			block.Draw(screen, dx, dy, 0.5)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	//g.drawSprites(screen)
	g.drawBlocks(screen)

	w := 400
	h := 400
	c := color.RGBA{
		R: 40,
		G: 60,
		B: 90,
		A: 100,
	}
	rect := ebiten.NewImage(w, h)
	rect.Fill(c)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(0))
	screen.DrawImage(rect, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenNoWidth, screenNoHeight int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("15 puzzle")

	var blocks []*dnd.Block
	for i := 0; i < int(blockSize*blockSize-1); i++ {
		block := dnd.MakeBlock(uint(i), color.White)
		posX, posY := block.GetRegular()
		block.CellnumberX = posX
		block.CellnumberY = posY
		blocks = append(blocks, block)
	}

	game := &Game{
		Tester2187: inu.Dog{
			Color: "白",
			Speed: 16.23,
			Power: 75.56,
		},
		Tester4893: inu.Dog{
			Color: "黒",
			Speed: 32.18,
			Power: 52.2,
		},
		Questionlist: tyoco.Xlist,

		strokes: map[*dnd.Stroke]struct{}{},
		sprites: []*dnd.Sprite{
			dnd.Primitivestripe(0, 0, 100, 100, color.White),
		},
		Blocks: blocks,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
