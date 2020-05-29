package draw

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// RecognizedDigit is recognized digit
var RecognizedDigit int

var (
	mousePos     rl.Vector2
	mousePosLast rl.Vector2
	canvas       rl.RenderTexture2D
	chosenColor  rl.Color
	brushSize    float32
	toClear      bool
)

const (
	windowHeight = 560
	windowWidth  = 560
)

// Run runs the application
func Run() {
	initialize()

	for !rl.WindowShouldClose() {
		update()
		draw()
	}

	// Unload texture from GPU memory
	rl.UnloadRenderTexture(canvas)
	rl.CloseWindow()
}

// UpdateDigit updates
func UpdateDigit(digit int) {
	RecognizedDigit = digit
}

func initialize() {
	rl.InitWindow(windowWidth, windowHeight, "GoDraw")
	rl.SetTargetFPS(120)

	// Load clean canvas
	canvas = rl.LoadRenderTexture(windowWidth, windowHeight)
	rl.BeginTextureMode(canvas)
	rl.ClearBackground(rl.White)
	rl.EndTextureMode()

	brushSize = 30
	chosenColor = rl.Black
}

func update() {
	mousePos = rl.GetMousePosition()

	// Save texture to png file
	if rl.IsKeyPressed(rl.KeyS) {
		image := rl.GetTextureData(canvas.Texture)
		rl.ImageFlipVertical(image)
		rl.ImageColorInvert(image)    // Invert color
		rl.ImageResize(image, 28, 28) // Resize image to size corresponding to MNIST dataset
		rl.ExportImage(*image, "images/image.png")
		rl.UnloadImage(image)

		// Clear canvas
		toClear = true
	}

	// Run neural network to recognize digit
	if rl.IsKeyPressed(rl.KeyR) {
		//RecognizedDigit = NN()
	}

	// Brush resizing
	if rl.IsKeyPressed(rl.KeyUp) {
		brushSize += 5
		if brushSize > 50 {
			brushSize = 50
		}
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		brushSize -= 5
		if brushSize < 10 {
			brushSize = 10
		}
	}
}

func draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)

	// Paint on a texture
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) || rl.GetGestureDetected() == rl.GestureDrag {
		rl.BeginTextureMode(canvas)
		// rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), brushSize, chosenColor)
		rl.DrawLineEx(mousePos, mousePosLast, brushSize*2, chosenColor)
		rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), brushSize, chosenColor)
		rl.EndTextureMode()
	}

	// Draw texture,
	// height is flipped because Raylib uses top-left corner starting point
	// but OpenGl uses bottom-left approach
	rl.DrawTextureRec(canvas.Texture, rl.Rectangle{X: 0, Y: 0, Width: windowWidth, Height: -windowHeight}, rl.Vector2{X: 0, Y: 0}, rl.White)

	// Draw brush "cursor" and outline
	rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), brushSize, chosenColor)
	rl.DrawCircleLines(int32(mousePos.X), int32(mousePos.Y), brushSize, rl.Black)
	rl.DrawText(fmt.Sprintf("%.f %.f", mousePos.X, mousePos.Y), 20, 20, 20, rl.Black)
	rl.DrawText(fmt.Sprintf("Recognized digit: %d", RecognizedDigit), 20, 50, 20, rl.Black)

	// Clear the canvas
	if rl.IsKeyPressed(rl.KeyC) || toClear {
		rl.BeginTextureMode(canvas)
		rl.ClearBackground(rl.White)
		rl.EndTextureMode()
		toClear = false
	}

	mousePosLast = mousePos

	rl.EndDrawing()
}
