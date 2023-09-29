package models

import "snake/core"

type Snake struct {
	X                []core.Point
	direction        core.Direction
	desiredDirection core.Direction
}

func NewSnake(x int, y int) Snake {
	return Snake{
		X:         []core.Point{{X: x, Y: y}},
		direction: core.DIRECTION_RIGHT,
	}
}

func (s *Snake) Grow() {
	for i := 0; i < 5; i++ {
		// TODO redesign snake data model to reduce allocations
		s.X = append([]core.Point{s.X[len(s.X)-1]}, s.X...)
	}
}

func (s *Snake) Head() core.Point {
	return s.X[len(s.X)-1]
}

func (s *Snake) Tail() core.Point {
	return s.X[0]
}

func (s *Snake) canChangeDirection(direction core.Direction) bool {
	diff := s.direction - direction
	if diff < 0 {
		diff = -diff
	}

	return diff != 1
}

func (s *Snake) SetDirection(direction core.Direction) {
	if direction != core.DIRECTION_NONE && s.canChangeDirection(direction) {
		s.desiredDirection = direction
	}
}

func (s *Snake) GetDirection() core.Direction {
	return s.direction
}

func (s *Snake) Move() {
	head := s.Head()
	tail := s.Tail()

	if ((s.desiredDirection == core.DIRECTION_UP || s.desiredDirection == core.DIRECTION_DOWN) && head.X%10 == 0) ||
		((s.desiredDirection == core.DIRECTION_LEFT || s.desiredDirection == core.DIRECTION_RIGHT) && head.Y%10 == 0) {
		s.direction = s.desiredDirection
	}

	switch s.direction {
	case core.DIRECTION_RIGHT:
		tail.X = head.X + 1
		tail.Y = head.Y
	case core.DIRECTION_LEFT:
		tail.X = head.X - 1
		tail.Y = head.Y
	case core.DIRECTION_UP:
		tail.X = head.X
		tail.Y = head.Y - 1
	case core.DIRECTION_DOWN:
		tail.X = head.X
		tail.Y = head.Y + 1
	}

	s.X = append(s.X[1:], tail)
}
