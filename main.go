package main

import (
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/inufuzei/puzzle2048/inu"
	"github.com/inufuzei/puzzle2048/tyoco"
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
	keys         []ebiten.Key
	Tester2187   inu.Dog
	Tester4893   inu.Dog
	Msg          string
	count        int
	Witch        bool
	Questionlist []tyoco.Tyoco
}

func (g *Game) Update() error {
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

func (g *Game) Draw(screen *ebiten.Image) {
	//text.Draw(screen, g.Msg, mPlus1pRegular_ttf, 120, 140, color.White)
	//for i, k := range g.keys {
	//posY := (1 + i) * 20
	//s := k.String()
	//text.Draw(screen, s, mPlus1pRegular_ttf, 100, posY, color.White)
	//}
	text.Draw(screen, "4+5=\n\n\n\nA,", mPlus1pRegular_ttf, 0, 24, color.White)
	if len(g.keys) > 0 {
		akey := g.keys[0]
		s := akey.String()
		if strings.HasPrefix(s, "Digit") {
			s = s[5:]
		}
		text.Draw(screen, s, mPlus1pRegular_ttf, 70, 24, color.White)
		if akey == ebiten.Key9 {
			red := color.RGBA{
				R: 255,
				G: 100,
				B: 60,
				A: 255,
			}
			text.Draw(screen, "正解", mPlus1pRegular_ttf, 70, 55, red)
		} else {
			blue2 := color.RGBA{
				R: 60,
				G: 150,
				B: 255,
				A: 255,
			}
			text.Draw(screen, "残念", mPlus1pRegular_ttf, 70, 55, blue2)
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
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
