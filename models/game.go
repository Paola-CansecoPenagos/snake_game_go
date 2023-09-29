package models

import (
	"math"
	"math/rand"
	"snake/core"
	"time"
)

type Game struct {
	Width        int
	Height       int
	Snake        Snake
	Food         Food
	Points       int
	Radius       int
	Walls        []Wall
	lastWallTime time.Time
}

func NewGame(width int, height int) Game {
	return Game{
		Width:        width,
		Height:       height,
		Snake:        NewSnake(width/2, height/2),
		Food:         NewFoodAtRandom(width, height),
		Points:       0,
		Radius:       SQUARE_SIZE,
		Walls:        make([]Wall, 0),
		lastWallTime: time.Now(),
	}
}

func (g *Game) checkBorderCollision() bool {
	snakePos := g.Snake.Head()
	return snakePos.X == -1 || snakePos.X == g.Width+1 || snakePos.Y == -1 || snakePos.Y == g.Height+1
}

func (g *Game) checkSnakeSelfCollision() bool {
	head := g.Snake.Head()
	count := 0

	for _, sp := range g.Snake.X {
		// TODO add radius to collision detection
		if sp.X == head.X && sp.Y == head.Y {
			count += 1
		}
	}

	return count == 2
}

func (g *Game) checkFoodCollision() bool {
	snakeHead := g.Snake.Head()

	radius := g.Radius / 2.0
	xDiff := float64(snakeHead.X + radius - (g.Food.X.X + radius))
	yDiff := float64(snakeHead.Y + radius - (g.Food.X.Y + radius))
	dist := math.Sqrt(xDiff*xDiff + yDiff*yDiff)

	return dist < float64(g.Radius)
}

func (g *Game) Update() bool {

	currentTime := time.Now()
	if currentTime.Sub(g.lastWallTime).Seconds() >= 10 {
		g.lastWallTime = currentTime
		go g.addWall()
	}

	if g.checkBorderCollision() || g.checkSnakeSelfCollision() {
		return false
	}

	go g.Snake.Move()

	if g.checkFoodCollision() {
		go g.Snake.Grow()
		g.Food = NewFoodAtRandom(g.Width, g.Height)
		g.Points += 1
	}

	if g.checkBorderCollision() || g.checkSnakeSelfCollision() || g.checkWallCollision() {
		return false
	}

	return true
}

func (g *Game) StartWallSpawner() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go g.addWall()
		}
	}
}

func (g *Game) addWall() {
	x := rand.Intn(g.Width)
	y := rand.Intn(g.Height)
	g.Walls = append(g.Walls, Wall{X: core.Point{X: x, Y: y}})
}

func (g *Game) checkWallCollision() bool {
	snakeHead := g.Snake.Head()

	// Comprueba si la cabeza de la serpiente colisiona con alguna pared
	for _, wall := range g.Walls {
		if snakeHead.X == wall.X.X && snakeHead.Y == wall.X.Y {
			return true
		}
	}

	return false
}
