package main

import (
	"log"

	"github.com/sombre-hombre/mosaic/app/cmd"
)

func init() {
	log.SetFlags(log.Ltime)
}

func main() {
	cmd.Execute()
}
