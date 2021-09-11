package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"game"
	"rendering"
	"utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var state utils.State
var debugMode = false

func init() {
	runtime.LockOSThread()
}

// Handle cli arguments
func init() {
	widthFlag := flag.Int("w", 800, "Define the window width")
	heightFlag := flag.Int("h", 600, "Define the window height")
	musicFlag := flag.Bool("music", true, "Enable or disable music")
	debugFlag := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	log.Printf("Running with flags: -w %d -h %d -music=%v", *widthFlag, *heightFlag, *musicFlag)

	debugMode = *debugFlag
	state = utils.State{
		Loading:      false,
		View:         utils.MAIN_MENU,
		RES:          utils.IVector2{X: int32(*widthFlag), Y: int32(*heightFlag)},
		Music:        *musicFlag,
		RenderAssets: nil,
	}
}

func main() {
	rl.InitWindow(state.RES.X, state.RES.Y, "go-raylib")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
	rl.SetExitKey(rl.KeyF4)

	utils.InitUtils(&state, debugMode)
	state.RenderAssets = rendering.LoadAssets()
	rl.InitAudioDevice()

	var gameState *game.GameState

	var menuMusic rl.Music
	if state.Music {
		menuMusic = rl.LoadMusicStream(utils.GetAssetPath(utils.MUSIC, "main_menu01.mp3"))
		rl.PlayMusicStream(menuMusic)
	}

	exitWindow := false
	for !exitWindow {
		exitWindow = rl.WindowShouldClose()
		rl.SetWindowTitle(fmt.Sprintf("kiikkupaskaa | %f fps %fms", rl.GetFPS(), rl.GetFrameTime()*1000.0))

		switch state.View {
		case utils.MAIN_MENU:
			if state.Music {
				rl.SetMusicVolume(menuMusic, 0.4)
				rl.UpdateMusicStream(menuMusic)
			}

			if rl.IsKeyPressed(rl.KeyEnter) {
				state.View = utils.IN_GAME
			}

			rl.BeginDrawing()

			rl.ClearBackground(rl.Black)
			utils.DrawMainText(rl.Vector2{X: float32(state.RES.X/2 - 170), Y: float32(state.RES.Y / 6)}, 96.0, "Main Menu", rl.RayWhite)
			start := utils.DrawButton(rl.NewVector2(float32(state.RES.X)/2.0, float32(state.RES.Y)/2.0+50.0), "START")
			exit := utils.DrawButton(rl.NewVector2(float32(state.RES.X)/2.0, float32(state.RES.Y)/2.0+100.0), "QUIT")

			rl.EndDrawing()
			if start {
				state.View = utils.IN_GAME
			}
			if exit {
				exitWindow = true
			}
		case utils.PAUSED:
			if rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyM) {
				state.View = utils.IN_GAME
			}

			if rl.IsKeyPressed(rl.KeyQ) {
				gameState.AppState = nil
				state.View = utils.MAIN_MENU
			}

			rl.BeginDrawing()

			rl.ClearBackground(rl.Black)
			utils.DrawMainText(rl.Vector2{X: float32(state.RES.X) / 2.25, Y: float32(state.RES.Y) / 6.0}, 96.0, "Paused", rl.RayWhite)

			resume := utils.DrawButton(rl.NewVector2(float32(state.RES.X)/2.0, float32(state.RES.Y)/2.0+50.0), "RESUME")
			exit := utils.DrawButton(rl.NewVector2(float32(state.RES.X)/2.0, float32(state.RES.Y)/2.0+100.0), "EXIT TO MENU")

			rl.EndDrawing()

			if resume {
				state.View = utils.IN_GAME
			}

			if exit {
				gameState.AppState = nil
				state.View = utils.MAIN_MENU
			}
		case utils.IN_GAME:
			game.GameUpdate(&state, &gameState)
		}
	}

	rendering.Cleanup()

	rl.CloseAudioDevice()
	rl.CloseWindow()
}
