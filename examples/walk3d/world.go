package main

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// box is an axis-aligned block standing on the floor.
type box struct {
	center rl.Vector3 // y is the vertical middle
	size   rl.Vector3
	color  rl.Color
}

func (b box) minX() float32 { return b.center.X - b.size.X/2 }
func (b box) maxX() float32 { return b.center.X + b.size.X/2 }
func (b box) minZ() float32 { return b.center.Z - b.size.Z/2 }
func (b box) maxZ() float32 { return b.center.Z + b.size.Z/2 }
func (b box) top() float32  { return b.center.Y + b.size.Y/2 }

// world is a flat floor with a deterministic scatter of blocks and a short
// staircase near the spawn point for trying out jumping.
type world struct {
	boxes []box
}

const (
	worldHalf  = 40 // the floor spans ±worldHalf on X and Z
	spawnClear = 5  // no random blocks this close to the origin
)

func newWorld() *world {
	w := &world{}

	// Stairs: 0.5 m steps along +X, to the right of the spawn.
	for i := range 4 {
		h := 0.5 * float32(i+1)
		w.add(rl.NewVector3(4+float32(i)*1.5, h/2, -3), rl.NewVector3(1.5, h, 3), i)
	}

	rng := rand.New(rand.NewPCG(7, 11))
	for i := range 45 {
		x := (rng.Float32()*2 - 1) * (worldHalf - 3)
		z := (rng.Float32()*2 - 1) * (worldHalf - 3)
		if math.Hypot(float64(x), float64(z)) < spawnClear {
			continue
		}
		sx, sz := 1+rng.Float32()*2.5, 1+rng.Float32()*2.5
		h := 0.8 + rng.Float32()*3.2
		w.add(rl.NewVector3(x, h/2, z), rl.NewVector3(sx, h, sz), i+4)
	}
	return w
}

func (w *world) add(center, size rl.Vector3, hueIndex int) {
	w.boxes = append(w.boxes, box{
		center: center,
		size:   size,
		color:  rl.ColorFromHSV(float32(hueIndex*37%360), 0.45, 0.85),
	})
}

func (w *world) Draw() {
	rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(worldHalf*2, worldHalf*2), rl.NewColor(96, 140, 90, 255))
	rl.DrawGrid(worldHalf*2, 1)
	for _, b := range w.boxes {
		rl.DrawCubeV(b.center, b.size, b.color)
		rl.DrawCubeWires(b.center, b.size.X, b.size.Y, b.size.Z, rl.NewColor(30, 30, 40, 255))
	}
}

// pushOut moves a vertical cylinder (feet at y, given radius) out of every
// block it overlaps sideways. Blocks whose top is within stepHeight of the
// feet are ignored, so the player can walk up small steps.
func (w *world) pushOut(pos *rl.Vector3, radius, stepHeight float32) {
	for _, b := range w.boxes {
		if pos.Y >= b.top()-stepHeight {
			continue
		}
		// Closest point on the block's footprint to the player's center.
		cx := rl.Clamp(pos.X, b.minX(), b.maxX())
		cz := rl.Clamp(pos.Z, b.minZ(), b.maxZ())
		dx, dz := pos.X-cx, pos.Z-cz
		dist := float32(math.Hypot(float64(dx), float64(dz)))
		if dist >= radius {
			continue
		}
		if dist == 0 {
			// Center is inside the footprint: leave through the nearest face.
			left, right := pos.X-b.minX(), b.maxX()-pos.X
			back, front := pos.Z-b.minZ(), b.maxZ()-pos.Z
			switch min(left, right, back, front) {
			case left:
				pos.X = b.minX() - radius
			case right:
				pos.X = b.maxX() + radius
			case back:
				pos.Z = b.minZ() - radius
			default:
				pos.Z = b.maxZ() + radius
			}
			continue
		}
		push := (radius - dist) / dist
		pos.X += dx * push
		pos.Z += dz * push
	}
}

// groundAt returns the height of the highest surface under (x, z) that is at
// or below maxY: a block top, or the floor at 0.
func (w *world) groundAt(x, z, maxY float32) float32 {
	ground := float32(0)
	for _, b := range w.boxes {
		if x < b.minX() || x > b.maxX() || z < b.minZ() || z > b.maxZ() {
			continue
		}
		if t := b.top(); t <= maxY && t > ground {
			ground = t
		}
	}
	return ground
}
