package draw

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Text can be additional information to display from the outside of package
var Text string

var (
	mousePos     rl.Vector2
	mousePosLast rl.Vector2
	canvas       rl.RenderTexture2D
	chosenColor  rl.Color
	brushSize    float32
	toClear      bool
	autoSave     bool
	autoSaveTime int
	frame        int
)

const (
	windowHeight = 560
	windowWidth  = 560
	targetFPS    = 120
)

// Run runs the application
func Run(autosave bool) {
	autoSave = autosave
	initialize()

	for !rl.WindowShouldClose() {
		update()
		draw()
	}

	// Unload texture from GPU memory
	rl.UnloadRenderTexture(canvas)
	rl.CloseWindow()
}

func initialize() {
	rl.InitWindow(windowWidth, windowHeight, "GoDraw")
	rl.SetTargetFPS(targetFPS)

	// Load clean canvas
	canvas = rl.LoadRenderTexture(windowWidth, windowHeight)
	rl.BeginTextureMode(canvas)
	rl.ClearBackground(rl.White)
	rl.EndTextureMode()
	if autoSave {
		autoSaveTime = targetFPS
	}
	autoSaveTime = 1
	brushSize = 25
	chosenColor = rl.Black
}

func update() {
	mousePos = rl.GetMousePosition()
	frame++

	// Save texture to png file
	if rl.IsKeyPressed(rl.KeyS) || frame%autoSaveTime == 0 {
		image := rl.GetTextureData(canvas.Texture)
		rl.ImageFlipVertical(image)
		rl.ImageColorInvert(image)    // Invert color
		rl.ImageResize(image, 28, 28) // Resize image to size corresponding to MNIST dataset
		// Check if images directory exists, create one if not
		if _, err := os.Stat("/images"); os.IsNotExist(err) {
			os.Mkdir("images", os.ModePerm)
		}
		rl.ExportImage(*image, "images/image.png")
		rl.UnloadImage(image)

		// Clear canvas
		// toClear = true
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
	rl.DrawText(fmt.Sprintf("%s", Text), 20, 50, 20, rl.Black)

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
