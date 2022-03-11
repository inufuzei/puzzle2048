package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/inufuzei/puzzle2048/inu"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	Tester2187 inu.Dog
	Tester4893 inu.Dog
	Msg        string
	count      int
	Witch      bool
}

func (g *Game) Update() error {
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

func (g *Game) Draw(screen *ebiten.Image) {

	text.Draw(screen, g.Msg, mPlus1pRegular_ttf, 90, 120, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetFPSMode()
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
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
