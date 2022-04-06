package dnd

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	BlockWidth  = 100
	BlockHeight = 100
	blockSize   = uint(4)
)

var (
	MPlus1pRegular_ttf font.Face
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
	MPlus1pRegular_ttf = ft
}

type Block struct {
	image       *ebiten.Image
	CellnumberX uint
	CellnumberY uint
	Number      uint
}

func MakeBlock(number uint, c color.Color) *Block {
	image := ebiten.NewImage(BlockWidth, BlockHeight)
	image.Fill(c)

	return &Block{
		image:  image,
		Number: number,
	}
}

func (b *Block) GetRegular() (uint, uint) {
	regularx := b.Number % blockSize
	regulary := b.Number / blockSize
	return regularx, regulary
}

func (b *Block) GetDotX() (uint, uint) {
	startx := b.CellnumberX * BlockWidth
	finishx := ((b.CellnumberX + 1) * BlockWidth) - 1
	return startx, finishx
}

func (b *Block) GetDotY() (uint, uint) {
	starty := b.CellnumberY * BlockWidth
	finishy := ((b.CellnumberY + 1) * BlockWidth) - 1
	return starty, finishy
}

func (b *Block) In(x, y uint) bool {
	startX, finishX := b.GetDotX()
	startY, finishY := b.GetDotY()
	inX := startX <= x && finishX >= x
	inY := startY <= y && finishY >= y
	return inX && inY
}

func (b *Block) MoveOn(x, y int, touch []*Block) {
	log.Println("MoveOn:", x, y)
	posX := int(b.CellnumberX)
	posY := int(b.CellnumberY)
	if x > y {
		posX = posX + 1
	}

	if y > x {
		posY = posY + 1
	}

	if -y > x {
		posY = posY - 1
	}

	if -x > y {
		posX = posX - 1
	}

	if posX < 0 {
		posX = 0
	}

	if posX >= int(blockSize) {
		posX = int(blockSize) - 1
	}

	if posY < 0 {
		posY = 0
	}

	if posY >= int(blockSize) {
		posY = int(blockSize) - 1
	}

	log.Printf("Moved block(%v) at cell %v, %v\n",
		b.Number, b.CellnumberX, b.CellnumberY)

	for _, ather := range touch {
		if uint(posX) == ather.CellnumberX &&
			uint(posY) == ather.CellnumberY {
			return
		}
	}

	b.CellnumberX = uint(posX)
	b.CellnumberY = uint(posY)

}

func PickBlock(blocks []*Block, posX, posY uint) (*Block, bool) {
	for _, ather := range blocks {
		if uint(posX) == ather.CellnumberX &&
			uint(posY) == ather.CellnumberY {
			return ather, true
		}
	}
	return nil, false
}

func FindHoll(blocks []*Block) (uint, uint) {
	for x := 0; x < int(blockSize); x++ {
		for y := 0; y < int(blockSize); y++ {
			_, ishere := PickBlock(blocks, uint(x), uint(y))
			if !ishere {
				return uint(x), uint(y)
			}
		}
	}
	log.Fatal("NO HOLE!!!!!!!")
	return 0, 0
}

func FindNeibs(blocks []*Block, posX, posY uint) []*Block {
	var results []*Block
	if res, isHere := PickBlock(blocks, posX+1, posY); isHere {
		results = append(results, res)
	}
	if res, isHere := PickBlock(blocks, posX-1, posY); isHere {
		results = append(results, res)
	}
	if res, isHere := PickBlock(blocks, posX, posY+1); isHere {
		results = append(results, res)
	}
	if res, isHere := PickBlock(blocks, posX, posY-1); isHere {
		results = append(results, res)
	}
	return results
}
func (b *Block) JustMove(blocks []*Block) {
	posX := b.CellnumberX
	posY := b.CellnumberY
	if posX < blockSize-1 {
		if _, isHere := PickBlock(blocks, posX+1, posY); !isHere {
			b.CellnumberX = posX + 1
			b.CellnumberY = posY
			return
		}
	}
	if posX > 0 {
		if _, isHere := PickBlock(blocks, posX-1, posY); !isHere {
			b.CellnumberX = posX - 1
			b.CellnumberY = posY
			return
		}
	}
	if posY < blockSize-1 {
		if _, isHere := PickBlock(blocks, posX, posY+1); !isHere {
			b.CellnumberX = posX
			b.CellnumberY = posY + 1
			return
		}
	}
	if posY > 0 {
		if _, isHere := PickBlock(blocks, posX, posY-1); !isHere {
			b.CellnumberX = posX
			b.CellnumberY = posY - 1
			return
		}
	}
}

func (b *Block) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	startX, _ := b.GetDotX()
	startY, finishY := b.GetDotY()
	op.GeoM.Translate(float64(startX), float64(startY))
	screen.DrawImage(b.image, op)
	moji := fmt.Sprintf("%v", b.Number)
	text.Draw(screen, moji, MPlus1pRegular_ttf,
		int(startX+45), int(finishY-45), color.Black)

	//log.Printf("Drawing block(%v) at %v, %v\n", b.Number, startX, startY)
}
