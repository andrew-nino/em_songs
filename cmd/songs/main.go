package main

import (
	"fmt"

	"github.com/andrew-nino/em_songs/config"
)

func main() {

	cfg := config.NewConfig()

	log := SetLogrus(cfg.Log.Level)

	log.Out.Write([]byte{65,66,'\n'})

	fmt.Println(cfg)
}
