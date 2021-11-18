package main

import "github.com/brianmartens/rover-images/cmd"

func main() {
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
