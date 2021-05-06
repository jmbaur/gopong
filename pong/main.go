package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

const (
	ballSpeedRatio      = 0.3
	ballSpeedMultiplier = 1.05
	ballDiameterRatio   = 0.03
	paddleSpeedRatio    = 1
)

var (
	frameRate  = 60 // Hz
	ballTheta0 float64
	side       float64
)

type pongBall struct {
	x, y, speed, angle, width, height float64
	element                           js.Value
}

func (b *pongBall) update() {
	b.element.Get("style").Call("setProperty", "left", b.x-(b.width/2))
	b.element.Get("style").Call("setProperty", "top", b.y-(b.height/2))
	b.element.Get("style").Call("setProperty", "width", b.width)
	b.element.Get("style").Call("setProperty", "height", b.height)
}

type pongPaddle struct {
	x, y, speed, angle, width, height float64
	element                           js.Value
}

func (b *pongPaddle) update() {
	b.element.Get("style").Call("setProperty", "left", b.x-(b.width/2))
	b.element.Get("style").Call("setProperty", "top", b.y-(b.height/2))
	b.element.Get("style").Call("setProperty", "width", b.width)
	b.element.Get("style").Call("setProperty", "height", b.height)
}

func (b *pongBall) getNextPosition() (x, y float64) {
	deltaX := b.speed * math.Cos(b.angle)
	deltaY := b.speed * math.Sin(b.angle)
	return b.x + deltaX, b.y + deltaY
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

func changeFPS(fps int, fpsDisplay js.Value, ball *pongBall, paddle *pongPaddle) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		frameRate = fps
		paddle.speed = side * paddleSpeedRatio / float64(frameRate)
		ball.speed = side * ballSpeedRatio / float64(frameRate)
		fpsDisplay.Set("innerHTML", fmt.Sprintf("%d Hz", frameRate))
		return nil
	})
}

func paddleCallback(p *pongPaddle, side float64) js.Func {
	return js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
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
				paddleEdge = p.y + p.speed + (p.height / 2)
				delta = p.speed
			case 33:
				fallthrough
			case 38:
				fallthrough
			case 75:
				// up
				paddleEdge = p.y - p.speed - (p.height / 2)
				delta = p.speed * -1
			}
			if 0 < paddleEdge && paddleEdge < side {
				p.y += delta
				// p.update()
			}
			return nil
		})
}

func main() {
	c := make(chan bool)

	rand.Seed(time.Now().UnixNano())
	// Get a random initial angle for the ball where it will be moving away
	// from the paddle within a 90 degree window.
	ballTheta0 = ((rand.Float64() * 90) - 45) * math.Pi / 180

	document := js.Global().Get("document")
	window := js.Global().Get("window")

	width := window.Get("innerWidth").Float()
	height := window.Get("innerHeight").Float()

	if width > height {
		side = height
	} else {
		side = width
	}
	side = side * 0.95

	root := document.Call("getElementById", "root")
	root.Get("style").Call("setProperty", "width", side)
	root.Get("style").Call("setProperty", "height", side)

	dashboard := document.Call("getElementById", "dashboard")

	message := document.Call("getElementById", "message")
	fps := document.Call("getElementById", "fps")
	fps.Set("innerHTML", fmt.Sprintf("%d Hz", frameRate))
	score := document.Call("getElementById", "score")
	scoreCount := 0
	score.Set("innerHTML", scoreCount)

	ball := &pongBall{
		x:       side / 2,
		y:       side / 2,
		width:   side * ballDiameterRatio,
		height:  side * ballDiameterRatio,
		speed:   side * ballSpeedRatio / float64(frameRate),
		angle:   ballTheta0,
		element: document.Call("getElementById", "ball"),
	}
	ballRadius := ball.width / 2
	ball.update()

	paddle := &pongPaddle{
		x:       side/75 + ((side / 25) / 2),
		y:       side / 2,
		speed:   side * paddleSpeedRatio / float64(frameRate),
		angle:   0.0,
		width:   side / 25,
		height:  side / 6,
		element: document.Call("getElementById", "paddle"),
	}
	paddle.update()

	fpsLo := document.Call("getElementById", "60")
	fpsLo.Call("addEventListener", "click", changeFPS(60, fps, ball, paddle))
	fpsMid := document.Call("getElementById", "144")
	fpsMid.Call("addEventListener", "click", changeFPS(144, fps, ball, paddle))
	fpsHi := document.Call("getElementById", "240")
	fpsHi.Call("addEventListener", "click", changeFPS(240, fps, ball, paddle))

	window.Call("addEventListener", "keydown", paddleCallback(paddle, side))

	for i := 3; i > 0; i-- {
		message.Set("innerHTML", i)
		time.Sleep(1 * time.Second)
	}

	message.Get("style").Call("setProperty", "visibility", "hidden")
	dashboard.Get("style").Call("setProperty", "visibility", "visible")

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
		paddle.update()

		time.Sleep(time.Duration(1000/frameRate) * time.Millisecond)
	}

	<-c
}
