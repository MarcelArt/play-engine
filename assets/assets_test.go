package assets

import (
	"strings"
	"testing"
)

func TestMissingAssetPanicsWithName(t *testing.T) {
	m := New()
	for name, fn := range map[string]func(){
		"texture": func() { m.Texture("playr") },
		"font":    func() { m.Font("mono") },
	} {
		func() {
			defer func() {
				r := recover()
				if r == nil || !strings.Contains(r.(string), `"`) {
					t.Errorf("%s: want panic naming the asset, got %v", name, r)
				}
			}()
			fn()
		}()
	}
}

func TestUnloadAllTwice(t *testing.T) {
	m := New()
	m.UnloadAll()
	m.UnloadAll()
}
