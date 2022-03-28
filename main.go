package main

import (
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/inufuzei/puzzle2048/dnd"
	"github.com/inufuzei/puzzle2048/inu"
	"github.com/inufuzei/puzzle2048/tyoco"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 480
	screenHeight = 360
)

var (
	mPlus1pRegular_ttf font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	ft, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mPlus1pRegular_ttf = ft
}

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

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawSprites(screen)

	//text.Draw(screen, g.Msg, mPlus1pRegular_ttf, 120, 140, color.White)
	//for i, k := range g.keys {
	//posY := (1 + i) * 20
	//s := k.String()
	//text.Draw(screen, s, mPlus1pRegular_ttf, 100, posY, color.White)
	//}
	t := g.Questionlist[g.Questionnumber]
	q := t.Question
	a := t.Answer
	text.Draw(screen, q, mPlus1pRegular_ttf, 0, 24, color.White)
	if len(g.keys) > 0 {
		akey := g.keys[0]
		s := akey.String()
		if strings.HasPrefix(s, "Digit") {
			s = s[5:]
		}
		text.Draw(screen, s, mPlus1pRegular_ttf, 70, 24, color.White)

		if s == a {
			g.Questionnumber = g.Questionnumber + 1
			red := color.RGBA{
				R: 255,
				G: 100,
				B: 60,
				A: 255,
			}
			text.Draw(screen, "正解", mPlus1pRegular_ttf, 70, 75, red)
		} else {
			blue2 := color.RGBA{
				R: 60,
				G: 150,
				B: 255,
				A: 255,
			}
			text.Draw(screen, "残念", mPlus1pRegular_ttf, 70, 755, blue2)
		}
		text.Draw(screen, "FINAL ANSWER?????", mPlus1pRegular_ttf, 150, 230, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 480, 360
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Dogs run")
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
			dnd.NewBlock(100, 100, 50, 50, color.White),
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
