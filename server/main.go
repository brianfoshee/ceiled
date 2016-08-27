package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type LED struct {
	Brightness uint8
	White      uint8
	Red        uint8
	Green      uint8
	Blue       uint8
	mu         sync.RWMutex
}

func (l *LED) Set(br, w, r, g, b uint8) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Brightness = br
	l.White = w
	l.Red = r
	l.Green = g
	l.Blue = b
}

func main() {
	port := flag.String("port", "8080", "Port for server to list on")
	flag.Parse()

	l := LED{
		Brightness: 32,
		White:      255,
		Red:        0,
		Green:      0,
		Blue:       0,
		mu:         sync.RWMutex{},
	}

	idxtempl, err := template.New("index").Parse(index)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			l.mu.RLock()
			defer l.mu.RUnlock()
			err := idxtempl.Execute(w, l)
			if err != nil {
				fmt.Println(err)
			}
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "error reading form values %s", err)
				return
			}

			bright, err := strconv.ParseUint(r.Form.Get("brightness"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing brightness %s", err)
				return
			}
			white, err := strconv.ParseUint(r.PostForm.Get("white"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing white %s", err)
				return
			}
			red, err := strconv.ParseUint(r.PostForm.Get("red"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing red %s", err)
				return
			}
			green, err := strconv.ParseUint(r.PostForm.Get("green"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing green %s", err)
				return
			}
			blue, err := strconv.ParseUint(r.PostForm.Get("blue"), 10, 8)
			if err != nil {
				fmt.Fprintf(w, "error parsing blue %s", err)
				return
			}

			l.Set(uint8(bright), uint8(white), uint8(red), uint8(green), uint8(blue))
			l.Render()

			http.Redirect(w, r, "/", http.StatusFound)
		}
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
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
    <form action="/" method="post">
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

      <button type="submit">Set</button>
    </form>
  </body>
</html>
`
