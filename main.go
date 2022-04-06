package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/inufuzei/puzzle2048/dnd"
)

const (
	screenWidth  = 400
	screenHeight = 400
	blockSize    = uint(4)
	shuffle      = 1000
)

type Game struct {
	touchIDs []ebiten.TouchID
	strokes  map[*dnd.Stroke]struct{}
	Blocks   []*dnd.Block
}

func (g *Game) Iscompleted() bool {
	for _, comp := range g.Blocks {
		regularX, regularY := comp.GetRegular()
		if regularX != comp.CellnumberX ||
			regularY != comp.CellnumberY {
			return false
		}
	}
	return true
}

func (g *Game) blockAt(x, y int) *dnd.Block {
	for i := len(g.Blocks) - 1; i >= 0; i-- {
		s := g.Blocks[i]
		if s.In(uint(x), uint(y)) {
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

	s := stroke.DraggingObject().(*dnd.Block)
	if s == nil {
		return
	}

	x, y := stroke.PositionDiff()
	s.MoveOn(x, y, g.Blocks)

	index := -1
	for i, ss := range g.Blocks {
		if ss == s {
			index = i
			break
		}
	}

	g.Blocks = append(g.Blocks[:index], g.Blocks[index+1:]...)
	g.Blocks = append(g.Blocks, s)
	stroke.SetDraggingObject(nil)
}

func (g *Game) updateStrokes() error {
	if g.Iscompleted() {
		return nil
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := dnd.NewStroke(&dnd.MouseStrokeSource{})
		s.SetDraggingObject(g.blockAt(s.Position()))
		g.strokes[s] = struct{}{}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := dnd.NewStroke(&dnd.TouchStrokeSource{ID: id})
		s.SetDraggingObject(g.blockAt(s.Position()))
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
	return nil
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
		if block := s.DraggingObject().(*dnd.Block); block != nil {
			block.Draw(screen)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBlocks(screen)
	if g.Iscompleted() {
		w := 400
		h := 400
		c := color.RGBA{
			R: 40,
			G: 60,
			B: 90,
			A: 200,
		}
		rect := ebiten.NewImage(w, h)
		rect.Fill(c)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(0), float64(0))
		screen.DrawImage(rect, op)

		text.Draw(screen, "完成", dnd.MPlus1pRegular_ttf, 170, 200, color.Black)
	}
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

	rn := rand.New(rand.NewSource(time.Now().UnixMicro()))
	for i := 0; i < shuffle; i++ {
		fh, fhh := dnd.FindHoll(blocks)
		fn := dnd.FindNeibs(blocks, fh, fhh)
		n := rn.Float64() * float64(len(fn))
		b := fn[int(n)]
		b.JustMove(blocks)
	}

	game := &Game{
		strokes: map[*dnd.Stroke]struct{}{},
		Blocks:  blocks,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
