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
	ballDiameter = 40 // pixels
	ballRadius   = ballDiameter / 2
	refreshRate  = (1000 / 60) * time.Millisecond
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

    width := window.Get("innerWidth").Float()
    height := window.Get("innerHeight").Float()

	ball := &entity{
		x:       (width) / 2,
		y:       (height) / 2,
		width:   ballDiameter,
		height:  ballDiameter,
		speed:   12.0,
		angle:   ballTheta0,
		element: document.Call("getElementById", "ball"),
	}
	ball.element.Get("style").Call("setProperty", "position", "absolute")
	ball.element.Get("style").Call("setProperty", "background-color", "#00ff00")
	ball.element.Get("style").Call("setProperty", "width", ballDiameter)
	ball.element.Get("style").Call("setProperty", "height", ballDiameter)
	ball.element.Get("style").Call("setProperty", "border-radius", "50%")
	ball.update()

	// paddle := &entity{
	// 	x:       0,
	// 	y:       0,
	// 	speed:   2.0,
	// 	angle:   0.0,
	// 	element: document.Call("getElementById", "paddle"),
	// }
	// paddle.element.Get("style").Call("setProperty", "position", "absolute")
	// paddle.element.Get("style").Call("setProperty", "background-color", "#ff0000")
	// paddle.element.Get("style").Call("setProperty", "width", "2em")
	// paddle.element.Get("style").Call("setProperty", "height", "10em")

	for {
		x, y := ball.getNextPosition()
        // collision with left wall
		if x-ballRadius <= 0 {
			fmt.Println("YOU LOSE")
            os.Exit(0)
		}
		if x+ballRadius >= width {
            // collision with right wall
            if ball.angle > 0 {
                ball.angle  = ball.angle + math.Pi - (2 * ball.angle)
            } else {
                ball.angle  = ball.angle - math.Pi - (2 * ball.angle)
            }
		} else if y-ballRadius <= 0 || y+ballRadius >= height {
            // collision with top or bottom wall
            ball.angle  = ball.angle * -1
		}
		ball.x = x
		ball.y = y
		ball.update()
		time.Sleep(refreshRate)
	}

	<-c
}
