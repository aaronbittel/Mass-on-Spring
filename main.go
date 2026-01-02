package main

import (
	"fmt"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WIDTH  int32 = 900
	HEIGHT int32 = 600

	FONTSIZE int32 = 20

	RECT_SIZE    int32 = 50
	GROUND_THICK int32 = 6

	EQUILIBRIUM_X float32 = float32(WIDTH)/2 - float32(RECT_SIZE)/2

	SPRING_NUM      int32   = 12
	SPRING_LEN      float32 = 30
	SPRING_THICK    float32 = 2
	SPRING_Y_OFFSET float32 = 5

	START_VELOCITY float32 = 40

	SPRING_STIFFNESS    float32 = 50
	OBJ_MASS            float32 = 1
	DAMPING_COEFFICIENT float32 = 2

	EPSILON   float64 = 0.1
	THRESHOLD float64 = 10
)

type State string

const (
	Idle              State = "Idle"
	RunningSimulation       = "Simulation"
)

var (
	state       = Idle
	DraggingObj = false
)

func main() {
	rl.InitWindow(WIDTH, HEIGHT, "Mass on Spring - Simulation")
	defer rl.CloseWindow()

	var (
		titleText = "Mass on Spring - Simulation"
		titleSize = rl.MeasureText(titleText, FONTSIZE)
		instText  = "Drag the Object. Release to start the Simulation."
		instSize  = rl.MeasureText(instText, FONTSIZE)
	)

	rl.SetTargetFPS(60)

	var (
		x float32 = EQUILIBRIUM_X
		v float32 = START_VELOCITY
	)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawFPS(10, 10)
		rl.DrawText(titleText, WIDTH/2-titleSize/2, int32(float64(HEIGHT)*0.1), FONTSIZE, rl.DarkGreen)
		rl.DrawText(instText, WIDTH/2-instSize/2, int32(float64(HEIGHT)*0.15), FONTSIZE, rl.DarkGray)

		dt := rl.GetFrameTime()

		switch state {
		case Idle:
			x = DrawIdle(x)
		case RunningSimulation:
			x, v = RunSimulation(x, v, dt)
		}

		status := fmt.Sprintf("Current State: %s", state)
		rl.DrawText(status, 3, HEIGHT-FONTSIZE, FONTSIZE, rl.DarkGray)

		rl.EndDrawing()
	}
}

func DrawIdle(x float32) float32 {
	objRect := rl.Rectangle{
		X:      x,
		Y:      float32(HEIGHT)/2 - float32(RECT_SIZE) - float32(GROUND_THICK)/2,
		Width:  float32(RECT_SIZE),
		Height: float32(RECT_SIZE),
	}

	color := rl.Gray
	mousePos := rl.GetMousePosition()

	if rl.CheckCollisionPointRec(mousePos, objRect) {
		rl.SetMouseCursor(rl.MouseCursorPointingHand)
		if rl.CheckCollisionPointRec(mousePos, objRect) {
			DraggingObj = true
			color = rl.Maroon
		}
	} else {
		if !DraggingObj {
			rl.SetMouseCursor(rl.MouseCursorDefault)
		}
	}

	if !rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		DraggingObj = false
	}

	if DraggingObj {
		x = rl.Clamp(mousePos.X, float32(RECT_SIZE), float32(WIDTH)-float32(RECT_SIZE))
	}

	drawGround()
	drawMass(x, color)
	drawSpring(x)

	if !DraggingObj && math.Abs(float64(x)-float64(EQUILIBRIUM_X)) > THRESHOLD {
		state = RunningSimulation
	}

	return x
}

func RunSimulation(x, v, dt float32) (float32, float32) {
	a := acceleration(x, v)
	if math.Abs(float64(x-EQUILIBRIUM_X)) < EPSILON && math.Abs(float64(v)) < EPSILON {
		state = Idle
	}
	v += a * dt
	x += v * dt

	drawGround()
	drawMass(x)
	drawSpring(x)

	return x, v
}

func drawMass(x float32, color ...color.RGBA) {
	c := rl.Maroon
	if len(color) >= 1 {
		c = color[0]
	}
	rl.DrawRectangle(int32(x), HEIGHT/2-RECT_SIZE-GROUND_THICK/2, RECT_SIZE, RECT_SIZE, c)
}

func drawSpring(x float32) {
	sectionWidth := x / (float32(SPRING_NUM) / 2)

	var (
		start = rl.Vector2{X: 0, Y: float32(HEIGHT)/2 - SPRING_LEN - SPRING_Y_OFFSET}
		end   = rl.Vector2{X: sectionWidth / 2, Y: float32(HEIGHT)/2 - SPRING_Y_OFFSET}
	)

	for range SPRING_NUM {
		rl.DrawLineEx(start, end, SPRING_THICK, rl.LightGray)
		start.X += sectionWidth
		start, end = end, start
	}

}

func drawGround() {
	var (
		start = rl.Vector2{X: 0, Y: float32(HEIGHT) / 2}
		end   = rl.Vector2{X: float32(WIDTH), Y: float32(HEIGHT) / 2}
	)
	rl.DrawLineEx(start, end, float32(SPRING_THICK), rl.LightGray)
}

func acceleration(x, v float32) float32 {
	displacement := x - EQUILIBRIUM_X
	return (-SPRING_STIFFNESS / OBJ_MASS * displacement) - (DAMPING_COEFFICIENT / OBJ_MASS * v)
}
