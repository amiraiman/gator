package main

import (
	"fmt"

	"github.com/amiraiman/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("Amir Aiman")

	cfg = config.Read()
	fmt.Print(cfg)
}
