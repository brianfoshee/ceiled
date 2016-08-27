package light

import "fmt"

func (x X) Render() {
	x.prepare()
	fmt.Println("rendering")
}

func (x X) Open() {
	fmt.Println("Opening")
}

func (x X) Close() {
	fmt.Println("Closing")
}
