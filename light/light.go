package light

import (
	"image/color"
	"sync"
)

type Light struct {
	Color color.RGBA
	White uint8
}

type Bar struct {
	Lights []Light
}

type X struct {
	Bars       []Bar
	Brightness int
	lights     []uint32
	mu         sync.RWMutex
}

func (x *X) prepare() {
	x.mu.Lock()
	defer x.mu.Unlock()

	// Combine all Light values on each Bar into a []uint32
	var c uint32
	var offset int
	for _, b := range x.Bars {
		for j, l := range b.Lights {
			c = uint32(l.White) << 24
			c = c | uint32(l.Color.R)<<16
			c = c | uint32(l.Color.G)<<8
			c = c | uint32(l.Color.B)
			// insert color at its correct index
			x.lights[offset+j] = c
		}
		offset += len(b.Lights)
	}
}
