package light

import (
	"image/color"
	"sync"
)

type BarStatus int

const (
	BarStatusOff BarStatus = iota
	BarStatusOn
)

type Light struct {
	Color color.RGBA
	White uint8
}

type Bar struct {
	Lights []Light
	status BarStatus
	sync.RWMutex
}

func (b *Bar) SetStatus(s BarStatus) {
	b.Lock()
	b.status = s
	b.Unlock()
}

func (b *Bar) GetStatus() BarStatus {
	b.RLock()
	defer b.RUnlock()
	return b.status
}

type X struct {
	Bars       []Bar
	Brightness int
	lights     []uint32
	Mu         *sync.RWMutex
}

func (x X) prepare() {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	// Combine all Light values on each Bar into a []uint32
	var c uint32
	var offset int
	for _, b := range x.Bars {
		for j, l := range b.Lights {
			if b.GetStatus() == BarStatusOff {
				c = 0x0
			} else {
				c = uint32(l.White) << 24
				c = c | uint32(l.Color.R)<<16
				c = c | uint32(l.Color.G)<<8
				c = c | uint32(l.Color.B)
			}
			x.lights[offset+j] = c
		}
		offset += len(b.Lights)
	}
}
