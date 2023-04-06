package main

import (
	"fmt"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"math/rand"
	"sync"
	"time"
)

const (
	GameW = 720.0
	GameH = 720.0
)

type Game struct {
	cv  *canvas.Canvas
	wnd *sdlcanvas.Window

	worldS   float64
	snake    *Snake
	gameOver bool
	speed    int
	needMove bool
	food     []Point
}

// NewGame - конструктор игры
func NewGame() *Game {
	wnd, cv, err := sdlcanvas.CreateWindow(1080, 750, "Hello, Snake!")
	if err != nil {
		panic(err)
	}
	g := &Game{
		cv:       cv,
		wnd:      wnd,
		speed:    300,
		gameOver: false,
	}
	return g
}

// SetSnake - инициализация змеи
func (g *Game) SetSnake(s *Snake) {
	g.snake = s
}

// CreateWorld - создание мира
func (g *Game) CreateWorld(s float64) {
	g.worldS = s
}

// Run - запуск игры
func (g *Game) Run() {
	go g.SnakeMovement()
	go g.foodGeneration()
	g.renderLoop()
}

func (g *Game) Exit() {
	defer g.wnd.Destroy()
}

func (g *Game) foodGeneration() {
	var foodTimer *time.Timer
	resetTimer := func() {
		foodTimer = time.NewTimer(3 * time.Second)
		_ = foodTimer
	}
	resetTimer()

	for {
		<-foodTimer.C
		if !g.gameOver {
			min := 1
			max := 20 - 1
			randX := rand.Intn(max-min) + min
			randY := rand.Intn(max-min) + min
			newPoint := Point{float64(randX), float64(randY)}

			check := true
			if g.snake.IsSnake(newPoint) {
				check = false
			}
			for _, p := range g.food {
				if p.X == newPoint.X && p.Y == newPoint.Y {
					check = false
					break
				}
			}

			if check {
				g.food = append(g.food, newPoint)
			}
		}
		resetTimer()
	}
}

func (g *Game) SnakeMovement() {
	var snakeTimer *time.Timer
	var snakeDir Dir = Right
	var snakeLock sync.Mutex

	// resetTime - обнуление таймера змейки
	resetTimer := func() {
		snakeTimer = time.NewTimer(time.Duration(g.speed) * time.Millisecond)
		_ = snakeTimer
	}
	resetTimer()

	// обработка нажатий клавиш
	g.wnd.KeyUp = func(code int, rn rune, name string) {

		if code < 79 && code > 82 || g.needMove {
			return
		}

		snakeLock.Lock()

		newDir := snakeDir
		switch code {
		case 80: //влево
			newDir = Left
		case 82: //вниз
			newDir = Bottom
		case 79: //направо
			newDir = Right
		case 81: //вверх
			newDir = Top
		}

		if !snakeDir.CheckParallel(newDir) {
			snakeDir = newDir
			g.needMove = true
		}

		snakeLock.Unlock()

	}
	// передвижение змеи (обновление)
	for {
		<-snakeTimer.C
		snakeLock.Lock()

		if !g.gameOver {
			newPos := snakeDir.Exec(g.snake.Head())
			if newPos.X <= 0 || newPos.X >= g.worldS-1 || newPos.Y <= 0 || newPos.Y >= g.worldS-1 {
				g.gameOver = true
			}

			g.snake.CutIfSnake(newPos)

			// змейка попала на еду
			isFood := false
			for i := range g.food {
				if newPos.X == g.food[i].X && newPos.Y == g.food[i].Y {
					g.food = append(g.food[:i], g.food[i+1:]...)
					g.snake.Add(newPos)
					g.speed -= 25
					isFood = true
					break
				}
			}

			if !isFood {
				g.snake.Move(snakeDir)
				g.needMove = false
			}

		}
		snakeLock.Unlock()
		resetTimer()
	}
}

