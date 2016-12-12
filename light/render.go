package light

import (
	"context"
	"image/color"
)

type Light struct {
	Color color.RGBA
	White uint8
	Raw   uint32
}

func (l Light) RGBA() (r, g, b, a uint32) {
	return l.Color.RGBA()
}

type Lights []Light

func (l Lights) Raw() []uint32 {
	a := make([]uint32, len(l))
	for _, x := range lg.Lights {
		a = append(a, x.Raw)
	}
	return a
}

// Renderer submits light change commands. Implementations can be physical or virtual.
type Renderer interface {
	SetBrightness(int)
	SetPixels([]uint32)
}

// RendererFunc takes in a context and a Renderer object and performs some rendering
// until ctx is stopped.
type RendererFunc func(context.Context, Renderer)

// playing with ideas about how to implement passing around RendererFuncs
func meh() {
	c := make(chan RendererFunc, 100)

	f := func(ctx context.Context, r Renderer) {
	renderLoop:
		for {
			select {
			case <-ctx.Done():
				break renderLoop
			default:
				// do rendering
			}
		}
	}

	c <- f

	go func() {
		for rf := range c {
			//rf(ctx,
		}

	}()
}