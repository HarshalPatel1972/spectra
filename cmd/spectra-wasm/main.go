//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/HarshalPatel1972/spectra"
	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// scanSpectra is the function exposed to the JavaScript host.
// It takes three arguments: code (string), language (string), filename (string).
// It returns a JSON string containing the ScanResult.
func scanSpectra(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return js.ValueOf(`{"error": "requires 3 arguments: code, language, filename"}`)
	}

	code := args[0].String()
	language := ""
	if len(args) > 1 {
		language = args[1].String()
	}
	filename := "unknown"
	if len(args) > 2 {
		filename = args[2].String()
	}

	// Load patterns from the embedded rules
	registry, err := detector.LoadPatternsFromBytes(spectra.DefaultPatternsYAML)
	if err != nil {
		// Fallback for WASM environment if file is missing (should not happen with embed)
		registry = &detector.PatternRegistry{ByLanguage: make(map[string][]detector.PatternEntry)}
	}

	// Create a new orchestrator and scan entirely in memory
	res, err := scanner.ScanString(code, filename, language, registry)
	if err != nil {
		return js.ValueOf(fmt.Sprintf(`{"error": "scan failed: %v"}`, err))
	}

	// Serialize result to JSON
	out, err := json.Marshal(res)
	if err != nil {
		return js.ValueOf(fmt.Sprintf(`{"error": "json marshal failed: %v"}`, err))
	}

	return js.ValueOf(string(out))
}

func main() {
	// Expose the function to the global JavaScript object
	js.Global().Set("scanSpectra", js.FuncOf(scanSpectra))
	
	// Block forever so the Go WASM runtime stays alive
	select {}
}
