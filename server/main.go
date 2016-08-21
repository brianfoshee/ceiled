package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "Port for server to list on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprintf(w, index)
		} else if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "error reading form values")
				return
			}

			bright := r.Form.Get("brightness")
			white := r.PostForm.Get("white")
			red := r.PostForm.Get("red")
			green := r.PostForm.Get("green")
			blue := r.PostForm.Get("blue")

			fmt.Printf("brightness: %s, white: %s, red: %s, green: %s, blue: %s\n",
				bright, white, red, green, blue)

			http.Redirect(w, r, "/", http.StatusFound)

			// TODO send to led package
		}
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

const index = `
<!DOCTYPE html>
<html>
  <head>
  </head>
  <body>
    <form action="/" method="post">
      <div>
        <input id="brightness" name="brightness" type="range" min="0" max="255" step="1" value="0" />
        <label for="brightness">Brightness</label>
      </div>

      <div>
        <input id="white" name="white" type="range" min="0" max="255" step="1" value="0" />
        <label for="white">W</label>
      </div>

      <div>
        <input id="red" name="red" type="range" min="0" max="255" step="1" value="0" />
        <label for="red">R</label>
      </div>

      <div>
        <input id="green" name="green" type="range" min="0" max="255" step="1" value="0" />
        <label for="green">G</label>
      </div>

      <div>
        <input id="blue" name="blue" type="range" min="0" max="255" step="1" value="0" />
        <label for="blue">B</label>
      </div>

      <button type="submit">Set</button>
    </form>
  </body>
</html>
`
