package light

import (
	"fmt"
	"sync"
)

func (x X) Render() {
	x.prepare()
	fmt.Println("rendering")
}

func (x *X) Open() {
	var count int
	for _, b := range x.Bars {
		count += len(b.Lights)
	}

	x.lights = make([]uint32, count, count)
	x.Mu = &sync.RWMutex{}
	fmt.Println("Opening")
}

func (x X) Close() {
	fmt.Println("Closing")
}
