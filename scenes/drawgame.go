package ui

import (
	"fmt"
	"snake/core"
	"snake/models"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameUi struct {
	settings       core.Settings
	game           models.Game
	keys           []ebiten.Key
	newKeys        []ebiten.Key
	paused         bool
	lastPauseEvent time.Time
	difficulty     int
	tick           int
	status         bool
}

func NewGame(settings core.Settings) GameUi {
	return GameUi{
		settings: settings,
		game: models.NewGame(
			settings.Width-core.SQUARE_SIZE*3,
			settings.Height-core.SQUARE_SIZE*3-core.TOP_BAR_HEIGHT,
		),
		keys:       []ebiten.Key{},
		newKeys:    []ebiten.Key{},
		paused:     false,
		difficulty: 3,
		tick:       0,
		status:     true,
	}
}

func (g *GameUi) pause() {
	g.paused = !g.paused
}

func (g *GameUi) restart() {
	g.game = models.NewGame(
		g.settings.Width-core.SQUARE_SIZE*3,
		g.settings.Height-core.SQUARE_SIZE*3-core.TOP_BAR_HEIGHT,
	)
	g.keys = g.keys[:0]
	g.newKeys = g.newKeys[:0]
	g.status = true
	g.tick = 0
}

func (g *GameUi) readLastKey() ebiten.Key {
	g.newKeys = inpututil.AppendPressedKeys(g.newKeys[:0])
	keyPressed := ebiten.Key(-1)

	if len(g.newKeys) == 1 {
		keyPressed = g.newKeys[0]
	} else if len(g.keys) == 1 && len(g.newKeys) >= 1 {
		lastKeyPressed := g.keys[0]

		for _, k := range g.newKeys {
			if k != lastKeyPressed {
				keyPressed = k
			}
		}
	}

	g.keys = g.keys[:0]
	g.keys = append(g.keys, g.newKeys...)

	return keyPressed
}

func (g *GameUi) handleSpaceKey() {
	if time.Since(g.lastPauseEvent).Milliseconds() > 250 {
		g.lastPauseEvent = time.Now()
		if !g.status {
			go g.restart()
		} else {
			go g.pause()
		}
	}
}

func (g *GameUi) Update() error {
	direction := core.Direction(core.DIRECTION_NONE)
	keyPressed := g.readLastKey()

	switch keyPressed {
	case ebiten.KeyArrowRight:
		direction = core.DIRECTION_RIGHT
	case ebiten.KeyD:
		direction = core.DIRECTION_RIGHT
	case ebiten.KeyArrowLeft:
		direction = core.DIRECTION_LEFT
	case ebiten.KeyA:
		direction = core.DIRECTION_LEFT
	case ebiten.KeyArrowUp:
		direction = core.DIRECTION_UP
	case ebiten.KeyW:
		direction = core.DIRECTION_UP
	case ebiten.KeyArrowDown:
		direction = core.DIRECTION_DOWN
	case ebiten.KeyS:
		direction = core.DIRECTION_DOWN
	case ebiten.KeySpace:
		g.handleSpaceKey()
	case ebiten.KeyEscape:
		g.handleSpaceKey()
	}

	if g.paused {
		return nil
	}

	if !g.status {
		return nil
	}

	g.game.Snake.SetDirection(direction)

	if g.tick == g.difficulty {
		g.tick = 0
		g.status = g.game.Update()
	}

	g.tick += 1

	return nil
}

func (g *GameUi) Draw(screen *ebiten.Image) {
	xOffset := core.SQUARE_SIZE
	yOffset := core.SQUARE_SIZE + core.TOP_BAR_HEIGHT

	g.drawBorder(screen)
	g.drawSnake(screen, xOffset, yOffset)
	g.drawFood(screen, xOffset, yOffset)

	if g.paused {
		ebitenutil.DebugPrintAt(screen, "Game paused", 120, 2)
	}

	if !g.status {
		ebitenutil.DebugPrintAt(screen, "Game over", (g.settings.Width/2)-25, g.settings.Height/2)
		ebitenutil.DebugPrintAt(screen, "Press space to restart", (g.settings.Width/2)-60, (g.settings.Height/2)+20)
	}
	g.drawWalls(screen, xOffset, yOffset)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Points %d", g.game.Points), 10, 2)
}

func (g *GameUi) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.settings.Width, g.settings.Height
}

func (g *GameUi) drawBorder(screen *ebiten.Image) {
	ebitenutil.DrawRect(
		screen,
		0,
		float64(core.TOP_BAR_HEIGHT),
		float64(g.settings.Width),
		float64(core.SQUARE_SIZE),
		g.settings.SnakeColor,
	)
	ebitenutil.DrawRect(
		screen,
		0,
		float64(g.settings.Height-core.SQUARE_SIZE),
		float64(g.settings.Width),
		float64(core.SQUARE_SIZE),
		g.settings.SnakeColor,
	)
	ebitenutil.DrawRect(
		screen,
		0,
		float64(core.TOP_BAR_HEIGHT),
		float64(core.SQUARE_SIZE),
		float64(g.settings.Height),
		g.settings.SnakeColor,
	)
	ebitenutil.DrawRect(
		screen,
		float64(g.settings.Width-core.SQUARE_SIZE),
		float64(core.TOP_BAR_HEIGHT),
		float64(core.SQUARE_SIZE),
		float64(g.settings.Height),
		g.settings.SnakeColor,
	)
}

func (g *GameUi) drawSnake(screen *ebiten.Image, xOffset int, yOffset int) {
	for _, sp := range g.game.Snake.X {
		ebitenutil.DrawRect(
			screen,
			float64(sp.X+xOffset),
			float64(sp.Y+yOffset),
			float64(core.SQUARE_SIZE),
			float64(core.SQUARE_SIZE),
			g.settings.SnakeColor,
		)
	}
}

func (g *GameUi) drawFood(screen *ebiten.Image, xOffset int, yOffset int) {
	ebitenutil.DrawRect(
		screen,
		float64(g.game.Food.X.X+xOffset),
		float64(g.game.Food.X.Y+yOffset),
		float64(core.SQUARE_SIZE),
		float64(core.SQUARE_SIZE),
		g.settings.FoodColor,
	)
}

func (g *GameUi) drawWalls(screen *ebiten.Image, xOffset int, yOffset int) {
	for _, wall := range g.game.Walls {
		ebitenutil.DrawRect(
			screen,
			float64(wall.X.X+xOffset),
			float64(wall.X.Y+yOffset),
			float64(core.SQUARE_SIZE),
			float64(core.SQUARE_SIZE),
			g.settings.SnakeColor,
		)
	}
}
