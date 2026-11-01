package entity

import "testing"

type dummy struct {
	Lifecycle
	updates  int
	onUpdate func()
}

func (d *dummy) Update(float32) {
	d.updates++
	if d.onUpdate != nil {
		d.onUpdate()
	}
}
func (d *dummy) Draw() {}

func TestSweepRemovesDead(t *testing.T) {
	var w World[*dummy]
	a, b, c := &dummy{}, &dummy{}, &dummy{}
	w.Add(a)
	w.Add(b)
	w.Add(c)
	b.Kill()
	w.Update(0.1)
	if w.Len() != 2 || b.updates != 0 || a.updates != 1 || c.updates != 1 {
		t.Fatalf("len=%d a=%d b=%d c=%d", w.Len(), a.updates, b.updates, c.updates)
	}
}

func TestAddDuringUpdateRunsNextFrame(t *testing.T) {
	var w World[*dummy]
	spawned := &dummy{}
	w.Add(&dummy{onUpdate: func() { w.Add(spawned) }})
	w.Update(0.1)
	if spawned.updates != 0 || w.Len() != 2 {
		t.Fatalf("spawned updated early: updates=%d len=%d", spawned.updates, w.Len())
	}
	w.Update(0.1)
	if spawned.updates != 1 {
		t.Fatalf("spawned.updates = %d", spawned.updates)
	}
}

func TestKillDuringUpdateSweeps(t *testing.T) {
	var w World[*dummy]
	d := &dummy{}
	d.onUpdate = d.Kill
	w.Add(d)
	w.Update(0.1)
	if w.Len() != 0 {
		t.Fatal("dead entity not swept")
	}
}
