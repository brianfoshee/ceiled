package main

import "fmt"

func (l LED) Render() {
	var color uint32
	color = uint32(l.White) << 24
	color = color | uint32(l.Red)<<16
	color = color | uint32(l.Green)<<8
	color = color | uint32(l.Blue)
	fmt.Printf("Color: %X, Brightness: %d, W: %d, R: %d, G: %d, B: %d\n",
		color, l.Brightness, l.White, l.Red, l.Green, l.Blue)
}
