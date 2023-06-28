package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := play(v)
	fmt.Fprint(w, resp)
}

func play(input ArenaUpdate) (response string) {
	log.Printf("IN: %#v", input)

	self := input.Links.Self.Href
	var myState PlayerState
	width := input.Arena.Dimensions[0]
	height := input.Arena.Dimensions[1]
	arenaMap := make([][]string, height)
	for i := range arenaMap {
		arenaMap[i] = make([]string, width)
	}

	for k, v := range input.Arena.State {
		if k == self {
			myState = v
		}

		x := v.X
		y := v.Y
		log.Printf("%d, %d is %s", y, x, k)
		arenaMap[y][x] = k
	}

	if isSomeoneInFront(myState, arenaMap, width, height) {
		return "T"
	}

	commands := []string{"F", "R", "L"}
	rand := rand2.Intn(3)
	return commands[rand]
}

func isSomeoneInFront(state PlayerState, arena [][]string, width int, height int) bool {
	d := state.Direction
	x := state.X
	y := state.Y

	log.Printf("direction: %s, x: %d, y: %d", d, x, y)

	if d == "N" {
		if y-1 >= 0 && arena[y-1][x] != "" {
			log.Printf("Target %s at %d, %d", arena[y-1][x], y-1, x)
			return true
		}
		if y-2 >= 0 && arena[y-2][x] != "" {
			log.Printf("Target %s at %d, %d", arena[y-2][x], y-2, x)
			return true
		}
		if y-3 >= 0 && arena[y-3][x] != "" {
			log.Printf("Target %s at %d, %d", arena[y-3][x], y-3, x)
			return true
		}
	} else if d == "W" {
		if x-1 >= 0 && arena[y][x-1] != "" {
			log.Printf("Target %s at %d, %d", arena[y][x-1], y, x-1)
			return true
		}
		if x-2 >= 0 && arena[y][x-2] != "" {
			log.Printf("Target %s at %d, %d", arena[y][x-2], y, x-2)
			return true
		}
		if x-3 >= 0 && arena[y][x-3] != "" {
			log.Printf("Target %s at %d, %d", arena[y][x-3], y, x-3)
			return true
		}
	} else if d == "E" {
		if x+1 < width && arena[y][x+1] != "" {
			log.Printf("Target %s at %d, %d", arena[y][x+1], y, x+1)
			return true
		}
		if x+2 < width && arena[y][x+2] != "" {
			log.Printf("Target %s at %d, %d", arena[y][x+2], y, x+2)
			return true
		}
		if x+3 < width && arena[y][x+3] != "" {
			log.Printf("Target %s at %d, %d", arena[y][x+3], y, x+3)
			return true
		}
	} else if d == "S" {
		if y+1 < height && arena[y+1][x] != "" {
			log.Printf("Target %s at %d, %d", arena[y+1][x], y+1, x)
			return true
		}
		if y+2 < height && arena[y+2][x] != "" {
			log.Printf("Target %s at %d, %d", arena[y+2][x], y+2, x)
			return true
		}
		if y+3 < height && arena[y+3][x] != "" {
			log.Printf("Target %s at %d, %d", arena[y+3][x], y+3, x)
			return true
		}
	}

	return false
}
