# Repository Guidelines

## Project Structure & Module Organization

This Go module, `github.com/MarcelArt/play-engine`, provides a scene-based game loop over raylib-go. `engine/` owns window configuration, scene lifecycle, and the fixed-timestep loop. Supporting packages are `assets/` (texture and font loading), `input/` (action bindings), `entity/` (entity lifecycle and collections), `save/` (JSON persistence), and `option/` (functional options). Tests live beside source files as `*_test.go`. `examples/asteroids/` contains a runnable game. The `assets/` directory contains loader code; game asset paths belong in manifests.

## Build, Test, and Development Commands

Use the Go version declared in `go.mod` (currently 1.27.3) and install raylib's native dependencies. Running games requires a graphical environment.

- `go build ./...`: compile all packages and the example.
- `go vet ./...`: run Go's static checks.
- `go test ./...`: run the complete test suite.
- `go test ./entity -run TestAddDuringUpdateRunsNextFrame`: run one focused test.
- `go run ./examples/asteroids` or `make asteroids`: launch the demo.

## Coding Style & Naming Conventions

Format Go changes with `gofmt`; use its standard tab indentation and import formatting. Keep package names lowercase, exported identifiers in PascalCase, and private identifiers in camelCase. Follow the example's `menu.scene.go` and `ship.entity.go` naming for scene and entity files. Document exported APIs with Go comments.

Preserve the architecture: register scene factories and input bindings in `main`, transition between scenes by name, initialize graphics-dependent state in `Enter`, and pass entities only the services they need. Draw directly with raylib; put simulation in `FixedUpdate`.

## Testing Guidelines

Use the standard `testing` package and descriptive `TestBehavior` names. Add regression tests for changed behavior, especially lifecycle ordering, missing assets, and save failures. Use `t.TempDir()` and `save.NewAt` for persistence tests. No coverage threshold is configured. Run tests and vet before submitting; report native-library or display limitations explicitly.

## Commit & Pull Request Guidelines

Git history is unavailable in this checkout, so historical commit conventions cannot be verified. Use concise, imperative commit subjects describing one focused change. In pull requests, explain the behavior change, link relevant issues, and list validation commands and results. Include screenshots or a short recording for visible demo changes. Update `README.md` and examples when public APIs change.
