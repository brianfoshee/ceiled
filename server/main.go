package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/brianfoshee/ceiled/light"
)

type LED struct {
	Brightness uint8
	White      uint8
	Red        uint8
	Green      uint8
	Blue       uint8
	sync.RWMutex

	Bar0, Bar1, Bar2, Bar3 light.BarStatus
}

func (l *LED) SetPower(x *light.X, b0, b1, b2, b3 light.BarStatus) {
	l.Lock()
	defer l.Unlock()

	l.Bar0 = b0
	l.Bar1 = b1
	l.Bar2 = b2
	l.Bar3 = b3

	x.Bars[0].SetStatus(b0)
	x.Bars[1].SetStatus(b1)
	x.Bars[2].SetStatus(b2)
	x.Bars[3].SetStatus(b3)
}

func (l *LED) SetColor(x *light.X, br, w, r, g, b uint8) {
	l.Lock()
	defer l.Unlock()

	l.Brightness = br
	l.White = w
	l.Red = r
	l.Green = g
	l.Blue = b

	x.Brightness = int(br)

	x.Mu.Lock()
	for i := 0; i < len(x.Bars); i++ {
		for j := 0; j < len(x.Bars[i].Lights); j++ {
			x.Bars[i].Lights[j].Color.R = r
			x.Bars[i].Lights[j].Color.G = g
			x.Bars[i].Lights[j].Color.B = b
			x.Bars[i].Lights[j].White = w
		}
	}
	x.Mu.Unlock()
}

func main() {
	port := flag.String("port", "8080", "Port for server to list on")
	flag.Parse()

	idxtempl, err := template.New("index").Parse(index)
	if err != nil {
		fmt.Println(err)
		return
	}

	x := light.X{
		Bars: []light.Bar{
			{Lights: make([]light.Light, 60, 60)},
			{Lights: make([]light.Light, 60, 60)},
			{Lights: make([]light.Light, 58, 58)},
			{Lights: make([]light.Light, 60, 60)},
		},
	}
	x.Open()

	l := LED{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			l.RLock()
			defer l.RUnlock()
			err := idxtempl.Execute(w, l)
			if err != nil {
				fmt.Println(err)
			}
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "error reading form values %s", err)
				return
			}

			r.ParseForm()
			r.ParseMultipartForm(1024 * 10)

			formValue := func(r *http.Request, e string) string {
				mf, ok := r.MultipartForm.Value["brightness"]
				if ok && len(mf) > 0 {
					return mf[0]
				}
				return r.Form.Get("brightness")
			}

			bright, err := strconv.ParseUint(formValue(r, "brightness"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing brightness %s", err)
				return
			}
			white, err := strconv.ParseUint(formValue(r, "white"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing white %s", err)
				return
			}
			red, err := strconv.ParseUint(formValue(r, "red"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing red %s", err)
				return
			}
			green, err := strconv.ParseUint(formValue(r, "green"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing green %s", err)
				return
			}
			blue, err := strconv.ParseUint(formValue(r, "blue"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing blue %s", err)
				return
			}

			var b0, b1, b2, b3 light.BarStatus

			v, ok := r.Form["bar0"]
			if ok && len(v) > 0 && v[0] == "on" {
				b0 = light.BarStatusOn
			}

			v, ok = r.Form["bar1"]
			if ok && len(v) > 0 && v[0] == "on" {
				b1 = light.BarStatusOn
			}

			v, ok = r.Form["bar2"]
			if ok && len(v) > 0 && v[0] == "on" {
				b2 = light.BarStatusOn
			}

			v, ok = r.Form["bar3"]
			if ok && len(v) > 0 && v[0] == "on" {
				b3 = light.BarStatusOn
			}

			l.SetPower(&x, b0, b1, b2, b3)
			l.SetColor(&x, uint8(bright), uint8(white), uint8(red), uint8(green), uint8(blue))
			go func() {
				x.Render()
			}()

			http.Redirect(w, r, "/", http.StatusFound)
		}
	})

	addr, err := net.ResolveTCPAddr("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		// Block until signal is received
		<-c
		// Close LED things
		x.Close()
		// Close TCP Listener
		if err := listener.Close(); err != nil {
			log.Println(err)
		}
		// Proceed with exiting the program
		done <- struct{}{}
	}()

	if err := http.Serve(listener, nil); err != nil {
		fmt.Println(err)
	}
	<-done
}

const index = `
<!DOCTYPE html>
<html>
  <head>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		form {
			width: 216px;
			position: relative;
		}

		label {
			float: left;
			text-align: right;
			margin-right: 15px;
			width: 68px;
		}

		button {
			float: right;
		}
	</style>
  </head>
  <body>
    <form id="lights" action="/" method="post">
      <!--
        <input id="brightness" name="brightness" type="range" min="0" max="255" step="1" value="128" oninput="brightnessInput.value=brightness.value"/>
        <input id="brightnessInput" name="brightnessInput" type="text" for="brightness" value="128" oninput="brightness.value=brightnessInput.value" />
      -->
      <div>
        <input id="brightness" name="brightness" type="range" min="0" max="255" step="1" value="{{ .Brightness }}" />
        <label for="brightness">Brightness</label>
      </div>

      <div>
        <input id="white" name="white" type="range" min="0" max="255" step="1" value="{{ .White }}" />
        <label for="white">White</label>
      </div>

      <div>
        <input id="red" name="red" type="range" min="0" max="255" step="1" value="{{ .Red }}" />
        <label for="red">Red</label>
      </div>

      <div>
        <input id="green" name="green" type="range" min="0" max="255" step="1" value="{{ .Green }}" />
        <label for="green">Green</label>
      </div>

      <div>
        <input id="blue" name="blue" type="range" min="0" max="255" step="1" value="{{ .Blue }}" />
        <label for="blue">Blue</label>
      </div>

	  <div>
		<input id="bar0" name="bar0" type="checkbox" {{ if eq .Bar0 1 }} checked {{ end }} />
		<input id="bar1" name="bar1" type="checkbox" {{ if eq .Bar1 1 }} checked {{ end }} />
		<input id="bar2" name="bar2" type="checkbox" {{ if eq .Bar2 1 }} checked {{ end }} />
		<input id="bar3" name="bar3" type="checkbox" {{ if eq .Bar3 1 }} checked {{ end }} />
	  </div>

      <button type="submit">Set</button>
    </form>

	<script>
		const form = document.getElementById("lights");

		function submit() {
			var xhr = new XMLHttpRequest();
			xhr.open(form.method, form.action, true);
			xhr.onload = function(){ console.log(xhr.responseText); }
			const formData = new FormData(form);
			xhr.send(formData);
		}

		const brightness = document.getElementById("brightness");
		brightness.onchange = submit;

		const red = document.getElementById("red");
		red.onchange = submit;

		const green = document.getElementById("green");
		green.onchange = submit;

		const blue = document.getElementById("blue");
		blue.onchange = submit;

		const white = document.getElementById("white");
		white.onchange = submit;
	</script>
  </body>
</html>
`
