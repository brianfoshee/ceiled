package main

import (
	"flag"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

var maxBrightness = 0
var leds = 0

func main() {
	bright := flag.Int("brightness", 32, "LED brightness")
	count := flag.Int("count", 60, "How many LEDs in the strip")
	flag.Parse()

	maxBrightness = *bright
	leds = *count

	ws2811.Init(18, leds, maxBrightness)
	defer ws2811.Fini()

	ws2811.Clear()
	ws2811.Render()
	ws2811.Wait()

	for i := 0; i < leds; i++ {
		ws2811.SetLed(i, 0xff000000)
	}
	ws2811.Render()
	ws2811.Wait()
}
