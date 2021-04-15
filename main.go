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
	fmt.Println(e.x, e.y)
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

	document := js.Global().Get("document")
	window := js.Global().Get("window")

	width := window.Get("innerWidth").Float()
	height := window.Get("innerHeight").Float()
	fmt.Println("width", width, "height", height)

	ball := &entity{
		x:       (width) / 2,
		y:       (height) / 2,
		width:   ballDiameter,
		height:  ballDiameter,
		speed:   10.0,
		angle:   ballTheta0,
		element: document.Call("getElementById", "ball"),
	}
	ball.element.Get("style").Call("setProperty", "position", "absolute")
	ball.element.Get("style").Call("setProperty", "background-color", "#00ff00")
	ball.element.Get("style").Call("setProperty", "width", ballDiameter)
	ball.element.Get("style").Call("setProperty", "height", ballDiameter)
	ball.element.Get("style").Call("setProperty", "border-radius", "50%")
	ball.update()

	paddle := &entity{
		x:       0,
		y:       0,
		speed:   2.0,
		angle:   0.0,
		element: document.Call("getElementById", "paddle"),
	}
	paddle.element.Get("style").Call("setProperty", "position", "absolute")
	paddle.element.Get("style").Call("setProperty", "background-color", "#ff0000")
	paddle.element.Get("style").Call("setProperty", "width", "2em")
	paddle.element.Get("style").Call("setProperty", "height", "10em")

	for {
		x, y := ball.getNextPosition()
		// TODO: check that ball hasn't collided
		if x-ballRadius <= 0 {
			fmt.Println("collide with left side of wall")
			os.Exit(0)
		}
		if x+ballRadius >= width {
			fmt.Println("collide with right side of wall")
			os.Exit(0)
		}
		if y-ballRadius <= 0 {
			fmt.Println("collide with top side of wall")
			os.Exit(0)
		}
		if y+ballRadius >= height {
			fmt.Println("collide with bottom side of wall")
			os.Exit(0)
		}
		ball.x = x
		ball.y = y
		ball.update()
		time.Sleep(refreshRate)
	}

	<-c
}
