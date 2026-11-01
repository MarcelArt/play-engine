package camera

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// maxPitch stops the view just short of straight up or down, where the look
// direction would flip.
const maxPitch = 89 * math.Pi / 180

// FirstPerson is an eye-level camera driven by yaw and pitch.
//
// Coordinates follow raylib: Y is up, the world is right-handed. Yaw 0 looks
// down -Z; positive yaw turns left. Pitch 0 is level; positive pitch looks up.
type FirstPerson struct {
	// Camera is the underlying raylib camera, ready for rl.BeginMode3D.
	Camera rl.Camera3D
	// Yaw and Pitch are in radians.
	Yaw, Pitch float32
	// Sensitivity is radians turned per pixel of mouse movement.
	Sensitivity float32

	eye rl.Vector3
}

// NewFirstPerson returns a camera at eye with the given vertical field of
// view in degrees.
func NewFirstPerson(eye rl.Vector3, fovy float32) *FirstPerson {
	c := &FirstPerson{
		Camera: rl.Camera3D{
			Up:         rl.NewVector3(0, 1, 0),
			Fovy:       fovy,
			Projection: rl.CameraPerspective,
		},
		Sensitivity: 0.0025,
	}
	c.SetEye(eye)
	return c
}

// Look turns the camera by a mouse movement in pixels. Moving the mouse right
// turns right and moving it down looks down. Pitch is clamped.
func (c *FirstPerson) Look(dx, dy float32) {
	c.Turn(-dx*c.Sensitivity, -dy*c.Sensitivity)
}

// Turn changes yaw and pitch by radians; use it for stick input, scaled by dt.
func (c *FirstPerson) Turn(yaw, pitch float32) {
	c.Yaw += yaw
	c.Pitch = rl.Clamp(c.Pitch+pitch, -maxPitch, maxPitch)
	c.sync()
}

// SetEye moves the camera to eye, keeping its direction.
func (c *FirstPerson) SetEye(eye rl.Vector3) {
	c.eye = eye
	c.sync()
}

// Eye is the camera position.
func (c *FirstPerson) Eye() rl.Vector3 { return c.eye }

// Direction is the unit vector the camera looks along.
func (c *FirstPerson) Direction() rl.Vector3 {
	sy, cy := float32(math.Sin(float64(c.Yaw))), float32(math.Cos(float64(c.Yaw)))
	sp, cp := float32(math.Sin(float64(c.Pitch))), float32(math.Cos(float64(c.Pitch)))
	return rl.NewVector3(-sy*cp, sp, -cy*cp)
}

// Forward is the unit look direction flattened onto the ground plane, so
// walking forward does not climb or dive with the pitch.
func (c *FirstPerson) Forward() rl.Vector3 {
	return rl.NewVector3(-float32(math.Sin(float64(c.Yaw))), 0, -float32(math.Cos(float64(c.Yaw))))
}

// Right is the unit vector to the camera's right on the ground plane.
func (c *FirstPerson) Right() rl.Vector3 {
	return rl.NewVector3(float32(math.Cos(float64(c.Yaw))), 0, -float32(math.Sin(float64(c.Yaw))))
}

// Begin starts drawing in 3D; End stops. Draw the HUD after End.
func (c *FirstPerson) Begin() { rl.BeginMode3D(c.Camera) }
func (c *FirstPerson) End()   { rl.EndMode3D() }

func (c *FirstPerson) sync() {
	c.Camera.Position = c.eye
	c.Camera.Target = rl.Vector3Add(c.eye, c.Direction())
}
