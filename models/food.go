package models

import (
	"math/rand"
	"snake/core"
)

type Food struct {
	X core.Point
}

const (
	SQUARE_SIZE = 10
)

func NewFoodAtRandom(xMax int, yMax int) Food {
	x := rand.Intn(xMax)
	y := rand.Intn(yMax)

	return Food{
		X: core.Point{
			X: x - x%SQUARE_SIZE,
			Y: y - y%SQUARE_SIZE,
		},
	}
}
