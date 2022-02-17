package main

import (
	"log"

	"github.com/ariary/cfuzz/pkg/config"
	"github.com/ariary/cfuzz/pkg/fuzz"
	"github.com/ariary/cfuzz/pkg/output"
	//"pkg/output"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) //set default logger
	// config & banner
	cfg := config.NewConfig()

	output.Banner()
	output.PrintConfig(cfg)

	if err := cfg.CheckConfig(); err != nil {
		log.Fatal(err)
	}

	fuzz.PerformFuzzing(cfg)

}
