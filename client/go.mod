module client

go 1.23.4

require (
	github.com/gorilla/websocket v1.5.3
	github.com/hajimehoshi/ebiten/v2 v2.8.8
	golang.org/x/image v0.20.0
// shared_types v0.0.0
// game_state v0.0.0
)

// replace shared_types => ../shared_types
// replace game_state => ../client/game_state

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
