package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

const (
	ballSpeedRatio      = 0.005
	ballSpeedMultiplier = 1.05
	ballDiameterRatio   = 0.03
	paddleSpeedRatio    = 0.02
	refreshRate         = 1000 / 120 * time.Millisecond
)

var (
	ballTheta0 float64
)

type entity struct {
	x, y, speed, angle, width, height float64
	element                           js.Value
}

func (e *entity) update() {
	e.element.Get("style").Call("setProperty", "left", e.x-(e.width/2))
	e.element.Get("style").Call("setProperty", "top", e.y-(e.height/2))
	e.element.Get("style").Call("setProperty", "width", e.width)
	e.element.Get("style").Call("setProperty", "height", e.height)
}

func (e *entity) getNextPosition() (x, y float64) {
	deltaX := e.speed * math.Cos(e.angle)
	deltaY := e.speed * math.Sin(e.angle)
	return e.x + deltaX, e.y + deltaY
}

func closestDivisibleNumber(n int, m int) int {
	q := n / m
	n1 := m * q

	var n2 int
	if n*m > 0 {
		n2 = m * (q + 1)
	} else {
		n2 = m * (q - 1)
	}

	if math.Abs(float64(n-n1)) < math.Abs(float64(n-n2)) {
		return n1
	}
	return n2
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

	width := window.Get("innerWidth").Float()
	height := window.Get("innerHeight").Float()
	var side float64
	if width > height {
		side = height
	} else {
		side = width
	}
	side = side * 0.95

	root := document.Call("getElementById", "root")
	root.Get("style").Call("setProperty", "width", side)
	root.Get("style").Call("setProperty", "height", side)

	message := document.Call("getElementById", "message")
	score := document.Call("getElementById", "score")
	scoreCount := 0
	score.Set("innerHTML", scoreCount)

	ball := &entity{
		x:       side / 2,
		y:       side / 2,
		width:   side * ballDiameterRatio,
		height:  side * ballDiameterRatio,
		speed:   side * ballSpeedRatio,
		angle:   ballTheta0,
		element: document.Call("getElementById", "ball"),
	}
	fmt.Println(ball.width)
	ballRadius := ball.width / 2
	ball.update()

	paddle := &entity{
		x:       side/75 + ((side / 25) / 2),
		y:       side / 2,
		speed:   side * paddleSpeedRatio,
		angle:   0.0,
		width:   side / 25,
		height:  side / 6,
		element: document.Call("getElementById", "paddle"),
	}
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
			paddleEdge = paddle.y + paddle.speed + (paddle.height / 2)
			delta = paddle.speed
		case 33:
			fallthrough
		case 38:
			fallthrough
		case 75:
			// up
			paddleEdge = paddle.y - paddle.speed - (paddle.height / 2)
			delta = paddle.speed * -1
		}
		if 0 < paddleEdge && paddleEdge < side {
			paddle.y += delta
			paddle.update()
		}
		return nil
	}))

	for i := 3; i > 0; i-- {
		message.Set("innerHTML", i)
		time.Sleep(1 * time.Second)
	}
	message.Get("style").Call("setProperty", "visibility", "hidden")
	score.Get("style").Call("setProperty", "visibility", "visible")
	for {
		x, y := ball.getNextPosition()
		if x-ballRadius <= 0 {
			// collision with left wall
			message.Set("innerHTML", "Game Over")
			message.Get("style").Call("setProperty", "visibility", "visible")
			break
		}
		if x-ballRadius <= paddle.x+(paddle.width/2) && (paddle.y-(paddle.height/2)) < y && y < (paddle.y+(paddle.height/2)) {
			// collision with paddle
			// Calculate an angle within a 90 degree window based on where the
			// ball hit the paddle.
			bounceScale := -1 * (paddle.y - y) / (paddle.height / 2)
			ball.angle = (bounceScale * 45) * (math.Pi / 180)
			scoreCount++
			score.Set("innerHTML", scoreCount)
			ball.speed = ball.speed * ballSpeedMultiplier
		} else if x+ballRadius >= side {
			// collision with right wall
			if ball.angle > 0 {
				ball.angle = ball.angle + math.Pi - (2 * ball.angle)
			} else {
				ball.angle = ball.angle - math.Pi - (2 * ball.angle)
			}
		} else if y-ballRadius <= 0 || y+ballRadius >= side {
			// collision with top or bottom wall
			ball.angle = ball.angle * -1
		}
		ball.x = x
		ball.y = y
		ball.update()
		time.Sleep(refreshRate)
	}

	<-c
}
