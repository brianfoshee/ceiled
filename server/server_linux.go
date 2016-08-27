package main

import (
	"fmt"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

func (l LED) Render() {
	var color uint32
	color = uint32(l.White) << 24
	color = color | uint32(l.Red)<<16
	color = color | uint32(l.Green)<<8
	color = color | uint32(l.Blue)

	ws2811.SetBrightness(int(l.Brightness))
	for i := 0; i < 238; i++ {
		ws2811.SetLed(i, color)
	}
	ws2811.Render()
	ws2811.Wait()
}

func (l LED) Open() {
	ws2811.Init(18, 238, int(l.Brightness))
	ws2811.Clear()
	ws2811.Render()
	ws2811.Wait()

	l.Render()
}

func (l LED) Close() {
	ws2811.Clear()
	ws2811.Render()
	ws2811.Wait()
	ws2811.Fini()
	fmt.Println("Closed LEDs")
}
