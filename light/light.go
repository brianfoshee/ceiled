package light

import (
	"fmt"
	"image/color"
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
}

func (x X) Render() {
	fmt.Println("rendering")
}

func (x X) Open() {
	fmt.Println("Opening")
}

func (x X) Close() {
	fmt.Println("Closing")
}

func (x X) prepare() {
	// loop through Bars.LED and combine the values
	var c uint32
	var l Light
	var offset int
	for i := 0; i < len(x.Bars); i++ {
		bl := len(x.Bars[i].Lights)
		for j := 0; j < bl; j++ {
			l = x.Bars[i].Lights[j]
			c = uint32(l.White) << 24
			c = c | uint32(l.Color.R)<<16
			c = c | uint32(l.Color.G)<<8
			c = c | uint32(l.Color.B)
			// insert color at its correct index
			x.lights[offset+j] = c
		}
		offset += bl
	}
}
