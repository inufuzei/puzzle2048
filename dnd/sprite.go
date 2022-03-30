package dnd

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite represents an image.
type Sprite struct {
	image *ebiten.Image
	x     int
	y     int
}

func NewBlock(x, y int, w, h int, c color.Color) *Sprite {
	image := ebiten.NewImage(w, h)
	image.Fill(c)

	return &Sprite{
		image: image,
		x:     x,
		y:     y,
	}
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (s *Sprite) In(x, y int) bool {
	// Check the actual color (alpha) value at the specified position
	// so that the result of In becomes natural to users.
	//
	// Note that this is not a good manner to use At for logic
	// since color from At might include some errors on some machines.
	// As this is not so important logic, it's ok to use it so far.
	return s.image.At(x-s.x, y-s.y).(color.RGBA).A > 0
}

// MoveBy moves the sprite by (x, y).
func (s *Sprite) MoveBy(screenWidth, screenHeight int, x, y int) {
	w, h := s.image.Size()

	s.x += x
	s.y += y
	if s.x < 0 {
		s.x = 0
	}
	if s.x > screenWidth-w {
		s.x = screenWidth - w
	}
	if s.y < 0 {
		s.y = 0
	}
	if s.y > screenHeight-h {
		s.y = screenHeight - h
	}
}

// Draw draws the sprite.
func (s *Sprite) Draw(screen *ebiten.Image, dx, dy int, alpha float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x+dx), float64(s.y+dy))
	op.ColorM.Scale(1, 1, 1, alpha)
	screen.DrawImage(s.image, op)
	screen.DrawImage(s.image, op)
}
