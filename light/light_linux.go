package light

import (
	"fmt"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

func (x X) Render() {
	x.prepare()

	x.mu.RLock()
	defer x.mu.RUnlock()

	ws2811.SetBrightness(x.Brightness)
	ws2811.SetBitmap(x.lights)
	ws2811.Render()
	ws2811.Wait()
}

func (x X) Open() {
	var count int
	for _, b := range x.Bars {
		count += len(b.Lights)
	}
	ws2811.Init(18, count, x.Brightness)
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
