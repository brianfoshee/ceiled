package main

import (
	"fmt"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

func (l LED) Render() {
	ws2811.SetBrightness(int(l.Brightness))
	ws2811.SetBitmap(l.leds)
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
