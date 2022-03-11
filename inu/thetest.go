package inu

import (
	"fmt"
)

type Dog struct {
	Color string
	Speed float64
	Power float64
}

func (x *Dog) Hashiru() string {
	return fmt.Sprintf("走っている %v", x.Color)
}
func (x *Dog) Stukareru() string {
	return fmt.Sprintf("疲れている %v", x.Color)
}

func main() {
	a := Dog{
		Color: "白",
		Speed: 16.23,
		Power: 75.56,
	}
	b := Dog{
		Color: "黒",
		Speed: 32.18,
		Power: 52.2,
	}
	c := Dog{
		Color: "マダラ",
		Speed: 36.24,
		Power: 49.98,
	}
	a.Hashiru()
	b.Stukareru()
	c.Hashiru()
}
