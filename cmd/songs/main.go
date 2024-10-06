package main

import (
	"fmt"

	"github.com/andrew-nino/em_songs/config"
)

func main() {

	cfg := config.NewConfig()

	fmt.Println(cfg)
}
