package palette

import (
	"fmt"

	"github.com/lucasb-eyer/go-colorful"
)

var globalPalette map[ID]Palette

// Init should be called before using this package
func Init(filePath string) {
	globalPalette = make(map[ID]Palette, 1)
	d := NewPalette(DefaultID, RGB)
	d.PutAt(colorful.Color{R: 1.0, G: 0.0, B: 0.0}, 0.0)
	d.PutAt(colorful.Color{R: 0.0, G: 1.0, B: 0.0}, 1.0/3.0)
	d.PutAt(colorful.Color{R: 0.0, G: 0.0, B: 1.0}, 2.0/3.0)
	d.PutAt(colorful.Color{R: 1.0, G: 0.0, B: 0.0}, 1.0)
	SetGlobal(d)
}

// GetGlobal gets all the global palettes
func GetGlobal() map[ID]Palette {
	return globalPalette
}

// GetGlobalByID gets the global palette by id
func GetGlobalByID(id ID) (Palette, error) {
	p, ok := globalPalette[id]
	if !ok {
		return globalPalette[DefaultID], fmt.Errorf("palette not found with id %s", id)
	}
	return p, nil
}

// SetGlobal gets the global palette by id
func SetGlobal(p Palette) {
	globalPalette[p.GetID()] = p
}
