package dnd

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BlockWidth  = 100
	BlockHeight = 100
)

type Block struct {
	image       *ebiten.Image
	CellnumberX uint
	CellnumberY uint
	Number      uint
}

func MakeBlock(x, y uint, c color.Color) *Block {
	image := ebiten.NewImage(BlockWidth, BlockHeight)
	image.Fill(c)

	return &Block{
		image:       image,
		CellnumberX: x,
		CellnumberY: y,
	}
}

func (b *Block) GetRegular() (uint, uint) {
	rgularx := Block

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

func (b *Block) MoveOn(x, y int) {
	positionX := int(b.CellnumberX)
	b.CellnumberX = uint(positionX + x)

	positionY := int(b.CellnumberY)
	b.CellnumberY = uint(positionY + x)
}

func (b *Block) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	startX, _ := b.GetDotX()
	startY, _ := b.GetDotY()
	op.GeoM.Translate(float64(startX), float64(startY))
	screen.DrawImage(b.image, op)
}
