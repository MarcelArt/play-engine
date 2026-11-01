# Repository Guidelines

## Project Structure & Module Organization

This Go module, `github.com/MarcelArt/play-engine`, provides a scene-based game loop over raylib-go. `engine/` owns window/display configuration, scene lifecycle, and the fixed-timestep loop. Supporting packages are `assets/` (texture/font loading), `input/` (action bindings), `entity/` (entity collections), `save/` (progress persistence), `settings/` (preferences with defaults), and `option/` (functional options). Tests live beside source as `*_test.go`.

`examples/asteroids/` uses wireframe graphics; `examples/asteroids_2d/` uses embedded PNGs and persistent display settings. Root `assets/` contains loader code; game assets belong inside their example package. `skills/` contains contributor-maintained agent skills and an API reference.

## Build, Test, and Development Commands

Use the Go version declared in `go.mod` and install raylib's native dependencies. Running games requires a graphical environment.

- `go build ./...`: compile all packages and both examples.
- `go vet ./...`: run Go's static checks.
- `go test ./...`: run the complete test suite.
- `go test ./entity -run TestAddDuringUpdateRunsNextFrame`: run one focused test.
- `make asteroids` / `make asteroids2d`: launch the wireframe / sprite demo.
- `go run ./examples/asteroids_2d`: launch the sprite demo directly.

## Coding Style & Naming Conventions

Format Go changes with `gofmt` (tabs). Use lowercase package names, PascalCase exported identifiers, camelCase private identifiers, and Go comments for exported APIs. Follow `menu.scene.go` and `ship.entity.go` naming; import raylib as `rl`.

Register scenes, bindings, and manifests in `main`; use constants for names. Initialize graphics-dependent state in `Enter`, read input in `Update`, simulate in `FixedUpdate`, and draw with raylib. Pass entities only needed services.

Use `WithDisplay` initially and `ctx.ApplyDisplay` for runtime window changes. Layouts should read current screen dimensions. Store preferences in `settings`, progress in `save`; preserve defaults for absent settings keys.

## Testing Guidelines

Use standard `testing` and descriptive `TestBehavior` names. Cover lifecycle ordering, missing assets, persistence failures, settings defaults, and display normalization. Use `t.TempDir()` with `save.NewAt` or `settings.NewAt` for persistence tests. No coverage threshold is configured. Run tests and vet before submitting; report native-library or display limitations.

## Commit & Pull Request Guidelines

History includes `feat: display settings` and `ai: make skills`; use concise, imperative subjects with a relevant prefix. PRs should explain behavior changes, link relevant issues, and list validation results. Include screenshots or recordings for visible demo changes. Update `README.md`, examples, and applicable `skills/` guidance when public APIs change.
