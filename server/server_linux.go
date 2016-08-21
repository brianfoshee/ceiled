package main

import "github.com/jgarff/rpi_ws281x/golang/ws2811"

func (l LED) Render() {
	ws2811.Init(18, 120, int(l.Brightness))
	defer ws2811.Fini()

	var color uint32
	color = uint32(l.White) << 24
	color = color | uint32(l.Red)<<16
	color = color | uint32(l.Green)<<8
	color = color | uint32(l.Blue)

	for i := 0; i < 120; i++ {
		ws2811.SetLed(i, color)
	}
	ws2811.Render()
	ws2811.Wait()
}
