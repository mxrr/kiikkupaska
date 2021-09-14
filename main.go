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
		Loading: false,
		View:    utils.MAIN_MENU,
		Settings: utils.Settings{
			PanelVisible: false,
			Resolution:   utils.IVector2{X: int32(*widthFlag), Y: int32(*heightFlag)},
			Music:        *musicFlag,
		},
		RenderAssets: nil,
	}
}

func main() {
	utils.InitUtils(&state, debugMode)
	rl.InitWindow(state.Settings.Resolution.X, state.Settings.Resolution.Y, "Kiikkupaskaa")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
	rl.SetExitKey(rl.KeyF4)

	icon := rl.LoadImage("assets/fav.png")
	rl.SetWindowIcon(*icon)

	state.RenderAssets = rendering.LoadAssets(&state)
	rl.InitAudioDevice()

	var gameState *game.GameState

	var menuMusic rl.Music
	menuMusic = rl.LoadMusicStream(utils.GetAssetPath(utils.MUSIC, "main_menu01.mp3"))
	rl.PlayMusicStream(menuMusic)

	exitWindow := false

	go func(exit *bool) {
		for !(*exit) {
			if state.Settings.Music {
				rl.SetMusicVolume(menuMusic, 0.4)
				rl.UpdateMusicStream(menuMusic)
			}
		}
	}(&exitWindow)

	for !exitWindow {
		exitWindow = rl.WindowShouldClose()
		rl.SetWindowTitle(fmt.Sprintf("Kiikkupaskaa | %f fps %fms", rl.GetFPS(), rl.GetFrameTime()*1000.0))

		if state.Loading {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			rendering.DrawDefaultText(
				rl.NewVector2(
					state.Settings.Resolution.ToVec2().X/2.0,
					state.Settings.Resolution.ToVec2().Y/2.0,
				),
				48.0,
				"LOADING...",
				rl.RayWhite,
			)

			rl.EndDrawing()
		} else {
			switch state.View {
			//*
			//*	Main Menu UI
			//*
			//*
			case utils.MAIN_MENU:
				gameState = nil
				if rl.IsKeyPressed(rl.KeyEnter) {
					state.View = utils.IN_GAME
				}

				rl.BeginDrawing()

				rl.ClearBackground(rl.Black)

				rendering.DrawMenuButtons(state.View, &exitWindow)

				rl.EndDrawing()

			//*
			//*	Pause Menu UI
			//*
			//*
			case utils.PAUSED:
				if rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyP) {
					state.View = utils.IN_GAME
				}

				if rl.IsKeyPressed(rl.KeyQ) {
					state.View = utils.MAIN_MENU
				}

				rl.BeginDrawing()

				rl.ClearBackground(rl.Black)
				rendering.DrawMenuButtons(state.View, &exitWindow)

				rl.EndDrawing()

			//*
			//*	Game loop
			//*
			//*
			case utils.IN_GAME:
				game.GameUpdate(&state, &gameState)
			}
		}
	}

	rendering.Cleanup()

	rl.CloseAudioDevice()
	rl.CloseWindow()
}
