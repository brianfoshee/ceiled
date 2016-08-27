package light

import (
	"fmt"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

func (x X) Render() {
	ws2811.SetBrightness(x.Brightness)
	ws2811.SetBitmap(x.lights)
	ws2811.Render()
	ws2811.Wait()
}

func (x X) Open() {
	ws2811.Init(18, len(x.lights), x.Brightness)
	ws2811.Clear()
	ws2811.Render()
	ws2811.Wait()

	x.Render()
}

func (x X) Close() {
	ws2811.Clear()
	ws2811.Render()
	ws2811.Wait()
	ws2811.Fini()
	fmt.Println("Closed LEDs")
}
