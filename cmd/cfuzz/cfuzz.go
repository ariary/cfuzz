package main

import (
	"log"

	"github.com/ariary/cfuzz/pkg/fuzz"
)

func main() {
	//log.SetFlags(log.Lshortfile) //set default logger
	log.SetFlags(0)

	// config & banner
	cfg := fuzz.NewConfig()

	fuzz.Banner()
	fuzz.PrintConfig(cfg)

	if err := cfg.CheckConfig(); err != nil {
		log.Fatal(err)
	}

	fuzz.PerformFuzzing(cfg)

}
