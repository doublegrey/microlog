package main

import (
	"log"

	"github.com/doublegrey/microlog/utils"
)

func main() {
	err := utils.Config.Parse()
	if err != nil {
		log.Fatalf("Failed to parse config: %v\n", err)
	}
}
