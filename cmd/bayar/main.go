package main

import (
	"log"

	"github.com/jamesssooi/expense/pkg/bayar"
)

func main() {
	cfg, err := bayar.LoadConfig()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	log.Fatal(bayar.ListenAndServe(cfg.HostAddress, cfg.PortNumber))
}
