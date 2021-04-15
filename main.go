package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"syscall/js"
	"time"
)

const (
	ballDiameter = 30               // pixels
	ballRadius   = ballDiameter / 2 // pixels
	paddleWidth  = 30               // pixels
	paddleHeight = 150              // pixels
	paddleStep   = 10               // pixels
	paddleBorder = 20               // pixels
	refreshRate  = 1000 / 60        // Hz
)

var (
	ballTheta0 float64
)

type entity struct {
	x, y, speed, angle, width, height float64
	element                           js.Value
}

func (e *entity) update() error {
	e.element.Get("style").Call("setProperty", "left", e.x-(e.width/2))
	e.element.Get("style").Call("setProperty", "top", e.y-(e.height/2))
	e.element.Get("style").Call("setProperty", "width", e.width)
	e.element.Get("style").Call("setProperty", "height", e.height)
	return nil
}

func (e *entity) getNextPosition() (x, y float64) {
	deltaX := e.speed * math.Cos(e.angle)
	deltaY := e.speed * math.Sin(e.angle)
	return e.x + deltaX, e.y + deltaY
}

func main() {
	c := make(chan bool)

	rand.Seed(time.Now().UnixNano())
	// Get a random initial angle for the ball where it will be moving away
	// from the paddle within a 90 degree window.
	ballTheta0 = ((rand.Float64() * 90) - 45) * math.Pi / 180
	// ballTheta0 = (180) * math.Pi / 180

	document := js.Global().Get("document")
	window := js.Global().Get("window")

	totalWidth := window.Get("innerWidth").Float()
	totalHeight := window.Get("innerHeight").Float()
	width := totalWidth - totalWidth/10
	height := totalHeight - totalHeight/10

	root := document.Call("getElementById", "root")
	// root.Get("style").Call("setProperty", "margin", "auto")
	root.Get("style").Call("setProperty", "position", "relative")
	root.Get("style").Call("setProperty", "background-color", "#222222")
	root.Get("style").Call("setProperty", "width", width)
	root.Get("style").Call("setProperty", "height", height)

	ball := &entity{
		x:       width / 2,
		y:       height / 2,
		width:   ballDiameter,
		height:  ballDiameter,
		speed:   5.0,
		angle:   ballTheta0,
		element: document.Call("getElementById", "ball"),
	}
	ball.element.Get("style").Call("setProperty", "position", "absolute")
	ball.element.Get("style").Call("setProperty", "background-color", "#00ff00")
	ball.element.Get("style").Call("setProperty", "border-radius", "50%")
	ball.update()

	paddle := &entity{
		x:       paddleBorder * 2,
		y:       height / 2,
		speed:   10.0,
		angle:   0.0,
		width:   paddleWidth,
		height:  paddleHeight,
		element: document.Call("getElementById", "paddle"),
	}
	paddle.element.Get("style").Call("setProperty", "position", "absolute")
	paddle.element.Get("style").Call("setProperty", "background-color", "#ff0000")
	paddle.update()

	window.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		keyCode := e.Get("keyCode").Int()
		var paddleEdge float64
		var delta float64
		switch keyCode {
		case 34:
			fallthrough
		case 40:
			fallthrough
		case 74:
			// down
			paddleEdge = paddle.y + paddleStep + (paddle.height / 2)
			delta = paddleStep
		case 33:
			fallthrough
		case 38:
			fallthrough
		case 75:
			// up
			paddleEdge = paddle.y - paddleStep - (paddle.height / 2)
			delta = paddleStep * -1
		}
		if 0 < paddleEdge && paddleEdge < height {
			paddle.y += delta
			paddle.update()
		}
		return nil
	}))

	for {
		x, y := ball.getNextPosition()
		// collision with left wall
		if x-ballRadius <= 0 {
			fmt.Println("YOU LOSE")
			os.Exit(0)
		}
		if x+ballRadius >= width || x-ballRadius <= paddle.x+(paddle.width/2) {
			// collision with right wall
			if ball.angle > 0 {
				ball.angle = ball.angle + math.Pi - (2 * ball.angle)
			} else {
				ball.angle = ball.angle - math.Pi - (2 * ball.angle)
			}
		} else if y-ballRadius <= 0 || y+ballRadius >= height {
			// collision with top or bottom wall
			ball.angle = ball.angle * -1
		}
		ball.x = x
		ball.y = y
		ball.update()
		time.Sleep(refreshRate * time.Millisecond)
	}

	<-c
}
