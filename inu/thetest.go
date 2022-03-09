package inu

import (
	"fmt"
)

type Dog struct {
	color string
	speed float64
	power float64
}

func (x *Dog) hashiru() {
	fmt.Println("走っている", x.color)
}
func (x *Dog) stukareru() {
	fmt.Println("疲れている", x.color)
}

func main() {
	a := Dog{
		color: "白",
		speed: 16.23,
		power: 75.56,
	}
	b := Dog{
		color: "黒",
		speed: 32.18,
		power: 52.2,
	}
	c := Dog{
		color: "マダラ",
		speed: 36.24,
		power: 49.98,
	}
	a.hashiru()
	b.stukareru()
	c.hashiru()
}
