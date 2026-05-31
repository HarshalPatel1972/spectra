package main

import (
	"fmt"
	"log"

	"github.com/HarshalPatel1972/spectra/internal/persistence"
)

func main() {
	store, err := persistence.NewStore("")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	scans, err := store.ListScans()
	if err != nil {
		log.Fatal(err)
	}

	if len(scans) > 0 {
		fmt.Print(scans[len(scans)-1].ID)
	}
}
