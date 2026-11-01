package engine

import "testing"

func TestNormalize(t *testing.T) {
	got := Display{Width: 0, Height: -5, FPS: -1, VSync: true}.Normalize()
	if got.Width != minDisplayWidth || got.Height != minDisplayHeight || got.FPS != 0 || !got.VSync {
		t.Fatalf("Normalize = %+v", got)
	}
	ok := Display{Width: 1280, Height: 720, FPS: 60}
	if ok.Normalize() != ok {
		t.Fatal("valid display changed")
	}
}

func TestNextResolution(t *testing.T) {
	if r := NextResolution(1280, 720); r != (Resolution{1600, 900}) {
		t.Fatalf("next of 1280x720 = %+v", r)
	}
	last := CommonResolutions[len(CommonResolutions)-1]
	if r := NextResolution(last.Width, last.Height); r != CommonResolutions[0] {
		t.Fatalf("no wrap: %+v", r)
	}
	if r := NextResolution(777, 333); r != CommonResolutions[0] {
		t.Fatalf("unknown size: %+v", r)
	}
}