// renderLoop - отвечает за рендеринг пространства
func (g *Game) renderLoop() {

	gameAreaSP := Point{15, 15}
	gameAreaEP := Point{15 + GameW, 15 + GameH}

	cellW := GameW / g.worldS
	cellH := GameH / g.worldS

	font, err := g.cv.LoadFont("./tahoma.ttf")
	if err != nil {
		panic(err)
	}
	g.wnd.MainLoop(func() {
		// очистка окна
		g.cv.ClearRect(0, 0, 1080, 750)
		// отрисовка мира
		g.cv.BeginPath()
		g.cv.SetFillStyle("#333")
		g.cv.FillRect(gameAreaSP.X, gameAreaSP.Y, gameAreaEP.X-15, gameAreaEP.Y-15)
		g.cv.Stroke()

		g.cv.BeginPath()
		g.cv.SetStrokeStyle("#FFF001")
		g.cv.SetLineWidth(1)
		for i := 0; i < int(g.worldS)+1; i++ {
			g.cv.MoveTo(gameAreaSP.X+float64(i)*cellH, gameAreaSP.Y)
			g.cv.LineTo(gameAreaSP.X+float64(i)*cellH, gameAreaEP.Y)
		}
		for i := 0; i < 20+1; i++ {
			g.cv.MoveTo(gameAreaSP.X, gameAreaSP.Y+float64(i)*cellW)
			g.cv.LineTo(gameAreaEP.X, gameAreaSP.Y+float64(i)*cellW)
		}
		g.cv.Stroke()

		// стенки
		g.cv.BeginPath()
		g.cv.SetFillStyle("#ccc")

		//верх
		for i := 0; i < int(g.worldS); i++ {
			g.cv.FillRect(
				gameAreaSP.X+float64(i)*cellW+1,
				gameAreaSP.Y,
				cellW-1*2,
				cellH)
		}

		//низ
		for i := 0; i < int(g.worldS); i++ {
			g.cv.FillRect(
				gameAreaSP.X+float64(i)*cellW+1,
				gameAreaSP.Y+cellH*(g.worldS-1),
				cellW-1*2,
				cellH)
		}

		//лево
		for i := 0; i < int(g.worldS)-1; i++ {
			g.cv.FillRect(
				gameAreaSP.X,
				gameAreaSP.Y+float64(i)*cellH+1,
				cellW,
				cellH-1*2)
		}

		//право
		for i := 0; i < int(g.worldS)-1; i++ {
			g.cv.FillRect(
				gameAreaSP.X+cellH*(g.worldS-1),
				gameAreaSP.Y+float64(i)*cellH+1,
				cellW,
				cellH-1*2)
		}
		g.cv.Stroke()
		// отрисовка змеи
		g.cv.BeginPath()
		g.cv.SetFillStyle("#FFF")
		for _, p := range g.snake.Parts {
			g.cv.FillRect(
				gameAreaSP.X+p.X*cellW+1,
				gameAreaSP.Y+p.Y*cellH+1,
				cellW-1*2,
				cellH-1*2,
			)
		}
		g.cv.Stroke()

		// отрисовка еды
		g.cv.BeginPath()
		g.cv.SetFillStyle("#F15555")
		for _, p := range g.food {
			g.cv.FillRect(
				gameAreaSP.X+p.X*cellW+1,
				gameAreaSP.Y+p.Y*cellH+1,
				cellW-1*2,
				cellH-1*2)
		}
		g.cv.Stroke()
		// отрисовка счёта
		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text := fmt.Sprintf("Score: %d", g.snake.Len())
		g.cv.FillText(text, 750+50, 50)

		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text = fmt.Sprintf("Food: %d", len(g.food))
		g.cv.FillText(text, 750+50, 85)

		g.cv.BeginPath()
		g.cv.SetFont(font, 25)
		text = fmt.Sprintf("Speed: %d", 350-g.speed)
		g.cv.FillText(text, 750+50, 120)

		if g.gameOver {
			g.cv.BeginPath()
			g.cv.SetFont(font, 30)
			text = fmt.Sprintf("Game Over :(")
			g.cv.FillText(text, 750+50, 175)

		}
	})
}
