package camera

import (
	"math"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func near(a, b float32) bool { return math.Abs(float64(a-b)) < 1e-4 }

func TestApproach(t *testing.T) {
	if got := approach(0, 10, 0, 0.016); got != 10 {
		t.Fatalf("rate 0 should snap, got %v", got)
	}
	if got := approach(0, 10, 8, 0); got != 0 {
		t.Fatalf("dt 0 should not move, got %v", got)
	}
	// Frame-rate independence: one 1/30 step equals two 1/60 steps.
	one := approach(0, 10, 8, 1.0/30)
	two := approach(approach(0, 10, 8, 1.0/60), 10, 8, 1.0/60)
	if !near(one, two) {
		t.Fatalf("not frame-rate independent: %v vs %v", one, two)
	}
}

func TestClampAxis(t *testing.T) {
	tests := []struct {
		name                          string
		target, min, size, view, want float32
	}{
		{"inside", 500, 0, 1000, 200, 500},
		{"left edge", 10, 0, 1000, 200, 100},
		{"right edge", 990, 0, 1000, 200, 900},
		{"world smaller than view centers", 50, 0, 100, 200, 50},
		{"offset world", 0, 100, 1000, 200, 200},
	}
	for _, tt := range tests {
		if got := clampAxis(tt.target, tt.min, tt.size, tt.view); !near(got, tt.want) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFirstPersonDirections(t *testing.T) {
	c := NewFirstPerson(rl.NewVector3(0, 2, 0), 70)

	d := c.Direction()
	if !near(d.X, 0) || !near(d.Y, 0) || !near(d.Z, -1) {
		t.Fatalf("yaw 0 should face -Z, got %+v", d)
	}
	if r := c.Right(); !near(r.X, 1) || !near(r.Z, 0) {
		t.Fatalf("right at yaw 0 = %+v", r)
	}

	// Mouse right turns right: forward gains +X.
	c.Look(100, 0)
	if f := c.Forward(); f.X <= 0 {
		t.Fatalf("mouse right should turn right, forward = %+v", f)
	}
	// Forward and Right stay perpendicular and unit length at any yaw.
	f, r := c.Forward(), c.Right()
	if !near(rl.Vector3DotProduct(f, r), 0) || !near(rl.Vector3Length(f), 1) || !near(rl.Vector3Length(r), 1) {
		t.Fatalf("forward %+v right %+v", f, r)
	}
}

func TestFirstPersonPitchClamp(t *testing.T) {
	c := NewFirstPerson(rl.Vector3Zero(), 70)
	c.Look(0, -1e6) // mouse far up
	if c.Pitch > maxPitch+1e-6 || c.Direction().Y >= 1 {
		t.Fatalf("pitch not clamped: %v", c.Pitch)
	}
	c.Look(0, 1e6) // mouse far down
	if c.Pitch < -maxPitch-1e-6 {
		t.Fatalf("pitch not clamped: %v", c.Pitch)
	}
}

func TestFirstPersonCameraFollowsEye(t *testing.T) {
	c := NewFirstPerson(rl.NewVector3(1, 2, 3), 70)
	c.SetEye(rl.NewVector3(5, 2, 5))
	if c.Camera.Position != rl.NewVector3(5, 2, 5) {
		t.Fatalf("position = %+v", c.Camera.Position)
	}
	want := rl.Vector3Add(c.Camera.Position, c.Direction())
	if c.Camera.Target != want {
		t.Fatalf("target = %+v, want %+v", c.Camera.Target, want)
	}
}
